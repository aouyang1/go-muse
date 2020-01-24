package muse

import (
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
		if score.PercentScore != expectedScores[i].PercentScore {
			t.Errorf("Expected %d score but got %d score", expectedScores[i].PercentScore, score.PercentScore)
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

	comp := []*Series{
		NewSeries([]float64{0, 0, 0, 0, 2, 4, 6, 6, 4, 2, 0, 0}, NewLabels(LabelMap{"graph": "perfectMatch"})),
		NewSeries([]float64{0, 0, 0, 0, 2, 4, 6, 4, 2, 0, 0, 0}, NewLabels(LabelMap{"graph": "slightlyLower"})),
		NewSeries([]float64{0, 0, 0, 2, 4, 2, 0, 0, 0, 0, 0, 0}, NewLabels(LabelMap{"graph": "evenLower"})),
		NewSeries([]float64{0, 0, 0, 0, 0, 0, 0, 0, 2, 3, 2, 0}, NewLabels(LabelMap{"graph": "evenLowerShiftedAhead"})),
		NewSeries([]float64{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}, NewLabels(LabelMap{"graph": "zeros"})),
	}

	expectedScores := Scores{
		Score{Labels: NewLabels(LabelMap{"graph": "perfectMatch"}), Lag: 0, PercentScore: 100},
		Score{Labels: NewLabels(LabelMap{"graph": "slightlyLower"}), Lag: 0, PercentScore: 93},
		Score{Labels: NewLabels(LabelMap{"graph": "evenLowerShiftedAhead"}), Lag: -3, PercentScore: 75},
		Score{Labels: NewLabels(LabelMap{"graph": "evenLower"}), Lag: 2, PercentScore: 73},
		Score{Labels: NewLabels(LabelMap{"graph": "zeros"}), Lag: 0, PercentScore: 0},
	}

	compGroup := NewGroup("targets")
	if err := compGroup.Add(comp...); err != nil {
		t.Fatalf("%v", err)
	}

	Concurrency = 10
	g, err := New(ref, compGroup, NewResults(10, 20, 0))
	if err != nil {
		t.Fatalf("%v", err)
	}
	g.Run([]string{"graph"})

	scores, _ := g.Results.Fetch()
	compareScores(scores, expectedScores, t)
}

func TestRunMultiDimensional(t *testing.T) {
	ref := NewSeries(
		[]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4},
		NewLabels(LabelMap{"graph": "graph1"}),
	)

	comp := []*Series{
		NewSeries([]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4}, NewLabels(LabelMap{"graph": "graph1", "host": "host1"})),
		NewSeries([]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.2, 0.1}, NewLabels(LabelMap{"graph": "graph1", "host": "host2"})),
		NewSeries([]float64{0.0, 0.0, 0.0, 0.0, 0.2, 0.4, 0.4, 0.8}, NewLabels(LabelMap{"graph": "graph2", "host": "host1"})),
		NewSeries([]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.22, 0.1}, NewLabels(LabelMap{"graph": "graph3", "host": "host1"})),
		NewSeries([]float64{0.0, 0.0, 0.0, 0.0, -0.2, -0.4, 0.0, -0.8}, NewLabels(LabelMap{"graph": "graph4", "host": "host1"})),
		NewSeries([]float64{0.0, 0.0, 0.0, -0.2, -0.4, -0.6, 1.0, 0.0}, NewLabels(LabelMap{"graph": "graph5", "host": "host1"})),
	}

	expectedScores := Scores{
		Score{Labels: NewLabels(LabelMap{"graph": "graph1", "host": "host1"}), Lag: 0, PercentScore: 100},
		Score{Labels: NewLabels(LabelMap{"graph": "graph2", "host": "host1"}), Lag: 0, PercentScore: 98},
		Score{Labels: NewLabels(LabelMap{"graph": "graph4", "host": "host1"}), Lag: 0, PercentScore: 76},
		Score{Labels: NewLabels(LabelMap{"graph": "graph5", "host": "host1"}), Lag: 2, PercentScore: 47},
		Score{Labels: NewLabels(LabelMap{"graph": "graph3", "host": "host1"}), Lag: 7, PercentScore: 21},
	}

	compGroup := NewGroup("targets")
	if err := compGroup.Add(comp...); err != nil {
		t.Fatalf("%v", err)
	}

	Concurrency = 10
	m, err := New(ref, compGroup, NewResults(10, 20, 0))
	if err != nil {
		t.Fatalf("%+v\n", err)
	}
	m.Run([]string{"graph"})

	scores, _ := m.Results.Fetch()
	compareScores(scores, expectedScores, t)
}
func TestRunWithLargerGroup(t *testing.T) {
	ref := NewSeries(
		[]float64{0, 1, 2, 3, 3, 2, 1, 0},
		NewLabels(LabelMap{"graph": "graph1"}),
	)

	comp := []*Series{
		NewSeries([]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 4, 6, 6, 4, 2, 0, 0}, NewLabels(LabelMap{"graph": "longer"})),
	}

	compGroup := NewGroup("targets")
	if err := compGroup.Add(comp...); err != nil {
		t.Fatalf("%v", err)
	}

	_, err := New(ref, compGroup, NewResults(10, 20, 0))
	if err == nil {
		t.Fatalf("Expected error with length mismatch of comparison and reference time series")
	}
}

