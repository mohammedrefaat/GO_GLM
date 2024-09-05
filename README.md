

# GO_GLM 🌀

GO_GLM is an advanced Go package built to dynamically manage concurrency with load balancing based on real-time CPU and memory usage. It’s perfect for applications needing efficient parallel processing and robust observability through **Prometheus** metrics.

## Key Features 🌟

- **Dynamic Goroutine Management:** Automatically adjust the number of goroutines based on system load (CPU and memory).
- **Load Balancing:** Ensures efficient execution by throttling or boosting the number of concurrent processes to optimize resource usage.
- **Error Handling:** Implements robust handling of panics and context cancellations to avoid resource leaks or crashes.
- **Prometheus Integration:** Exposes a set of metrics for real-time monitoring and performance analysis.
- **Flexible Use:** Easily adaptable to any processing logic requiring parallel execution.

## Why GO_GLM? 🤔

- **Performance:** Prevents overloading the system by monitoring resource usage and adjusting dynamically.
- **Scalability:** Helps scale applications by handling an extensive number of parallel processes without manually managing goroutines.
- **Observability:** Built-in Prometheus metrics make it easy to monitor and optimize application performance.

## Installation ⚙️

To install **GO_GLM**, use:

```bash
go get github.com/mohammedrefaat/GO_GLM
```

## How to Use 🚀

### Import the Package:

```go
import (
    "context"
    "github.com/mohammedrefaat/GO_GLM"
)
```

### Define Your Processing Function:

```go
func processItem(ctx context.Context, item YourType) error {
    // Your processing logic
    return nil
}
```

### Execute with GO_GLM:

```go
func main() {
    ctx := context.Background()
    items := []YourType{ /* your data */ }

    err := GO_GLM.Go(ctx, "processItem", items, processItem)
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

### Monitor Prometheus Metrics:

Expose metrics via HTTP for tracking:

```go
http.Handle("/metrics", promhttp.Handler())
log.Fatal(http.ListenAndServe(":8080", nil))
```

## Prometheus Metrics 📊

GO_GLM provides rich metrics for in-depth monitoring:

| Metric                             | Description                                                 |
|------------------------------------|-------------------------------------------------------------|
| `goroutines_in_progress`           | Number of active goroutines.                                |
| `errors_total`                     | Total number of errors during execution.                    |
| `execution_duration_seconds`       | Duration of function executions.                            |
| `semaphore_wait_duration_seconds`  | Time spent waiting for semaphore resources.                 |
| `successful_executions_total`      | Total number of successful executions.                      |
| `goroutine_queue_length`           | Length of the queue for goroutines awaiting execution.       |
| `total_processing_time_seconds`    | Total time spent on processing all items.                   |

These metrics provide insight into system performance, concurrency levels, and the efficiency of task processing.

## Example Use Case 🛠️

Let’s say you’re processing a large batch of data where the processing time can vary significantly. GO_GLM dynamically adjusts how many items are processed in parallel to prevent overloading your system, ensuring high throughput and stability.

## Contributing 🛠️

Contributions are more than welcome! To get involved:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Commit your changes.
4. Submit a pull request.

Check out the [contributing guidelines](CONTRIBUTING.md) for more information.

## Acknowledgments 🙌

This project wouldn't have been possible without the incredible work of various open-source projects and contributors. Special thanks to:

- **[gopsutil](https://github.com/shirou/gopsutil)** for providing system information, such as CPU and memory stats, which allow GO_GLM to adjust dynamically.
- **[Prometheus Go client](https://github.com/prometheus/client_golang)** for the invaluable monitoring and metrics.
- All contributors who helped improve the codebase and documentation.

## License 📜

This project is licensed under the [MIT License](LICENSE), allowing open use and contribution.

