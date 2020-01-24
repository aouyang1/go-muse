[![Build Status](https://travis-ci.com/aouyang1/go-muse.svg?branch=master)](https://travis-ci.com/aouyang1/go-muse)
[![codecov](https://codecov.io/gh/aouyang1/go-muse/branch/master/graph/badge.svg)](https://codecov.io/gh/aouyang1/go-muse)
[![Go Report Card](https://goreportcard.com/badge/github.com/aouyang1/go-muse)](https://goreportcard.com/report/github.com/aouyang1/go-muse)
[![GoDoc](https://godoc.org/github.com/aouyang1/go-muse?status.svg)](https://godoc.org/github.com/aouyang1/go-muse)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

# go-muse
Golang library for comparing one time series with a group of other labeled time series. This library supports arbitrary slicing and dicing of the labeled time series for a more iterative approach to finding visibly similar timeseries. Comparison between two timeseries is done with z-normalizing both series and cross correlating the two using Fast Fourier Transforms (FFT). This library also support parallelization by setting the `Concurrency` variable in the package.

### Motivation
A common problem in the operations world is finding all the graphs that look like a particular alert or incident. For example a Site Reliability Engineer (SRE) receives an alert which indicates that something is broken. The SRE generally will open up the graph that triggered the alert which is likely one graph of many in a dashboard. Next, the SRE begins scrolling through this dashboard looking for anything that resembles the waveform of the received alert. Once the SRE has filtered down the set of graphs that looks the original alert graph, he/she begins building a story as to why the alert fired and root causing the incident. This whole process can be time consuming depending on the size and complexity of the dashboards. This library aims to provide a first pass filtering of the existing graphs or time series, so that an engineer can focus just on what looks similar.

This library will filter results down to anything that is positively or negatively correlated with the input reference series. You can limit the number of results returned and also specify a score between 0 to 1 with 1 being perfectly correlated. You can also filter down to a number of samples before and after the input reference series to filter our strong matches that are outside your window of interest.

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
$ make all
```

## Quick start
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
		), NewLabels(LabelMap{"graph": "CallTime99Pct", "host": "host1"}),
	)

	// create a comparison group of time series that the reference will query against
	comp := NewGroup("comparison")
	comp.Add(
		ref,
		NewSeries(
			siggen.Add(
				siggen.Rect(1.5, 242, 7, sampleRate, duration),
				siggen.Noise(0.1, int(sampleRate*duration)),
			), NewLabels(LabelMap{"graph": "CallTime99Pct", "host": "host2"}),
		),
		NewSeries(
			siggen.Add(
				siggen.Rect(43, 240, 10, sampleRate, duration),
				siggen.Noise(0.1, int(sampleRate*duration)),
			), NewLabels(LabelMap{"graph": "ErrorRate", "host": "host1"}),
		),
		NewSeries(
			siggen.Add(
				siggen.Line(0, 0.1, int(sampleRate*duration)),
				siggen.Noise(0.1, int(sampleRate*duration)),
			), NewLabels(LabelMap{"graph": "ErrorRate", "host": "host2"}),
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

## Benchmarks
Benchmark name                      | NumReps |    Time/Rep    |   Memory/Rep  |     Alloc/Rep   |
-----------------------------------:|--------:|---------------:|--------------:|----------------:|
BenchmarkMuseRun-4                  |    50000|     37319 ns/op|      9626 B/op|    112 allocs/op| 
BenchmarkMuseRunLarge-4             |       10| 184963980 ns/op| 133454180 B/op|  32001 allocs/op|
BenchmarkFilterByLabelValues-4      |  2000000|       569 ns/op|       496 B/op|      8 allocs/op|
BenchmarkIndexLabelValues-4         |   500000|      2709 ns/op|      2152 B/op|     38 allocs/op|
BenchmarkZPad-4                     | 30000000|      40.8 ns/op|        80 B/op|      1 allocs/op|
BenchmarkZNormalize-4               | 20000000|      63.4 ns/op|         0 B/op|      0 allocs/op|
BenchmarkXCorr-4                    |      300|   5651196 ns/op|   2114464 B/op|      7 allocs/op|
BenchmarkXCorrWithX-4               |      500|   3508246 ns/op|    799391 B/op|      3 allocs/op|

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
