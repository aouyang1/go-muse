package muse

import (
	"errors"
	"fmt"
	"math"

	"gonum.org/v1/gonum/dsp/fourier"
	"gonum.org/v1/gonum/floats"
)

// Muse is the primary struct to setup and run a z-normalized cross correlation between a
// reference series against an individual comparison series while tracking the resulting scores
type Muse struct {
	refN    int          // length of the input reference
	n       int          // fourier transform length
	x       []complex128 // z-normalized fourier transform of the reference to be reused
	Results *Results
}

// New creates a new Muse instance with a set reference timeseries, and a comparison
// timeseries, and results
func New(ref *Series, results *Results) (*Muse, error) {
	if ref.Length() < 1 {
		return nil, errors.New("Reference series length must be greater than zero")
	}
	n := nextPowOf2(float64(ref.Length()))
	ft := fourier.NewFFT(n)
	x, err := zNormalize(ref.Values())
	if err != nil {
		return nil, fmt.Errorf("Invalid input query, %v", err)
	}
	floats.Scale(1/float64(len(x)-1), x)
	x = zeroPad(x, n)

	return &Muse{
		refN:    ref.Length(),
		n:       n,
		x:       ft.Coefficients(nil, x),
		Results: results,
	}, nil
}

// Run compares a single comparison series against the reference series and updates
// the score results
func (m *Muse) Run(compGraphs []*Series) error {
	if len(compGraphs) == 0 {
		// nothing to compare so don't allocate anything
		return nil
	}

	var compScore Score
	var maxVal float64
	var lag int

	maxScore := Score{}
	ft := fourier.NewFFT(m.n)
	coefScratch := make([]complex128, m.n/2+1)
	seqScratch := make([]float64, m.n)

	// for each time series, store the time series with highest relationship
	// with the reference time series
	for _, compTs := range compGraphs {
		// calculates the cross correlation lag and value between the reference and
		// comparison time series. boolean value specifies that we are normalizing
		// the the time series so that the power of of the reference and comparison
		// is equivalent. output value will range between 0 and 1 due to normalizing
		if compTs.Length() != m.refN {
			return fmt.Errorf("Encountered a comparison graph with differing length than the reference, %+v", compTs.Labels())
		}
		_, lag, maxVal = xCorrWithX(m.x, compTs.Values(), ft, coefScratch, seqScratch)
		if maxVal > 1.0 {
			maxVal = 1.0
		} else if maxVal < -1.0 {
			maxVal = -1.0
		}

		compScore = Score{
			Labels:       compTs.Labels(),
			Lag:          lag,
			PercentScore: maxVal,
		}

		// retain the score if it's the highest recorded scoring time series for the
		// current graph
		if math.Abs(compScore.PercentScore) > math.Abs(maxScore.PercentScore) || maxScore.Labels == nil {
			maxScore = compScore
		}
	}
	m.Results.Update(maxScore)
	return nil
}
