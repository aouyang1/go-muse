package muse

import (
	"testing"
)

func TestRunSimple(t *testing.T) {
	ref := NewSeries(
		[]float64{0, 0, 0, 0, 1, 2, 3, 3, 2, 1, 0, 0},
		map[string]string{"graph": "graph1"},
	)
	compGroup := NewGroup("targets")

	data := []struct {
		y      []float64
		labels Labels
	}{
		{[]float64{0, 0, 0, 0, 2, 4, 6, 6, 4, 2, 0, 0}, Labels{"graph": "perfectMatch"}},
		{[]float64{0, 0, 0, 0, 2, 4, 6, 4, 2, 0, 0, 0}, Labels{"graph": "slightlyLower"}},
		{[]float64{0, 0, 0, 2, 4, 2, 0, 0, 0, 0, 0, 0}, Labels{"graph": "evenLower"}},
	}

	for _, d := range data {
		comp := NewSeries(d.y, d.labels)
		if err := compGroup.Add(comp); err != nil {
			t.Fatalf("%v", err)
		}
	}

	Concurrency = 10
	g := New(ref, compGroup, NewResults(10, 20, 0))
	g.Run([]string{"graph"})

	expectedScores := Scores{
		Score{Labels: map[string]string{"graph": "evenLower"}, Lag: 2, PercentScore: 83},
		Score{Labels: map[string]string{"graph": "slightlyLower"}, Lag: 0, PercentScore: 95},
		Score{Labels: map[string]string{"graph": "perfectMatch"}, Lag: 0, PercentScore: 100},
	}

	scores, _ := g.Results.Fetch()
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

func TestRunMultiDimensional(t *testing.T) {
	ref := NewSeries(
		[]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4},
		map[string]string{"graph": "graph1"},
	)
	compGroup := NewGroup("targets")

	data := []struct {
		y      []float64
		labels Labels
	}{
		{[]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4}, Labels{"graph": "graph1", "host": "host1"}},
		{[]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.2, 0.1}, Labels{"graph": "graph1", "host": "host2"}},
		{[]float64{0.0, 0.0, 0.0, 0.0, 0.2, 0.4, 0.5, 0.8}, Labels{"graph": "graph2", "host": "host1"}},
		{[]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.22, 0.1}, Labels{"graph": "graph3", "host": "host1"}},
		{[]float64{0.0, 0.0, 0.0, 0.0, -0.2, -0.4, 0.0, -0.8}, Labels{"graph": "graph4", "host": "host1"}},
		{[]float64{0.0, 0.0, 0.0, -0.2, -0.4, -0.6, 1.0, 0.0}, Labels{"graph": "graph5", "host": "host1"}},
	}

	for _, d := range data {
		comp := NewSeries(d.y, d.labels)
		if err := compGroup.Add(comp); err != nil {
			t.Fatalf("%v", err)
		}
	}

	Concurrency = 10
	m := New(ref, compGroup, NewResults(10, 20, 0))
	m.Run([]string{"graph"})

	expectedScores := Scores{
		Score{Labels: map[string]string{"graph": "graph3", "host": "host1"}, Lag: 1, PercentScore: 55},
		Score{Labels: map[string]string{"graph": "graph5", "host": "host1"}, Lag: 2, PercentScore: 63},
		Score{Labels: map[string]string{"graph": "graph4", "host": "host1"}, Lag: 0, PercentScore: 80},
		Score{Labels: map[string]string{"graph": "graph2", "host": "host1"}, Lag: 0, PercentScore: 99},
		Score{Labels: map[string]string{"graph": "graph1", "host": "host1"}, Lag: 0, PercentScore: 100},
	}

	scores, _ := m.Results.Fetch()
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

func BenchmarkMuseRun(b *testing.B) {

	ref := NewSeries(
		[]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4},
		map[string]string{"graph": "graph1"},
	)
	compGroup := NewGroup("targets")

	data := []struct {
		y      []float64
		labels Labels
	}{
		{[]float64{0.0, 0.0, 0.0, 0.0, 0.1, 0.2, 0.3, 0.4}, Labels{"graph": "graph1", "host": "host1"}},
		{[]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.2, 0.1}, Labels{"graph": "graph1", "host": "host2"}},
		{[]float64{0.0, 0.0, 0.0, 0.0, 0.2, 0.4, 0.5, 0.8}, Labels{"graph": "graph2", "host": "host1"}},
		{[]float64{0.2, 0.1, 0.2, 0.1, 0.2, 0.1, 0.22, 0.1}, Labels{"graph": "graph3", "host": "host1"}},
		{[]float64{0.0, 0.0, 0.0, 0.0, -0.2, -0.4, 0.0, -0.8}, Labels{"graph": "graph4", "host": "host1"}},
		{[]float64{0.0, 0.0, 0.0, -0.2, -0.4, -0.6, 1.0, 0.0}, Labels{"graph": "graph5", "host": "host1"}},
	}

	for _, d := range data {
		comp := NewSeries(d.y, d.labels)
		if err := compGroup.Add(comp); err != nil {
			b.Fatalf("%v", err)
		}
	}

	Concurrency = 10

	for i := 0; i < b.N; i++ {
		g := New(ref, compGroup, NewResults(10, 20, 0))
		g.Run([]string{"graph"})
	}
}
