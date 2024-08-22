// loadmanager_test.go
package loadmanager

import (
	"testing"
	"time"
)

func TestCalculateLoad(t *testing.T) {
	load := CalculateLoad()
	if load < 0 || load > 100 {
		t.Errorf("Calculated load is out of range: %.2f%%", load)
	}
}

func TestAdjustConcurrency(t *testing.T) {
	maxConcurrency := 10
	loadThreshold := 80.0
	concurrency, liveLoad, last1Min, last5Min, last15Min := AdjustConcurrency(maxConcurrency, loadThreshold)

	if concurrency < 1 || concurrency > maxConcurrency {
		t.Errorf("Concurrency is out of range: %d", concurrency)
	}
	if liveLoad < 0 || liveLoad > 100 {
		t.Errorf("Live load is out of range: %.2f%%", liveLoad)
	}
	if last1Min < 0 || last1Min > 100 || last5Min < 0 || last5Min > 100 || last15Min < 0 || last15Min > 100 {
		t.Errorf("Average loads are out of range: last1Min=%.2f, last5Min=%.2f, last15Min=%.2f",
			last1Min, last5Min, last15Min)
	}
}

// This test simulates the AdjustConcurrency functionality over time
func TestConcurrencyAdjustment(t *testing.T) {
	maxConcurrency := 10
	loadThreshold := 80.0
	go AdjustConcurrency(maxConcurrency, loadThreshold)

	time.Sleep(5 * time.Second) // Wait for some load adjustments

	concurrency, liveLoad, last1Min, last5Min, last15Min := AdjustConcurrency(maxConcurrency, loadThreshold)

	if concurrency < 1 || concurrency > maxConcurrency {
		t.Errorf("Concurrency is out of range: %d", concurrency)
	}
	if liveLoad < 0 || liveLoad > 100 {
		t.Errorf("Live load is out of range: %.2f%%", liveLoad)
	}
	if last1Min < 0 || last1Min > 100 || last5Min < 0 || last5Min > 100 || last15Min < 0 || last15Min > 100 {
		t.Errorf("Average loads are out of range: last1Min=%.2f, last5Min=%.2f, last15Min=%.2f",
			last1Min, last5Min, last15Min)
	}
}
