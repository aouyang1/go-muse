package muse

import (
	"errors"
	"fmt"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/fourier"
)

// Muse is the primary struct to setup and run a z-normalized cross correlation between a
// reference series against an individual comparison series while tracking the resulting scores
type Muse struct {
	Reference *Series
	N         int          // fourier transform length
	X         []complex128 // z-normalized fourier transform of the reference to be reused
	Results   *Results
}

// New creates a new Muse instance with a set reference timeseries, and a comparison
// timeseries, and results
func New(ref *Series, results *Results) (*Muse, error) {
	if ref.Length() < 1 {
		return nil, errors.New("Reference series length must be greater than zero")
	}
	n := calculateN(ref.Length())
	ft := fourier.NewFFT(n)
	x, err := zNormalize(ref.Values())
	if err != nil {
		return nil, fmt.Errorf("Invalid input query, %v", err)
	}
	floats.Scale(1/float64(len(x)-1), x)
	x = zeroPad(x, n)
	X := ft.Coefficients(nil, x)

	m := &Muse{
		Reference: ref,
		N:         n,
		X:         X,
		Results:   results,
	}

	return m, nil
}

// Run compares a single comparison series against the reference series and updates
// the score results
func (m *Muse) Run(comp *Series) error {
	//m.scoreSingle()
	//m.Results.Update(s)
	return nil
}

func calculateN(refLen int) int {
	n := nextPowOf2(float64(refLen))
	if n < 2*refLen {
		n = nextPowOf2(float64(n + 1))
	}
	return n
}
