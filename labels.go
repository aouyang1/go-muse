package muse

import "sort"

const (
	// DefaultLabel is the label name if a series is specified without any labels
	DefaultLabel = "uid"
)

// LabelMap is a map of label keys to values
type LabelMap map[string]string

// Labels is a map of label names to label values
type Labels struct {
	labels LabelMap
	keys   []string
}

// NewLabels creates a Label storing the map and sorted keys
func NewLabels(labels LabelMap) *Labels {
	l := &Labels{
		labels: labels,
		keys:   make([]string, 0, len(labels)),
	}
	for k := range labels {
		l.keys = append(l.keys, k)
	}
	sort.Strings(l.keys)
	return l
}

// Len returns the number of labels
func (l *Labels) Len() int {
	return len(l.labels)
}

// Keys returns the sorted keys of the labels
func (l *Labels) Keys() []string {
	return l.keys
}

// Get returns the value of a specified key. Returns an empty
// string if the key is not present
func (l *Labels) Get(key string) (string, bool) {
	if v, exists := l.labels[key]; exists {
		return v, true
	}
	return "", false
}

// ID constructs the unique identifier based on an input set of labels.
// This does not have to be all the unique label names. Format will have the
// following "key1:val1,key2:val2" and so on
func (l Labels) ID(labels []string) string {
	var labelValues string
	if len(labels) == 0 {
		labels = l.Keys()
	} else {
		sort.Strings(labels)
	}
	for _, label := range labels {
		if v, exists := l.Get(label); exists {
			labelValues = labelValues + label + ":" + v + ","
		}
	}

	// trim last comma if present
	if len(labelValues) > 0 && labelValues[len(labelValues)-1] == ',' {
		return labelValues[:len(labelValues)-1]
	}

	return labelValues
}
