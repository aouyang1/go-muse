package muse

import (
	"fmt"

	"github.com/matrix-profile-foundation/go-matrixprofile/siggen"
)

func Example() {
	sampleRate := 1.0 // once per minute
	duration := 480.0 // minutes

	ref := NewSeries(
		siggen.Add(
			siggen.Rect(1.5, 240, 10, sampleRate, duration),
			siggen.Noise(0.1, int(sampleRate*duration)),
		), Labels{"graph": "CallTime99Pct", "host": "host1"},
	)

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
	m.Run(nil)
	fmt.Println(m.Results.Fetch())

	m.Run([]string{"graph"})
	fmt.Println(m.Results.Fetch())

	m.Run([]string{"host"})
	fmt.Println(m.Results.Fetch())

	// Output:
	// [{map[graph:CallTime99Pct host:host2] -3 82} {map[graph:ErrorRate host:host1] 0 99} {map[graph:CallTime99Pct host:host1] 0 100}] 93.66666666666667
	// [{map[graph:ErrorRate host:host1] 0 99} {map[graph:CallTime99Pct host:host1] 0 100}] 99.5
	// [{map[graph:CallTime99Pct host:host2] -3 82} {map[graph:CallTime99Pct host:host1] 0 100}] 91
}
