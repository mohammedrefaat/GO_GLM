# Load Manager

Load Manager is a Go package that dynamically adjusts concurrency based on the system's CPU and memory load. It calculates combined load metrics and enables effective management of concurrency based on predefined load thresholds.

## Features

- **Load Calculation**: Calculates the combined CPU and memory load.
- **Dynamic Concurrency Adjustment**: Automatically adjusts concurrency based on live load metrics.
- **Load History**: Maintains load history for the last 1, 5, and 15 minutes.
- **Customizable Thresholds**: Allows configuration of maximum concurrency and load thresholds to optimize performance.

## Installation

To install the `loadmanager` package, run the following command:

```bash
go get github.com/mohammedrefaat/loadmanager
```

## Usage

### Example

Here’s an example of how to use the `loadmanager` package in a Go application:

```go
package main

import (
    "fmt"
    "time"

    "github.com/mohammedrefaat/loadmanager"
)

func main() {
    maxConcurrency := 10
    loadThreshold := 80.0

    for {
        concurrency, liveLoad, last1Min, last5Min, last15Min := loadmanager.AdjustConcurrency(maxConcurrency, loadThreshold)
        fmt.Printf("Live Load: %.2f%%, Last 1-Min Load: %.2f%%, Last 5-Min Load: %.2f%%, Last 15-Min Load: %.2f%%, Concurrency: %d
",
            liveLoad, last1Min, last5Min, last15Min, concurrency)

        time.Sleep(1 * time.Second) // Adjust as needed
    }
}
```

### Functions

- **CalculateLoad() float64**
  - **Description**: Calculates the combined CPU and memory load as a percentage.
  - **Returns**: A float64 value representing the combined load.

- **AdjustConcurrency(maxConcurrency int, loadThreshold float64) (int, float64, float64, float64, float64)**
  - **Description**: Adjusts the concurrency based on the live load and returns the current concurrency along with calculated loads for the last 1, 5, and 15 minutes.
  - **Parameters**:
    - `maxConcurrency`: The maximum allowable concurrency.
    - `loadThreshold`: The load threshold for adjusting concurrency.
  - **Returns**: 
    - An integer representing the current concurrency.
    - Float64 values representing live load, last 1-minute load, last 5-minute load, and last 15-minute load.

## Testing

To run the tests for the `loadmanager` package, use the following command:

```bash
go test ./loadmanager
```

### Test Coverage

The package includes tests for load calculation and concurrency adjustment. You can extend the tests to cover additional scenarios as needed. 

## Contributing

Contributions are welcome! If you have suggestions or improvements, please feel free to open an issue or submit a pull request. When contributing, please adhere to the following guidelines:

1. **Fork the repository**: Create a fork of the repository.
2. **Create a feature branch**: Develop your changes in a new branch.
3. **Commit your changes**: Ensure your commits are clear and concise.
4. **Push your changes**: Push your changes back to your fork.
5. **Open a pull request**: Submit a pull request detailing your changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Acknowledgements

- [gopsutil](https://github.com/shirou/gopsutil): A Go library for retrieving system information, including CPU and memory statistics.
- **Contributors**: Thank you to everyone who has contributed to this project!

## Contact

For questions or inquiries, please reach out to [your email or contact method].
