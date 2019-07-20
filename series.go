package muse

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
)

var (
	// DefaultLabel is the label name if a series is specified without any labels
	DefaultLabel = "uid"
)

// Labels is a map of label names to label values
type Labels map[string]string

// Keys returns the sorted slice of label names
func (l Labels) Keys() []string {
	keys := make([]string, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// UID constructs the unique identifier based on an input set of labels.
// This does not have to be all the unique label names. Format will have the
// following "key1:val1,key2:val2" and so on
func (l Labels) UID(labels []string) string {
	var labelValues, lv string
	if len(labels) == 0 {
		labels = l.Keys()
	} else {
		sort.Strings(labels)
	}
	for _, label := range labels {
		if _, exists := l[label]; exists {
			lv = label + ":" + l[label] + ","
		}
		labelValues = labelValues + lv
	}

	// trim last comma if present
	if len(labelValues) > 0 && labelValues[len(labelValues)-1] == ',' {
		return labelValues[:len(labelValues)-1]
	}

	return labelValues
}

// Series is the general representation of timeseries containing only values.
type Series struct {
	y      []float64
	labels Labels
}

// NewSeries creates a new Series with a set of labels. If not labels are
// specified a unique ID is automatically generated
func NewSeries(y []float64, labels Labels) *Series {
	if len(labels) == 0 {
		labels = Labels{DefaultLabel: uuid.New().String()}
	}

	return &Series{y: y, labels: labels}
}

// Length returns the length of the timeseries
func (s *Series) Length() int {
	return len(s.y)
}

// Values returns the y or series values
func (s *Series) Values() []float64 {
	return s.y
}

// GetLabelValue returns whether the label name exists in the timeseries labels
// and returns the value if present
func (s *Series) GetLabelValue(label string) (bool, string) {
	for k, v := range s.labels {
		if k == label {
			return true, v
		}
	}
	return false, ""
}

// Labels returns the map of label to values for the timeseries
func (s *Series) Labels() Labels {
	return s.labels
}

// UID generates the unique identifier string that represents this particular
// timeseries. This must be unique within a timeseries Group
func (s *Series) UID() string {
	return s.labels.UID(s.labels.Keys())
}

// Group is a collection of timeseries keeping track of all labeled timeseries,
// All timeseries must be unique regarding their label value pairs
type Group struct {
	Name     string
	n        int                 // length of each timeseries in the group
	index    map[string][]string // mapping of the grouped labels to a slice of Series UIDs with the same group label
	registry map[string]*Series  // stores a mapping of the Series UID to the Series instance
}

// NewGroup creates a new Group and initializes the timeseries label registry
func NewGroup(name string) *Group {
	return &Group{
		Name:     name,
		index:    make(map[string][]string),
		registry: make(map[string]*Series),
	}
}

// Length returns the length of all timeseries. All timeseries have the same length
func (g *Group) Length() int {
	return g.n
}

// Add will register a time series with its labels into the current groups
// registry. If the timeseries with the exact same label values already exists,
// an error will be returned
func (g *Group) Add(series ...*Series) error {
	for _, s := range series {
		labels := s.labels.Keys()
		if len(labels) == 0 {
			return fmt.Errorf("Invalid Series with no labels, %v", s)
		}

		uid := s.UID()
		if _, exists := g.registry[uid]; exists {
			return fmt.Errorf("Series with label:values, %v, already exists within group, %s", uid, g.Name)
		}

		// set the length of timeseries for this group or check if the added timeseries
		// has the same length
		if len(g.registry) == 0 {
			g.n = s.Length()
		} else {
			if s.Length() != g.n {
				return fmt.Errorf("Timeseries has length %d, but current group has length %d", s.Length(), g.n)
			}
		}

		g.registry[uid] = s
	}
	return nil
}

// FilterByLabelValues returns the slice of timeseries filtered by specified label
// value pairs
func (g *Group) FilterByLabelValues(labels Labels) []*Series {
	var filteredSeries []*Series

	guid := labels.UID(labels.Keys())
	if _, exists := g.index[guid]; exists {
		filteredSeries = make([]*Series, 0, len(g.index[guid]))
		for _, uid := range g.index[guid] {
			filteredSeries = append(filteredSeries, g.registry[uid])
		}
	}
	return filteredSeries
}

// indexLabelValues return a slice of all the distinct combinations of the
// input label values while ignoring labels not being specified. If no labels
// are specified then each series will be treated separately.
func (g *Group) indexLabelValues(groupByLabels []string) []Labels {
	var distinctLabelValues []Labels
	var guid string

	// clear index
	g.index = make(map[string][]string)

	for uid, s := range g.registry {
		if len(groupByLabels) != 0 {
			guid = s.labels.UID(groupByLabels)
		} else {
			guid = uid
			groupByLabels = s.Labels().Keys()
		}
		if _, exists := g.index[guid]; !exists {
			lv := make(Labels)
			for _, name := range groupByLabels {
				lv[name] = s.labels[name]
			}
			distinctLabelValues = append(distinctLabelValues, lv)
		}

		g.index[guid] = append(g.index[guid], uid)
	}

	return distinctLabelValues
}
