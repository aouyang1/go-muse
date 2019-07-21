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
		switch {
		case score.Lag != expectedScores[i].Lag:
			t.Fatalf("Expected %d lag but got %d lag", expectedScores[i].Lag, score.Lag)
		case score.PercentScore != expectedScores[i].PercentScore:
			t.Fatalf("Expected %d score but got %d score", expectedScores[i].PercentScore, score.PercentScore)
		case len(score.Labels) != len(expectedScores[i].Labels):
			t.Fatalf("Expected %d labels but got %d labels", len(expectedScores[i].Labels), len(score.Labels))
		}
		for j, k := range score.Labels.Keys() {
			switch {
			case k != expectedScores[i].Labels.Keys()[j]:
				t.Fatalf("Expected %s label but got %s label", expectedScores[i].Labels.Keys()[j], k)
			case score.Labels[k] != expectedScores[i].Labels[expectedScores[i].Labels.Keys()[j]]:
				t.Fatalf("Expected %s value but got %s value", expectedScores[i].Labels[expectedScores[i].Labels.Keys()[j]], score.Labels[k])
			}
		}
	}
}

func TestRunSimple(t *testing.T) {
	ref := NewSeries(
		[]float64{0, 0, 0, 0, 1, 2, 3, 3, 2, 1, 0, 0},
		map[string]string{"graph": "graph1"},
	)

	comp := []*Series{
		NewSeries([]float64{0, 0, 0, 0, 2, 4, 6, 6, 4, 2, 0, 0}, Labels{"graph": "perfectMatch"}),
		NewSeries([]float64{0, 0, 0, 0, 2, 4, 6, 4, 2, 0, 0, 0}, Labels{"graph": "slightlyLower"}),
		NewSeries([]float64{0, 0, 0, 2, 4, 2, 0, 0, 0, 0, 0, 0}, Labels{"graph": "evenLower"}),
	}

	expectedScores := Scores{
		Score{Labels: map[string]string{"graph": "evenLower"}, Lag: 2, PercentScore: 83},
		Score{Labels: map[string]string{"graph": "slightlyLower"}, Lag: 0, PercentScore: 95},
		Score{Labels: map[string]string{"graph": "perfectMatch"}, Lag: 0, PercentScore: 100},
	}

	compGroup := NewGroup("targets")
	if err := compGroup.Add(comp...); err != nil {
		t.Fatalf("%v", err)
	}

	Concurrency = 10
	g := New(ref, compGroup, NewResults(10, 20, 0))
	g.Run([]string{"graph"})

	scores, _ := g.Results.Fetch()
	compareScores(scores, expectedScores, t)
}

