package muse

import "sort"

const (
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

// ID constructs the unique identifier based on an input set of labels.
// This does not have to be all the unique label names. Format will have the
// following "key1:val1,key2:val2" and so on
func (l Labels) ID(labels []string) string {
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
