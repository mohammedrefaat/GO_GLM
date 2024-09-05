
# GO_GLM

GO_GLM is a Go package designed for managing concurrency in Go applications with dynamic load adjustment based on CPU and memory usage. It integrates Prometheus for performance metrics, allowing for monitoring and adjustment of goroutine execution based on system load.

## Features

- Dynamic adjustment of goroutines based on CPU and memory load.
- Prometheus metrics integration to monitor:
  - Number of goroutines in progress
  - Total errors encountered
  - Execution duration
  - Semaphore wait duration
  - Total processing time
- Graceful handling of context cancellation and panic recovery.

## Installation

To use GO_GLM in your project, run:

```bash
go get github.com/mohammedrefaat/GO_GLM
```


## Usage

Here’s a basic example of how to use the `GO_GLM` package:

### 1. Import the package

```go
import (
    "context"
    "fmt"
    "github.com/mohammedrefaat/GO_GLM"
)
```

### 2. Define your function

Define a function that will be executed by the goroutines:

```go
func yourFunction(ctx context.Context, item YourType) error {
    // Your processing logic here
    return nil
}
```

### 3. Execute with dynamic concurrency

Use the `Go` function to run your tasks dynamically:

```go
func main() {
    ctx := context.Background()
    items := []YourType{ /* your data here */ }

    err := GO_GLM.Go(ctx, "yourFunction", items, yourFunction)
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

### 4. Monitor metrics

Prometheus metrics are registered in the package, and you can expose them using an HTTP server for monitoring.

```go
http.Handle("/metrics", promhttp.Handler())
log.Fatal(http.ListenAndServe(":8080", nil))
```

## Metrics

The following metrics are available:

- **goroutines_in_progress**: Current number of goroutines in progress.
- **errors_total**: Total number of errors encountered during execution.
- **execution_duration_seconds**: Duration of function executions in seconds.
- **semaphore_wait_duration_seconds**: Time each goroutine waits for a semaphore.
- **successful_executions_total**: Count of successfully executed goroutines.
- **goroutine_queue_length**: Current length of the goroutine queue.
- **total_processing_time_seconds**: Total cumulative processing time of all goroutines.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss improvements or changes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Prometheus](https://prometheus.io/) for monitoring metrics.
- [gopsutil](https://github.com/shirou/gopsutil) for gathering system metrics.
