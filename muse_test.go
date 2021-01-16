package muse

import (
	"math"
	"sync"
	"testing"

	"github.com/matrix-profile-foundation/go-matrixprofile/siggen"
)

func compareScores(scores, expectedScores Scores, t *testing.T) {
	if len(scores) != len(expectedScores) {
		t.Fatalf("Got %d scores, but expected %d scores", len(scores), len(expectedScores))
	}
	for i, score := range scores {
		if score.Lag != expectedScores[i].Lag {
			t.Errorf("Expected %d lag but got %d lag", expectedScores[i].Lag, score.Lag)
			continue
		}
		if math.Abs(score.PercentScore-expectedScores[i].PercentScore) > 1e-3 {
			t.Errorf("Expected %.3f score but got %.3f score", expectedScores[i].PercentScore, score.PercentScore)
			continue
		}
		if score.Labels.Len() != expectedScores[i].Labels.Len() {
			t.Errorf("Expected %d labels but got %d labels", expectedScores[i].Labels.Len(), score.Labels.Len())
			continue
		}
		for j, k := range score.Labels.Keys() {
			if k != expectedScores[i].Labels.Keys()[j] {
				t.Fatalf("Expected %s label but got %s label", expectedScores[i].Labels.Keys()[j], k)
			}
			v, _ := score.Labels.Get(k)
			expectedV, _ := expectedScores[i].Labels.Get(expectedScores[i].Labels.Keys()[j])
			if v != expectedV {
				t.Fatalf("Expected %s value but got %s value", expectedV, v)
			}
		}
	}
}

func TestRunSimple(t *testing.T) {
	ref := NewSeries(
		[]float64{0, 0, 0, 0, 1, 2, 3, 3, 2, 1, 0, 0},
		NewLabels(LabelMap{"graph": "graph1"}),
	)

	comp := [][]*Series{
		[]*Series{NewSeries([]float64{0, 0, 0, 0, 2, 4, 6, 6, 4, 2, 0, 0}, NewLabels(LabelMap{"graph": "perfectMatch"}))},
		[]*Series{NewSeries([]float64{0, 0, 0, 0, 2, 4, 6, 4, 2, 0, 0, 0}, NewLabels(LabelMap{"graph": "slightlyLower"}))},
		[]*Series{NewSeries([]float64{0, 0, 0, 2, 4, 2, 0, 0, 0, 0, 0, 0}, NewLabels(LabelMap{"graph": "evenLower"}))},
		[]*Series{NewSeries([]float64{0, 0, 0, 0, 0, 0, 0, 0, -2, -3, -2, 0}, NewLabels(LabelMap{"graph": "evenLowerShiftedAhead"}))},
		[]*Series{NewSeries([]float64{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}, NewLabels(LabelMap{"graph": "zeros"}))},
	}

	expectedScores := Scores{
		Score{Labels: NewLabels(LabelMap{"graph": "perfectMatch"}), Lag: 0, PercentScore: 1.000},
		Score{Labels: NewLabels(LabelMap{"graph": "slightlyLower"}), Lag: 0, PercentScore: 0.929},
		Score{Labels: NewLabels(LabelMap{"graph": "evenLowerShiftedAhead"}), Lag: -3, PercentScore: -0.754},
		Score{Labels: NewLabels(LabelMap{"graph": "evenLower"}), Lag: 2, PercentScore: 0.733},
		Score{Labels: NewLabels(LabelMap{"graph": "zeros"}), Lag: 0, PercentScore: 0},
	}

	g, err := New(ref, NewResults(10, 20, 0))
	if err != nil {
		t.Fatalf("%v", err)
	}
	for _, c := range comp {
		g.Run(c)
	}

	scores, _ := g.Results.Fetch()
	compareScores(scores, expectedScores, t)
}

func TestRunNoInput(t *testing.T) {
	ref := NewSeries(
		[]float64{0, 0, 0, 0, 1, 2, 3, 3, 2, 1, 0, 0},
		NewLabels(LabelMap{"graph": "graph1"}),
	)

	var comp []*Series

	expectedScores := Scores{}

	g, err := New(ref, NewResults(10, 20, 0))
	if err != nil {
		t.Fatalf("%v", err)
	}
	if err := g.Run(comp); err != nil {
		t.Fatalf("%v", err)
	}

	scores, _ := g.Results.Fetch()
	compareScores(scores, expectedScores, t)
}

func BenchmarkMuseRun(b *testing.B) {

	ref := NewSeries(
		[]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4},
		NewLabels(LabelMap{"graph": "graph1"}),
	)

	comp := [][]*Series{
		[]*Series{
			NewSeries([]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4}, NewLabels(LabelMap{"graph": "graph1", "host": "host1"})),
			NewSeries([]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.2, 0.1}, NewLabels(LabelMap{"graph": "graph1", "host": "host2"})),
		},
		[]*Series{
			NewSeries([]float64{0.0, 0.0, 0.0, 0.0, 0.2, 0.4, 0.5, 0.8}, NewLabels(LabelMap{"graph": "graph2", "host": "host1"})),
		},
		[]*Series{
			NewSeries([]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.22, 0.1}, NewLabels(LabelMap{"graph": "graph3", "host": "host1"})),
		},
		[]*Series{
			NewSeries([]float64{0.0, 0.0, 0.0, 0.0, -0.2, -0.4, 0.0, -0.8}, NewLabels(LabelMap{"graph": "graph4", "host": "host1"})),
		},
		[]*Series{
			NewSeries([]float64{0.0, 0.0, 0.0, -0.2, -0.4, -0.6, 1.0, 0.0}, NewLabels(LabelMap{"graph": "graph5", "host": "host1"})),
		},
	}

	g, err := New(ref, NewResults(10, 20, 0))
	if err != nil {
		b.Fatalf("%v\n", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range comp {
			g.Run(c)
		}
	}
}

func BenchmarkMuseRunLarge(b *testing.B) {
	n := 480
	ref := NewSeries(siggen.Noise(0.1, n), nil)
	numGraphs := 100
	numHosts := 50

	comp := make([][]*Series, numGraphs)
	for i := 0; i < numGraphs; i++ {
		comp[i] = make([]*Series, numHosts)
		for j := 0; j < numHosts; j++ {
			comp[i][j] = NewSeries(
				siggen.Noise(0.1, n),
				NewLabels(LabelMap{"graph": "graph" + string(i), "host": "host" + string(j)}),
			)
		}
	}

	g, err := New(ref, NewResults(10, 20, 0))
	if err != nil {
		b.Fatalf("%+v\n", err)
	}
	var wg sync.WaitGroup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(len(comp))
		for _, c := range comp {
			go func(cc []*Series) {
				defer wg.Done()
				g.Run(cc)
			}(c)
		}
		wg.Wait()
	}
}
