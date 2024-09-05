// GO_GLM_test.go
package GO_GLM

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// Mock function for testing
func mockFunction(ctx context.Context, input int) error {
	time.Sleep(10 * time.Millisecond) // Simulate some work
	if input < 0 {
		return fmt.Errorf("invalid input: %d", input)
	}
	return nil
}

func TestCalculateLoad(t *testing.T) {
	load := CalculateLoad()
	if load < 0 || load > 100 {
		t.Errorf("Calculated load is out of range: %.2f%%", load)
	}
}

func TestAdjustConcurrency(t *testing.T) {
	maxConcurrency := 10
	loadThreshold := 80.0
	response := AdjustConcurrency(maxConcurrency, loadThreshold)

	if response.concurrency < 1 || response.concurrency > maxConcurrency {
		t.Errorf("Concurrency is out of range: %d", response.concurrency)
	}
	if response.liveLoad < 0 || response.liveLoad > 100 {
		t.Errorf("Live load is out of range: %.2f%%", response.liveLoad)
	}
	if response.last1Min < 0 || response.last1Min > 100 || response.last5Min < 0 || response.last5Min > 100 || response.last15Min < 0 || response.last15Min > 100 {
		t.Errorf("Average loads are out of range: last1Min=%.2f, last5Min=%.2f, last15Min=%.2f",
			response.last1Min, response.last5Min, response.last15Min)
	}
}

// This test simulates the AdjustConcurrency functionality over time
func TestConcurrencyAdjustment(t *testing.T) {
	maxConcurrency := 10
	loadThreshold := 80.0
	go AdjustConcurrency(maxConcurrency, loadThreshold)

	time.Sleep(5 * time.Second) // Wait for some load adjustments

	response := AdjustConcurrency(maxConcurrency, loadThreshold)

	if response.concurrency < 1 || response.concurrency > maxConcurrency {
		t.Errorf("Concurrency is out of range: %d", response.concurrency)
	}
	if response.liveLoad < 0 || response.liveLoad > 100 {
		t.Errorf("Live load is out of range: %.2f%%", response.liveLoad)
	}
	if response.last1Min < 0 || response.last1Min > 100 || response.last5Min < 0 || response.last5Min > 100 || response.last15Min < 0 || response.last15Min > 100 {
		t.Errorf("Average loads are out of range: last1Min=%.2f, last5Min=%.2f, last15Min=%.2f",
			response.last1Min, response.last5Min, response.last15Min)
	}
}

// Test for GoWithCustomLimit function
func TestGoWithCustomLimit(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := GoWithCustomLimit(ctx, "TestFunction", arr, 3, mockFunction)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// Test for GoWithCustomLimit with invalid input
func TestGoWithCustomLimitInvalidInput(t *testing.T) {
	arr := []int{-1, 2, 3, 4, 5} // Include an invalid input

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := GoWithCustomLimit(ctx, "TestFunction", arr, 3, mockFunction)
	if err == nil {
		t.Error("expected an error due to invalid input, got nil")
	}
}

// New test for GoWithDynamicLoadAdjustment function
func TestGoWithDynamicLoadAdjustment(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := Go(ctx, "b", arr, mockFunction)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// Test for GoWithDynamicLoadAdjustment with invalid input
func TestGoWithDynamicLoadAdjustmentInvalidInput(t *testing.T) {
	arr := []int{-1, 2, 3, 4, 5} // Include an invalid input

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := Go(ctx, "b", arr, mockFunction)
	if err == nil {
		t.Error("expected an error due to invalid input, got nil")
	}
}
