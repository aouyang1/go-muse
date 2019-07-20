package muse

import (
	"reflect"
	"testing"
)

var (
	y = []float64{0.1, 0.2, 0.3}
)

func TestLabelsUID(t *testing.T) {
	data := []struct {
		labels        Labels
		groupByLabels []string
		expectedUID   string
	}{
		{
			Labels{
				"a": "v1",
			},
			[]string{"a"},
			"a:v1",
		},
		{
			Labels{
				"a": "v1",
				"c": "v2",
				"b": "v3",
			},
			nil,
			"a:v1,b:v3,c:v2",
		},
		{
			Labels{
				"a": "v1",
				"c": "v2",
				"b": "v3",
			},
			[]string{"a", "b", "c"},
			"a:v1,b:v3,c:v2",
		},
		{
			Labels{
				"a": "v1",
				"A": "v2",
				"b": "v3",
			},
			[]string{"a", "A", "b"},
			"A:v2,a:v1,b:v3",
		},
		{
			Labels{
				"a": "v1",
				"c": "v2",
				"b": "v3",
			},
			[]string{"a", "c"},
			"a:v1,c:v2",
		},
		{
			Labels{
				"a": "v1",
				"c": "v2",
				"b": "v3",
			},
			[]string{"d"},
			"",
		},
	}

	var uid string
	for _, d := range data {
		uid = d.labels.UID(d.groupByLabels)
		if uid != d.expectedUID {
			t.Fatalf("Expected %s but got %s", d.expectedUID, uid)
		}
	}
}

func TestNewSeries(t *testing.T) {
	data := []struct {
		y                 []float64
		labels            Labels
		expectedLabelKeys []string
	}{
		{y, Labels{"a": "v1"}, []string{"a"}},
		{y, nil, []string{DefaultLabel}},
	}

	var s *Series
	for _, d := range data {
		s = NewSeries(d.y, d.labels)
		if !reflect.DeepEqual(s.Labels().Keys(), d.expectedLabelKeys) {
			t.Fatalf("Expected a %v label keys but got %v", d.expectedLabelKeys, s.Labels().Keys())
		}
	}
}

func TestSeriesUID(t *testing.T) {
	data := []struct {
		labels      Labels
		expectedUID string
	}{
		{
			Labels{
				"a": "v1",
			},
			"a:v1",
		},
		{
			Labels{
				"a": "v1",
				"c": "v2",
				"b": "v3",
			},
			"a:v1,b:v3,c:v2",
		},
	}

	var uid string
	var s *Series
	for _, d := range data {
		s = NewSeries(y, d.labels)
		uid = s.UID()
		if uid != d.expectedUID {
			t.Fatalf("Expected %s but got %s", d.expectedUID, uid)
		}
	}
}

func TestLabels(t *testing.T) {
	data := []struct {
		labels         Labels
		expectedLabels []string
	}{
		{
			Labels{
				"a": "v1",
				"c": "v2",
				"b": "v3",
			},
			[]string{"a", "b", "c"},
		},
		{
			Labels{
				"a": "v1",
				"A": "v2",
				"b": "v3",
			},
			[]string{"A", "a", "b"},
		},
	}

	var labels []string
	var s *Series
	for _, d := range data {
		s = NewSeries(y, d.labels)
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
		labels      Labels
		expectError bool
	}{
		{Labels{"a": "v1"}, false},
		{Labels{"a": "v2"}, false},
		{Labels{"a": "v1", "c": "v2", "b": "v3"}, false},
		{Labels{"a": "v1", "A": "v2", "b": "v3"}, false},
		{Labels{"a": "v1"}, true},
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

	var dl []Labels
	for _, p := range testParams {
		dl = g.indexLabelValues(p.labelNames)
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
		s = NewSeries(y, l)
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
		g.indexLabelValues(p.labels.Keys())
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
		s = NewSeries(y, l)
		if err = g.Add(s); err != nil {
			b.Fatalf("%v", err)
		}
	}

	g.indexLabelValues([]string{"graph"})
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
		s = NewSeries(y, l)
		if err = g.Add(s); err != nil {
			b.Fatalf("%v", err)
		}
	}

	for i := 0; i < b.N; i++ {
		g.indexLabelValues([]string{"graph"})
	}
}
