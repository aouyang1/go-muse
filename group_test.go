package muse

import "testing"

func TestGroupAdd(t *testing.T) {
	data := []struct {
		labels      *Labels
		expectError bool
	}{
		{NewLabels(LabelMap{"a": "v1"}), false},
		{NewLabels(LabelMap{"a": "v2"}), false},
		{NewLabels(LabelMap{"a": "v1", "c": "v2", "b": "v3"}), false},
		{NewLabels(LabelMap{"a": "v1", "A": "v2", "b": "v3"}), false},
		{NewLabels(LabelMap{"a": "v1"}), true},
	}
	g := NewGroup("test")

	var err error
	var s *Series
	for _, d := range data {
		s = NewSeries(y, d.labels)
		if err = g.Add(s); (err != nil) != d.expectError {
			t.Fatalf("Expected %t error for label %v", d.expectError, d.labels)
		}
	}
}

func TestIndexLabelValues(t *testing.T) {
	g := NewGroup("test")

	labels := []*Labels{
		NewLabels(LabelMap{"graph": "graph1", "host": "host1", "colo": "colo1"}),
		NewLabels(LabelMap{"graph": "graph1", "host": "host1", "colo": "colo2"}),
		NewLabels(LabelMap{"graph": "graph1", "host": "host2", "colo": "colo1"}),
		NewLabels(LabelMap{"graph": "graph1", "host": "host2", "colo": "colo2"}),
		NewLabels(LabelMap{"graph": "graph2", "host": "host1", "colo": "colo1"}),
		NewLabels(LabelMap{"graph": "graph3", "host": "host2", "colo": "colo1"}),
	}

	var s *Series
	var err error
	for _, l := range labels {
		s = NewSeries(y, l)
		if err = g.Add(s); err != nil {
			t.Fatalf("%v", err)
		}
	}

	testParams := []struct {
		labelNames        []string
		expectedNumLabels int
	}{
		{[]string{"graph"}, 3},
		{[]string{"host"}, 2},
		{[]string{"colo"}, 2},
		{[]string{"graph", "host"}, 4},
		{[]string{"host", "colo"}, 4},
		{[]string{"graph", "colo"}, 4},
		{[]string{"graph", "host", "colo"}, 6},
	}

	var dl []*Labels
	for _, p := range testParams {
		dl = g.indexLabelValues(p.labelNames)
		if len(dl) != p.expectedNumLabels {
			t.Fatalf("Expected %d distinct labels grouped by %v, but got %d", p.expectedNumLabels, p.labelNames, len(dl))
		}
	}
}

func TestFilterByLabelValues(t *testing.T) {
	g := NewGroup("test")

	labels := []*Labels{
		NewLabels(LabelMap{"graph": "graph1", "host": "host1", "colo": "colo1"}),
		NewLabels(LabelMap{"graph": "graph1", "host": "host1", "colo": "colo2"}),
		NewLabels(LabelMap{"graph": "graph1", "host": "host2", "colo": "colo1"}),
		NewLabels(LabelMap{"graph": "graph1", "host": "host2", "colo": "colo2"}),
		NewLabels(LabelMap{"graph": "graph2", "host": "host1", "colo": "colo1"}),
		NewLabels(LabelMap{"graph": "graph3", "host": "host2", "colo": "colo1"}),
	}

	var s *Series
	var err error
	for _, l := range labels {
		s = NewSeries(y, l)
		if err = g.Add(s); err != nil {
			t.Fatalf("%v", err)
		}
	}

	testParams := []struct {
		labels            *Labels
		expectedNumSeries int
	}{
		{NewLabels(LabelMap{"graph": "graph1"}), 4},
		{NewLabels(LabelMap{"graph": "graph2"}), 1},
		{NewLabels(LabelMap{"host": "host1"}), 3},
		{NewLabels(LabelMap{"host": "host2"}), 3},
		{NewLabels(LabelMap{"host": "host0"}), 0},
		{NewLabels(LabelMap{"graph": "graph1", "host": "host2"}), 2},
		{NewLabels(LabelMap{"graph": "graph1", "host": "host0"}), 0},
	}

	var series []*Series
	for _, p := range testParams {
		g.indexLabelValues(p.labels.Keys())
		series = g.FilterByLabelValues(p.labels)
		if len(series) != p.expectedNumSeries {
			t.Fatalf("Expected %d series filtered by %v, but got %d", p.expectedNumSeries, p.labels, len(series))
		}
	}
}

func BenchmarkFilterByLabelValues(b *testing.B) {
	g := NewGroup("test")

	labels := []*Labels{
		NewLabels(LabelMap{"graph": "graph1", "host": "host1", "colo": "colo1"}),
		NewLabels(LabelMap{"graph": "graph1", "host": "host1", "colo": "colo2"}),
		NewLabels(LabelMap{"graph": "graph1", "host": "host2", "colo": "colo1"}),
		NewLabels(LabelMap{"graph": "graph1", "host": "host2", "colo": "colo2"}),
		NewLabels(LabelMap{"graph": "graph2", "host": "host1", "colo": "colo1"}),
		NewLabels(LabelMap{"graph": "graph3", "host": "host2", "colo": "colo1"}),
	}

	var s *Series
	var err error
	for _, l := range labels {
		s = NewSeries(y, l)
		if err = g.Add(s); err != nil {
			b.Fatalf("%v", err)
		}
	}

	g.indexLabelValues([]string{"graph"})
	for i := 0; i < b.N; i++ {
		g.FilterByLabelValues(NewLabels(LabelMap{"graph": "graph1"}))
	}
}

func BenchmarkIndexLabelValues(b *testing.B) {
	g := NewGroup("test")

	labels := []*Labels{
		NewLabels(LabelMap{"graph": "graph1", "host": "host1", "colo": "colo1"}),
		NewLabels(LabelMap{"graph": "graph1", "host": "host1", "colo": "colo2"}),
		NewLabels(LabelMap{"graph": "graph1", "host": "host2", "colo": "colo1"}),
		NewLabels(LabelMap{"graph": "graph1", "host": "host2", "colo": "colo2"}),
		NewLabels(LabelMap{"graph": "graph2", "host": "host1", "colo": "colo1"}),
		NewLabels(LabelMap{"graph": "graph3", "host": "host2", "colo": "colo1"}),
	}

	var s *Series
	var err error
	for _, l := range labels {
		s = NewSeries(y, l)
		if err = g.Add(s); err != nil {
			b.Fatalf("%v", err)
		}
	}

	for i := 0; i < b.N; i++ {
		g.indexLabelValues([]string{"graph"})
	}
}
