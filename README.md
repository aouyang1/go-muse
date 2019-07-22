[![Build Status](https://travis-ci.com/aouyang1/go-muse.svg?branch=master)](https://travis-ci.com/aouyang1/go-muse)
[![codecov](https://codecov.io/gh/aouyang1/go-muse/branch/master/graph/badge.svg)](https://codecov.io/gh/aouyang1/go-muse)
[![Go Report Card](https://goreportcard.com/badge/github.com/aouyang1/go-muse)](https://goreportcard.com/report/github.com/aouyang1/go-muse)
[![GoDoc](https://godoc.org/github.com/aouyang1/go-muse?status.svg)](https://godoc.org/github.com/aouyang1/go-muse)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

# go-muse

Golang library for comparing one time series with a group of other labeled time series. This library supports arbitrary slicing and dicing of the labeled time series for a more iterative approach to finding visibly similar timeseries. Comparison between two timeseries is done with z-normalizing both series and cross correlating the two using Fast Fourier Transforms (FFT). This library also support parallelization by setting the `Concurrency` variable in the package.

## Contents
- [Installation](#installation)
- [Quick start](#quick-start)
- [Benchmarks](#benchmarks)
- [Contributing](#contributing)
- [Testing](#testing)
- [Contact](#contact)
- [License](#license)

## Installation
```sh
$ go get -u github.com/aouyang1/go-muse
$ cd $GOPATH/src/github.com/aouyang1/go-muse
$ make setup
```

## Quick start
```sh
$ cat example_test.go
```
```go
package main

import (
	"fmt"

	"github.com/aouyang1/go-muse"
	"github.com/matrix-profile-foundation/go-matrixprofile/siggen"
)

func main() {
	sampleRate := 1.0 // once per minute
	duration := 480.0 // minutes

	// create a reference rectangular time series with an amplitude of 1.5 centered
	// at 240 minutes and a width of 10 minutes
	ref := NewSeries(
		siggen.Add(
			siggen.Rect(1.5, 240, 10, sampleRate, duration),
			siggen.Noise(0.1, int(sampleRate*duration)),
		), Labels{"graph": "CallTime99Pct", "host": "host1"},
	)

	// create a comparison group of time series that the reference will query against
	comp := NewGroup("comparison")
	comp.Add(
		ref,
		NewSeries(
			siggen.Add(
				siggen.Rect(1.5, 242, 7, sampleRate, duration),
				siggen.Noise(0.1, int(sampleRate*duration)),
			), Labels{"graph": "CallTime99Pct", "host": "host2"},
		),
		NewSeries(
			siggen.Add(
				siggen.Rect(43, 240, 10, sampleRate, duration),
				siggen.Noise(0.1, int(sampleRate*duration)),
			), Labels{"graph": "ErrorRate", "host": "host1"},
		),
		NewSeries(
			siggen.Add(
				siggen.Line(0, 0.1, int(sampleRate*duration)),
				siggen.Noise(0.1, int(sampleRate*duration)),
			), Labels{"graph": "ErrorRate", "host": "host2"},
		),
	)

	maxLag := 15.0   // minutes
	topN := 4        // top 4 grouped series
	threshold := 0.5 // correlation threshold
	m := New(ref, comp, NewResults(int(maxLag/sampleRate), topN, threshold))

	// Rank each individual time series in the comparison group
	m.Run(nil)
	fmt.Println(m.Results.Fetch())

	// Rank time series grouped by the graph label
	m.Run([]string{"graph"})
	fmt.Println(m.Results.Fetch())

	// Rank time series grouped by the host label
	m.Run([]string{"host"})
	fmt.Println(m.Results.Fetch())
}
```
```sh
$ go run example_test.go
[{map[graph:CallTime99Pct host:host2] -3 82} {map[graph:ErrorRate host:host1] 0 99} {map[graph:CallTime99Pct host:host1] 0 100}] 93.66666666666667
[{map[graph:ErrorRate host:host1] 0 99} {map[graph:CallTime99Pct host:host1] 0 100}] 99.5
[{map[graph:CallTime99Pct host:host2] -3 82} {map[graph:CallTime99Pct host:host1] 0 100}] 91
```

## Benchmarks
Benchmark name                      | NumReps |    Time/Rep    |   Memory/Rep  |     Alloc/Rep   |
-----------------------------------:|--------:|---------------:|--------------:|----------------:|
BenchmarkMuseRun-4                  |    30000|     40344 ns/op|      9465 B/op|    107 allocs/op| 
BenchmarkMuseRunLarge-4             |        5| 201152104 ns/op| 135877486 B/op|  39965 allocs/op|
BenchmarkFilterByLabelValues-4      |  3000000|       446 ns/op|       128 B/op|      5 allocs/op|
BenchmarkIndexLabelValues-4         |   500000|      2513 ns/op|      1912 B/op|     29 allocs/op|
BenchmarkZPad-4                     | 30000000|      40.8 ns/op|        80 B/op|      1 allocs/op|
BenchmarkZNormalize-4               | 20000000|      63.4 ns/op|         0 B/op|      0 allocs/op|
BenchmarkXCorr-4                    |      300|   5651196 ns/op|   2114464 B/op|      7 allocs/op|
BenchmarkXCorrNormalize-4           |      300|   5863611 ns/op|   2114465 B/op|      7 allocs/op|
BenchmarkXCorrWithX-4               |      500|   3508246 ns/op|    799391 B/op|      3 allocs/op|
BenchmarkXCorrWithXNormalize-4      |      500|   3690897 ns/op|    799391 B/op|      3 allocs/op|

Ran on a 2018 MacBookAir on Jul 21, 2019
```sh
    Processor: 1.6 GHz Intel Core i5
       Memory: 8GB 2133 MHz LPDDR3
           OS: macOS Mojave v10.14.2
 Logical CPUs: 4
Physical CPUs: 2
```
```sh
$ make bench
```

## Contributing
* Fork the repository
* Create a new branch (feature_\* or bug_\*)for the new feature or bug fix
* Run tests
* Commit your changes
* Push code and open a new pull request

## Testing
Run all tests including benchmarks
```sh
$ make all
```
Just run benchmarks
```sh
$ make bench
```
Just run tests
```sh
$ make test
```

## Contact
* Austin Ouyang (aouyang1@gmail.com)

## License
The MIT License (MIT). See [LICENSE](https://github.com/aouyang1/go-muse/blob/master/LICENSE) for more details.

Copyright (c) 2018 Austin Ouyang