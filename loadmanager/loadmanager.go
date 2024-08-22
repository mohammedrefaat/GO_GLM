// loadmanager.go
package loadmanager

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// Constants for load calculation and thresholds
const (
	cpuWeight             = 0.7
	memWeight             = 0.3
	criticalLoadThreshold = 90.0
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

// AdjustConcurrency adjusts concurrency based on load in the background
// and returns the calculated loads in the last 1, 5, and 15 minutes.
func AdjustConcurrency(maxConcurrency int, loadThreshold float64) (int, float64, float64, float64, float64) {
	concurrencyChan := make(chan int)

	go func() {
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

			concurrency := <-concurrencyChan
			if liveLoad > loadThreshold && concurrency < maxConcurrency {
				concurrency++
			} else if liveLoad < loadThreshold && concurrency > 1 {
				concurrency--
			}
			concurrencyChan <- concurrency
			time.Sleep(time.Second)
		}
	}()

	concurrency := maxConcurrency / 2
	concurrencyChan <- concurrency

	for {
		concurrency = <-concurrencyChan
		concurrencyChan <- concurrency

		if len(last1MinLoadHistory) == 60 || len(last5MinLoadHistory) == 300 || len(last15MinLoadHistory) == 900 {
			last1Min := average(last1MinLoadHistory)
			last5Min := average(last5MinLoadHistory)
			last15Min := average(last15MinLoadHistory)
			return concurrency, liveLoad, last1Min, last5Min, last15Min
		}
	}
}

// Function to calculate the average of a slice of floats
func average(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}
