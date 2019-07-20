package muse

import (
	"fmt"

	"github.com/matrix-profile-foundation/go-matrixprofile/siggen"
)

func Example() {
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
