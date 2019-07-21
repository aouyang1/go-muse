package muse

import (
	"reflect"
	"testing"
)

var (
	y = []float64{0.1, 0.2, 0.3}
)

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

func TestSeriesLabels(t *testing.T) {
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
