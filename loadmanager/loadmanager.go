// GO_GLM.go
package GO_GLM

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/fatih/color"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
)

var (
	goroutinesInProgress = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "goroutines_in_progress",
		Help: "Number of goroutines currently in progress.",
	})
	errorsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "errors_total",
		Help: "Total number of errors encountered.",
	})
	executionDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "execution_duration_seconds",
		Help:    "Duration of function execution in seconds.",
		Buckets: prometheus.DefBuckets,
	})
	semaphoreWaitDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "semaphore_wait_duration_seconds",
		Help:    "Duration each goroutine waits for a semaphore in seconds.",
		Buckets: prometheus.DefBuckets,
	})
	successfulExecutions = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "successful_executions_total",
		Help: "Total number of successfully executed goroutines.",
	})
	goroutineQueueLength = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "goroutine_queue_length",
		Help: "Number of goroutines waiting to acquire the semaphore.",
	})
	totalProcessingTime = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "total_processing_time_seconds",
		Help: "Cumulative processing time of all goroutines in seconds.",
	})
)

func init() {
	prometheus.MustRegister(goroutinesInProgress, errorsCounter, executionDuration, semaphoreWaitDuration, successfulExecutions, goroutineQueueLength, totalProcessingTime)
}

// Constants for load calculation and thresholds
const (
	cpuWeight             = 0.7
	memWeight             = 0.3
	criticalLoadThreshold = 90.0
	minGoroutines         = 1    // Minimum goroutines allowed
	maxGoroutines         = 100  // Maximum goroutines allowed
	highCPULoad           = 80.0 // CPU usage threshold to start scaling down
	lowCPULoad            = 30.0 // CPU usage threshold to start scaling up
)

// Slices to store load history
var (
	last1MinLoadHistory  []float64
	last5MinLoadHistory  []float64
	last15MinLoadHistory []float64
)

// Variable to store the live load
var liveLoad float64

// CalculateLoad calculates the combined CPU and memory load.
func CalculateLoad() float64 {
	cpuPercent, _ := cpu.Percent(time.Second, false)
	memInfo, _ := mem.VirtualMemory()

	// Calculate memory usage percentage
	memUsagePercent := (float64(memInfo.Used) / float64(memInfo.Total)) * 100

	// Calculate combined load
	return cpuPercent[0]*cpuWeight + memUsagePercent*memWeight
}

type AdjustConcurrencyResponse struct {
	concurrency int
	liveLoad    float64
	last1Min    float64
	last5Min    float64
	last15Min   float64
}

// AdjustConcurrency adjusts concurrency based on load in the background
// and returns the calculated loads in the last 1, 5, and 15 minutes.
func AdjustConcurrency(maxConcurrency int, loadThreshold float64) AdjustConcurrencyResponse {
	concurrencyChan := make(chan int)
	if loadThreshold == 0 {
		loadThreshold = highCPULoad
	}
	go func() {
		concurrency := maxConcurrency / 2
		concurrencyChan <- concurrency

		for {
			liveLoad = CalculateLoad()
			last1MinLoadHistory = append(last1MinLoadHistory, liveLoad)
			last5MinLoadHistory = append(last5MinLoadHistory, liveLoad)
			last15MinLoadHistory = append(last15MinLoadHistory, liveLoad)

			// Trim history slices to maintain the desired time windows
			if len(last1MinLoadHistory) > 60 {
				last1MinLoadHistory = last1MinLoadHistory[1:]
			}
			if len(last5MinLoadHistory) > 300 {
				last5MinLoadHistory = last5MinLoadHistory[1:]
			}
			if len(last15MinLoadHistory) > 900 {
				last15MinLoadHistory = last15MinLoadHistory[1:]
			}

			if liveLoad > loadThreshold && concurrency < maxConcurrency {
				concurrency++
			} else if liveLoad < loadThreshold && concurrency > minGoroutines {
				concurrency--
			}
			concurrencyChan <- concurrency
			time.Sleep(time.Second)
		}
	}()

	for {
		concurrency := <-concurrencyChan

		if len(last1MinLoadHistory) == 60 || len(last5MinLoadHistory) == 300 || len(last15MinLoadHistory) == 900 {
			last1Min := average(last1MinLoadHistory)
			last5Min := average(last5MinLoadHistory)
			last15Min := average(last15MinLoadHistory)
			return AdjustConcurrencyResponse{
				concurrency: concurrency,
				liveLoad:    liveLoad,
				last1Min:    last1Min,
				last5Min:    last5Min,
				last15Min:   last15Min,
			}
		}
	}
}

// Function to calculate the average of a slice of floats
func average(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// Go runs without taking a limit input and dynamically adjusts concurrency limit.
func Go[T any](orgCtx context.Context, funcName string, arr []T, f func(context.Context, T) error) error {
	// Set an initial limit to be dynamically adjusted based on load
	initialLimit := 10 // Starting point
	return GoWithCustomLimit(orgCtx, funcName, arr, initialLimit, f)
}

// GoWithCustomLimit dynamically adjusts concurrency limit based on load.
func GoWithCustomLimit[T any](ctx context.Context, funcName string, arr []T, limit int, f func(context.Context, T) error) error {
	if limit < minGoroutines {
		limit = minGoroutines
	}

	g, ctx := errgroup.WithContext(ctx)
	semaphore := make(chan struct{}, limit)
	errChan := make(chan error, 1)

	// Set up a ticker to adjust the goroutine limit dynamically every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			res := AdjustConcurrency(limit, 0)
			limit = res.concurrency
		}
	}()

	for _, ar := range arr {
		ar := ar

		g.Go(func() error {
			goroutinesInProgress.Inc()       // Increment the number of in-progress goroutines
			defer goroutinesInProgress.Dec() // Decrement when goroutine completes

			select {
			case <-ctx.Done():
				return ctx.Err()
			case semaphore <- struct{}{}: // Acquire semaphore
			}

			defer func() {
				<-semaphore // Release semaphore
			}()

			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("recover %v", r)
					fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
					select {
					case errChan <- err:
						errorsCounter.Inc() // Increment error counter
					default:
					}
				}
			}()

			if err := f(ctx, ar); err != nil {
				select {
				case errChan <- err:
					color.Red(err.Error())
					errorsCounter.Inc() // Increment error counter
				default:
				}
				return err
			}
			return nil
		})
	}

	color.Green("now waiting ...  %s", funcName)
	if err := g.Wait(); err != nil {
		return err
	}
	color.Green("done waiting ... %s", funcName)

	close(errChan)
	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
}

func GetBackgroundContextWithMd(ctx context.Context) (context.Context, context.CancelFunc, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md, ok = metadata.FromOutgoingContext(ctx)
		if !ok {
			return nil, nil, errors.New("request rejected due to internal server error")
		}
	}
	mdCopy := md.Copy()
	newCtx := metadata.NewOutgoingContext(context.Background(), mdCopy)
	newCtx, cancel := context.WithCancel(newCtx)
	return newCtx, cancel, nil
}
