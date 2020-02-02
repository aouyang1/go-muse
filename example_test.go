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
		NewSeries(
			siggen.Line(0, 0.1, int(sampleRate*duration)),
			NewLabels(LabelMap{"graph": "ErrorRate", "host": "host3"}),
		),
	)

	maxLag := 15.0   // minutes
	topN := 4        // top 4 grouped series
	threshold := 0.0 // correlation threshold
	m, err := NewBatch(ref, comp, NewResults(int(maxLag/sampleRate), topN, threshold), 2)
	if err != nil {
		panic(err)
	}

	// Rank each individual time series in the comparison group
	m.Run(nil)
	res, _ := m.Results.Fetch()
	fmt.Println("Unique")
	for _, s := range res {
		fmt.Printf("%s, Lag: %d, Score: %.3f\n", s.Labels.ID(s.Labels.Keys()), s.Lag, s.PercentScore)
	}

	// Rank time series grouped by the graph label
	m.Run([]string{"graph"})
	res, _ = m.Results.Fetch()
	fmt.Println("By Graph")
	for _, s := range res {
		fmt.Printf("%s, Lag: %d, Score: %.3f\n", s.Labels.ID(s.Labels.Keys()), s.Lag, s.PercentScore)
	}

	// Rank time series grouped by the host label
	m.Run([]string{"host"})
	res, _ = m.Results.Fetch()
	fmt.Println("By Host")
	for _, s := range res {
		fmt.Printf("%s, Lag: %d, Score: %.3f\n", s.Labels.ID(s.Labels.Keys()), s.Lag, s.PercentScore)
	}

	// Output: Unique
	// graph:CallTime99Pct,host:host1, Lag: 0, Score: 1.000
	// graph:ErrorRate,host:host1, Lag: 0, Score: 0.991
	// graph:CallTime99Pct,host:host2, Lag: -3, Score: 0.822
	// graph:ErrorRate,host:host3, Lag: 0, Score: 0.000
	// By Graph
	// graph:CallTime99Pct,host:host1, Lag: 0, Score: 1.000
	// graph:ErrorRate,host:host1, Lag: 0, Score: 0.991
	// By Host
	// graph:CallTime99Pct,host:host1, Lag: 0, Score: 1.000
	// graph:CallTime99Pct,host:host2, Lag: -3, Score: 0.822
	// graph:ErrorRate,host:host3, Lag: 0, Score: 0.000
}
