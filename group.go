package muse

import "fmt"

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
func (g *Group) FilterByLabelValues(labels *Labels) []*Series {
	var filteredSeries []*Series

	guid := labels.ID(labels.Keys())
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
func (g *Group) indexLabelValues(groupByLabels []string) []*Labels {
	var distinctLabelValues []*Labels
	var guid string

	// clear index
	g.index = make(map[string][]string)

	for uid, s := range g.registry {
		if len(groupByLabels) != 0 {
			guid = s.labels.ID(groupByLabels)
		} else {
			guid = uid
			groupByLabels = s.Labels().Keys()
		}
		if _, exists := g.index[guid]; !exists {
			lv := make(LabelMap)
			for _, name := range groupByLabels {
				if v, exists := s.labels.Get(name); exists {
					lv[name] = v
				}
			}
			distinctLabelValues = append(distinctLabelValues, NewLabels(lv))
		}

		g.index[guid] = append(g.index[guid], uid)
	}

	return distinctLabelValues
}
