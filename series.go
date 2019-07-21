package muse

import (
	"github.com/google/uuid"
)

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

// Labels returns the map of label to values for the timeseries
func (s *Series) Labels() Labels {
	return s.labels
}

// UID generates the unique identifier string that represents this particular
// timeseries. This must be unique within a timeseries Group
func (s *Series) UID() string {
	return s.labels.ID(s.labels.Keys())
}
