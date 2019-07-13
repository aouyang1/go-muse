package muse

import (
	"testing"
)

var (
	y = []float64{0.1, 0.2, 0.3}
)

func TestNewSeries(t *testing.T) {
	data := []struct {
		y             []float64
		labels        map[string]string
		expectedError bool
	}{
		{y, map[string]string{"a": "v1"}, false},
		{y, nil, true},
	}

	var err error
	for _, d := range data {
		if _, err = NewSeries(d.y, d.labels); (err != nil) != d.expectedError {
			t.Fatalf("Expected a %t error for %v", d.expectedError, d)
		}
	}
}

func TestUID(t *testing.T) {
	data := []struct {
		labels      map[string]string
		expectedUID string
	}{
		{
			map[string]string{
				"a": "v1",
			},
			"a:v1",
		},
		{
			map[string]string{
				"a": "v1",
				"c": "v2",
				"b": "v3",
			},
			"a:v1,b:v3,c:v2",
		},
		{
			map[string]string{
				"a": "v1",
				"A": "v2",
				"b": "v3",
			},
			"A:v2,a:v1,b:v3",
		},
	}

	var uid string
	var s *Series
	var err error
	for _, d := range data {
		s, err = NewSeries(y, d.labels)
		if err != nil {
			t.Fatalf("Failed to create timeseries, %v", err)
		}
		uid = s.UID()
		if uid != d.expectedUID {
			t.Fatalf("Expected %s but got %s", d.expectedUID, uid)
		}
	}

}

func TestLabels(t *testing.T) {
	data := []struct {
		labels         map[string]string
		expectedLabels []string
	}{
		{
			map[string]string{
				"a": "v1",
				"c": "v2",
				"b": "v3",
			},
			[]string{"a", "b", "c"},
		},
		{
			map[string]string{
				"a": "v1",
				"A": "v2",
				"b": "v3",
			},
			[]string{"A", "a", "b"},
		},
	}

	var labels []string
	var s *Series
	var err error
	for _, d := range data {
		s, err = NewSeries(y, d.labels)
		if err != nil {
			t.Fatalf("Failed to create timeseries, %v", err)
		}
		labels = s.labels.Keys()
		if len(labels) != len(d.expectedLabels) {
			t.Fatalf("Expected %d labels but got %d", len(d.expectedLabels), len(labels))
		}
		for i, v := range labels {
			if v != d.expectedLabels[i] {
				t.Fatalf("Expected %s in index %d, but got %s", d.expectedLabels[i], i, v)
			}
		}
	}
}

func TestGroupAdd(t *testing.T) {
	data := []struct {
		labels      map[string]string
		expectError bool
	}{
		{map[string]string{"a": "v1"}, false},
		{map[string]string{"a": "v2"}, false},
		{map[string]string{"a": "v1", "c": "v2", "b": "v3"}, false},
		{map[string]string{"a": "v1", "A": "v2", "b": "v3"}, false},
		{map[string]string{"a": "v1"}, true},
	}
	g := NewGroup("test")

	var err error
	var s *Series
	for _, d := range data {
		s, err = NewSeries(y, d.labels)
		if err != nil {
			t.Fatalf("%v", err)
		}
		if err = g.Add(s); (err != nil) != d.expectError {
			t.Fatalf("Expected %t error for label %v", d.expectError, d.labels)
		}
	}
}