func BenchmarkMuseRun(b *testing.B) {

	ref := NewSeries(
		[]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4},
		NewLabels(LabelMap{"graph": "graph1"}),
	)

	comp := []*Series{
		NewSeries([]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4}, NewLabels(LabelMap{"graph": "graph1", "host": "host1"})),
		NewSeries([]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.2, 0.1}, NewLabels(LabelMap{"graph": "graph1", "host": "host2"})),
		NewSeries([]float64{0.0, 0.0, 0.0, 0.0, 0.2, 0.4, 0.5, 0.8}, NewLabels(LabelMap{"graph": "graph2", "host": "host1"})),
		NewSeries([]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.22, 0.1}, NewLabels(LabelMap{"graph": "graph3", "host": "host1"})),
		NewSeries([]float64{0.0, 0.0, 0.0, 0.0, -0.2, -0.4, 0.0, -0.8}, NewLabels(LabelMap{"graph": "graph4", "host": "host1"})),
		NewSeries([]float64{0.0, 0.0, 0.0, -0.2, -0.4, -0.6, 1.0, 0.0}, NewLabels(LabelMap{"graph": "graph5", "host": "host1"})),
	}

	compGroup := NewGroup("targets")
	if err := compGroup.Add(comp...); err != nil {
		b.Fatalf("%v", err)
	}

	Concurrency = 10

	for i := 0; i < b.N; i++ {
		g, err := New(ref, compGroup, NewResults(10, 20, 0))
		if err != nil {
			b.Fatalf("%v\n", err)
		}
		g.Run([]string{"graph"})
	}
}

func BenchmarkMuseRunLarge(b *testing.B) {
	n := 480
	ref := NewSeries(siggen.Noise(0.1, n), nil)

	compGroup := NewGroup("targets")

	for i := 0; i < 100; i++ {
		for j := 0; j < 50; j++ {
			if err := compGroup.Add(
				NewSeries(
					siggen.Noise(0.1, n),
					NewLabels(LabelMap{"graph": "graph" + string(i), "host": "host" + string(j)}),
				),
			); err != nil {
				b.Fatalf("%v", err)
			}
		}
	}

	Concurrency = 10

	for i := 0; i < b.N; i++ {
		g, err := New(ref, compGroup, NewResults(10, 20, 0))
		if err != nil {
			b.Fatalf("%+v\n", err)
		}
		g.Run([]string{"host"})
	}
}