func TestRunMultiDimensional(t *testing.T) {
	ref := NewSeries(
		[]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4},
		map[string]string{"graph": "graph1"},
	)

	comp := []*Series{
		NewSeries([]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4}, Labels{"graph": "graph1", "host": "host1"}),
		NewSeries([]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.2, 0.1}, Labels{"graph": "graph1", "host": "host2"}),
		NewSeries([]float64{0.0, 0.0, 0.0, 0.0, 0.2, 0.4, 0.5, 0.8}, Labels{"graph": "graph2", "host": "host1"}),
		NewSeries([]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.22, 0.1}, Labels{"graph": "graph3", "host": "host1"}),
		NewSeries([]float64{0.0, 0.0, 0.0, 0.0, -0.2, -0.4, 0.0, -0.8}, Labels{"graph": "graph4", "host": "host1"}),
		NewSeries([]float64{0.0, 0.0, 0.0, -0.2, -0.4, -0.6, 1.0, 0.0}, Labels{"graph": "graph5", "host": "host1"}),
	}

	expectedScores := Scores{
		Score{Labels: map[string]string{"graph": "graph3", "host": "host1"}, Lag: 1, PercentScore: 55},
		Score{Labels: map[string]string{"graph": "graph5", "host": "host1"}, Lag: 2, PercentScore: 63},
		Score{Labels: map[string]string{"graph": "graph4", "host": "host1"}, Lag: 0, PercentScore: 80},
		Score{Labels: map[string]string{"graph": "graph2", "host": "host1"}, Lag: 0, PercentScore: 99},
		Score{Labels: map[string]string{"graph": "graph1", "host": "host1"}, Lag: 0, PercentScore: 100},
	}

	compGroup := NewGroup("targets")
	if err := compGroup.Add(comp...); err != nil {
		t.Fatalf("%v", err)
	}

	Concurrency = 10
	m := New(ref, compGroup, NewResults(10, 20, 0))
	m.Run([]string{"graph"})

	scores, _ := m.Results.Fetch()
	compareScores(scores, expectedScores, t)
}
func TestRunWithLargerGroup(t *testing.T) {
	ref := NewSeries(
		[]float64{0, 1, 2, 3, 3, 2, 1, 0},
		map[string]string{"graph": "graph1"},
	)

	comp := []*Series{
		NewSeries([]float64{0, 0, 0, 0, 0, 0, 0, 0, 2, 4, 6, 6, 4, 2, 0, 0}, Labels{"graph": "perfectMatch"}),
		NewSeries([]float64{0, 0, 0, 0, 0, 0, 0, 0, 2, 4, 6, 4, 2, 0, 0, 0}, Labels{"graph": "slightlyLower"}),
		NewSeries([]float64{0, 0, 0, 0, 0, 0, 0, 2, 4, 2, 0, 0, 0, 0, 0, 0}, Labels{"graph": "evenLower"}),
	}

	expectedScores := Scores{
		Score{Labels: map[string]string{"graph": "evenLower"}, Lag: -4, PercentScore: 82},
		Score{Labels: map[string]string{"graph": "slightlyLower"}, Lag: -7, PercentScore: 93},
		Score{Labels: map[string]string{"graph": "perfectMatch"}, Lag: -7, PercentScore: 100},
	}

	compGroup := NewGroup("targets")
	if err := compGroup.Add(comp...); err != nil {
		t.Fatalf("%v", err)
	}

	Concurrency = 10
	g := New(ref, compGroup, NewResults(10, 20, 0))
	g.Run([]string{"graph"})

	scores, _ := g.Results.Fetch()
	compareScores(scores, expectedScores, t)
}

func BenchmarkMuseRun(b *testing.B) {

	ref := NewSeries(
		[]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4},
		map[string]string{"graph": "graph1"},
	)

	comp := []*Series{
		NewSeries([]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4}, Labels{"graph": "graph1", "host": "host1"}),
		NewSeries([]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.2, 0.1}, Labels{"graph": "graph1", "host": "host2"}),
		NewSeries([]float64{0.0, 0.0, 0.0, 0.0, 0.2, 0.4, 0.5, 0.8}, Labels{"graph": "graph2", "host": "host1"}),
		NewSeries([]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.22, 0.1}, Labels{"graph": "graph3", "host": "host1"}),
		NewSeries([]float64{0.0, 0.0, 0.0, 0.0, -0.2, -0.4, 0.0, -0.8}, Labels{"graph": "graph4", "host": "host1"}),
		NewSeries([]float64{0.0, 0.0, 0.0, -0.2, -0.4, -0.6, 1.0, 0.0}, Labels{"graph": "graph5", "host": "host1"}),
	}

	compGroup := NewGroup("targets")
	if err := compGroup.Add(comp...); err != nil {
		b.Fatalf("%v", err)
	}

	Concurrency = 10

	for i := 0; i < b.N; i++ {
		g := New(ref, compGroup, NewResults(10, 20, 0))
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
					Labels{"graph": "graph" + string(i), "host": "host" + string(j)},
				),
			); err != nil {
				b.Fatalf("%v", err)
			}
		}
	}

	Concurrency = 10

	for i := 0; i < b.N; i++ {
		g := New(ref, compGroup, NewResults(10, 20, 0))
		g.Run([]string{"host"})
	}
}