func TestIndexLabelValues(t *testing.T) {
	g := NewGroup("test")

	labels := []Labels{
		Labels{"graph": "graph1", "host": "host1", "colo": "colo1"},
		Labels{"graph": "graph1", "host": "host1", "colo": "colo2"},
		Labels{"graph": "graph1", "host": "host2", "colo": "colo1"},
		Labels{"graph": "graph1", "host": "host2", "colo": "colo2"},
		Labels{"graph": "graph2", "host": "host1", "colo": "colo1"},
		Labels{"graph": "graph3", "host": "host2", "colo": "colo1"},
	}

	var s *Series
	var err error
	for _, l := range labels {
		s, err = NewSeries(y, l)
		if err != nil {
			t.Fatalf("%v", err)
		}
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

	var dl []Labels
	for _, p := range testParams {
		dl = g.IndexLabelValues(p.labelNames)
		if len(dl) != p.expectedNumLabels {
			t.Fatalf("Expected %d distinct labels grouped by %v, but got %d", p.expectedNumLabels, p.labelNames, len(dl))
		}
	}
}

func TestFilterByLabelValues(t *testing.T) {
	g := NewGroup("test")

	labels := []Labels{
		Labels{"graph": "graph1", "host": "host1", "colo": "colo1"},
		Labels{"graph": "graph1", "host": "host1", "colo": "colo2"},
		Labels{"graph": "graph1", "host": "host2", "colo": "colo1"},
		Labels{"graph": "graph1", "host": "host2", "colo": "colo2"},
		Labels{"graph": "graph2", "host": "host1", "colo": "colo1"},
		Labels{"graph": "graph3", "host": "host2", "colo": "colo1"},
	}

	var s *Series
	var err error
	for _, l := range labels {
		s, err = NewSeries(y, l)
		if err != nil {
			t.Fatalf("%v", err)
		}
		if err = g.Add(s); err != nil {
			t.Fatalf("%v", err)
		}
	}

	testParams := []struct {
		labels            Labels
		expectedNumSeries int
	}{
		{Labels{"graph": "graph1"}, 4},
		{Labels{"graph": "graph2"}, 1},
		{Labels{"host": "host1"}, 3},
		{Labels{"host": "host2"}, 3},
		{Labels{"host": "host0"}, 0},
		{Labels{"graph": "graph1", "host": "host2"}, 2},
		{Labels{"graph": "graph1", "host": "host0"}, 0},
	}

	var series []*Series
	for _, p := range testParams {
		g.IndexLabelValues(p.labels.Keys())
		series = g.FilterByLabelValues(p.labels)
		if len(series) != p.expectedNumSeries {
			t.Fatalf("Expected %d series filtered by %v, but got %d", p.expectedNumSeries, p.labels, len(series))
		}
	}
}

func BenchmarkFilterByLabelValues(b *testing.B) {
	g := NewGroup("test")

	labels := []Labels{
		Labels{"graph": "graph1", "host": "host1", "colo": "colo1"},
		Labels{"graph": "graph1", "host": "host1", "colo": "colo2"},
		Labels{"graph": "graph1", "host": "host2", "colo": "colo1"},
		Labels{"graph": "graph1", "host": "host2", "colo": "colo2"},
		Labels{"graph": "graph2", "host": "host1", "colo": "colo1"},
		Labels{"graph": "graph3", "host": "host2", "colo": "colo1"},
	}

	var s *Series
	var err error
	for _, l := range labels {
		s, err = NewSeries(y, l)
		if err != nil {
			b.Fatalf("%v", err)
		}
		if err = g.Add(s); err != nil {
			b.Fatalf("%v", err)
		}
	}

	g.IndexLabelValues([]string{"graph"})
	for i := 0; i < b.N; i++ {
		g.FilterByLabelValues(Labels{"graph": "graph1"})
	}
}

func BenchmarkIndexLabelValues(b *testing.B) {
	g := NewGroup("test")

	labels := []Labels{
		Labels{"graph": "graph1", "host": "host1", "colo": "colo1"},
		Labels{"graph": "graph1", "host": "host1", "colo": "colo2"},
		Labels{"graph": "graph1", "host": "host2", "colo": "colo1"},
		Labels{"graph": "graph1", "host": "host2", "colo": "colo2"},
		Labels{"graph": "graph2", "host": "host1", "colo": "colo1"},
		Labels{"graph": "graph3", "host": "host2", "colo": "colo1"},
	}

	var s *Series
	var err error
	for _, l := range labels {
		s, err = NewSeries(y, l)
		if err != nil {
			b.Fatalf("%v", err)
		}
		if err = g.Add(s); err != nil {
			b.Fatalf("%v", err)
		}
	}

	for i := 0; i < b.N; i++ {
		g.IndexLabelValues([]string{"graph"})
	}
}
