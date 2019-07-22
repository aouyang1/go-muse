package muse

import "testing"

func TestLabelsUID(t *testing.T) {
	data := []struct {
		labels        *Labels
		groupByLabels []string
		expectedUID   string
	}{
		{
			NewLabels(LabelMap{
				"a": "v1",
			}),
			[]string{"a"},
			"a:v1",
		},
		{
			NewLabels(LabelMap{
				"a": "v1",
				"c": "v2",
				"b": "v3",
			}),
			nil,
			"a:v1,b:v3,c:v2",
		},
		{
			NewLabels(LabelMap{
				"a": "v1",
				"c": "v2",
				"b": "v3",
			}),
			[]string{"a", "b", "c"},
			"a:v1,b:v3,c:v2",
		},
		{
			NewLabels(LabelMap{
				"a": "v1",
				"A": "v2",
				"b": "v3",
			}),
			[]string{"a", "A", "b"},
			"A:v2,a:v1,b:v3",
		},
		{
			NewLabels(LabelMap{
				"a": "v1",
				"c": "v2",
				"b": "v3",
			}),
			[]string{"a", "c"},
			"a:v1,c:v2",
		},
		{
			NewLabels(LabelMap{
				"a": "v1",
				"c": "v2",
				"b": "v3",
			}),
			[]string{"d"},
			"",
		},
	}

	var uid string
	for _, d := range data {
		uid = d.labels.ID(d.groupByLabels)
		if uid != d.expectedUID {
			t.Fatalf("Expected %s but got %s", d.expectedUID, uid)
		}
	}
}

func BenchmarkLabelKeys(b *testing.B) {
	l := NewLabels(LabelMap{"graph": "graph1", "host": "host1", "label": "label1"})
	for i := 0; i < b.N; i++ {
		l.Keys()
	}
}

func BenchmarkLabelID(b *testing.B) {
	l := NewLabels(LabelMap{"graph": "graph1", "host": "host1", "label": "label1"})
	for i := 0; i < b.N; i++ {
		l.ID([]string{"graph", "host"})
	}
}
