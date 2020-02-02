package muse

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/fourier"
)

// Batch is used to setup and run a z-normalized cross correlation between a
// reference series against each individual comparison series while tracking the resulting scores
type Batch struct {
	n           int
	x           []complex128
	Comparison  *Group
	Results     *Results
	Concurrency int
}

// NewBatch creates a new Muse instance with a set reference timeseries, a
// comparison group of timeseries, and results
func NewBatch(ref *Series, comp *Group, results *Results, cc int) (*Batch, error) {
	for uid, s := range comp.registry {
		if ref.Length() != s.Length() {
			return nil, fmt.Errorf("%s from comparison group series does not have the same length as the reference", uid)
		}
	}
	if cc < 1 {
		cc = 1
	}

	// Find the next power 2 that's at least twice as long as the the number of values
	// in the reference time series
	n := calculateN(ref.Length())

	ft := fourier.NewFFT(n)
	x, err := zNormalize(ref.Values())
	if err != nil {
		return nil, fmt.Errorf("Invalid input query, %v", err)
	}
	floats.Scale(1/float64(len(x)-1), x)
	x = zeroPad(x, n)

	return &Batch{
		n:           n,
		x:           ft.Coefficients(nil, x),
		Comparison:  comp,
		Results:     results,
		Concurrency: cc,
	}, nil
}

// scoreSingle calculates the highest score for a single set of label values given
// a reference time series
func (b *Batch) scoreSingle(idx int, labelValues *Labels, sem chan struct{}, graphScores []chan Score) {
	var compScore Score
	var maxVal float64
	var lag int

	maxScore := Score{}
	ft := fourier.NewFFT(b.n)

	compGraphs := b.Comparison.FilterByLabelValues(labelValues)
	// for each time series, store the time series with highest relationship
	// with the reference time series
	for _, compTs := range compGraphs {
		// calculates the cross correlation lag and value between the reference and
		// comparison time series. boolean value specifies that we are normalizing
		// the the time series so that the power of of the reference and comparison
		// is equivalent. output value will range between 0 and 1 due to normalizing
		_, lag, maxVal = xCorrWithX(b.x, compTs.Values(), ft)
		compScore = Score{
			Labels:       compTs.Labels(),
			Lag:          lag,
			PercentScore: math.Abs(maxVal),
		}

		// retain the score if it's the highest recorded scoring time series for the
		// current graph
		if compScore.PercentScore > maxScore.PercentScore || maxScore.Labels == nil {
			maxScore = compScore
		}
	}
	<-sem
	graphScores[idx] <- maxScore
}

// Run calculates the top N graphs with the highest scores given a reference time
// series and a group of comparison time series. Number of scores will be the number
// of unique labels specified in the input. If no groupByLabels is specified, then
// each timeseries will receive its own score.
func (m *Batch) Run(groupByLabels []string) error {
	labelValuesSet := m.Comparison.indexLabelValues(groupByLabels)

	// Slice of score channels will handle the output of the concurrent cross correlation
	// comparison
	graphScores := make([]chan Score, len(labelValuesSet))
	for i := range graphScores {
		graphScores[i] = make(chan Score)
	}

	// Sem channel is used to rate limit the number of concurrent go routines for cross
	// correlation comparison
	var sem = make(chan struct{}, m.Concurrency)
	var graphIdx int

	// Iterate over all the comparison graphs and determines the highest score a graph has
	// compared to the reference time series and stores into the slice of score channels
	for _, lv := range labelValuesSet {
		select {
		case sem <- struct{}{}:
			go m.scoreSingle(graphIdx, lv, sem, graphScores)
			graphIdx++
		}
	}

	var s Score
	for _, scoreCh := range graphScores {
		s = <-scoreCh
		m.Results.Update(s)
	}
	return nil
}
