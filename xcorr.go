package muse

import (
	"fmt"
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/fourier"
	"gonum.org/v1/gonum/stat"
)

func nextPowOf2(val float64) int {
	if val <= 0 {
		return 0
	}
	return int(math.Pow(2.0, math.Ceil(math.Log(val)/math.Log(2))))
}

func prettyClose(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if math.Abs(v-b[i]) > 1E-8 {
			return false
		}
	}
	return true
}

// mult multiplies two slices element by element
func mult(x []complex128, y []complex128) []complex128 {
	if len(x) != len(y) {
		panic(fmt.Sprintf("Non equivalent length of slices, x: %d, y: %d", len(x), len(y)))
	}
	out := make([]complex128, len(x))
	for i, v := range x {
		out[i] = v * y[i]
	}

	return out
}

// maxAbsIndex finds the index with the largest absolute value
func maxAbsIndex(x []float64) int {
	var maxIndex int
	var maxVal float64
	for i, v := range x {
		if math.Abs(v) > math.Abs(maxVal) {
			maxVal = v
			maxIndex = i
		}
	}

	return maxIndex
}

// conj returns the complex conjugate of a slice of complex values
func conj(x []complex128) []complex128 {
	out := make([]complex128, len(x))
	for i, v := range x {
		out[i] = cmplx.Conj(v)
	}

	return out
}

// zeroPad re-slices the input array to a size n leaving trailing zeroes
func zeroPad(x []float64, n int) []float64 {
	if n < len(x) {
		return x
	}

	xpad := make([]float64, n)
	for i := 0; i < len(x); i++ {
		xpad[i] = x[i]
	}
	return xpad
}

// zNormalize removes the mean and divides each value by the standard
// deviation of the resulting series
func zNormalize(x []float64) []float64 {
	meanX := floats.Sum(x) / float64(len(x))

	floats.AddConst(-meanX, x)

	weights := make([]float64, len(x))
	for i := 0; i < len(x); i++ {
		weights[i] = float64(len(x))
	}

	stdX := stat.StdDev(x, weights)

	if stdX != 0 {
		floats.Scale(1.0/stdX, x)
	}

	return x
}

// xCorr computes the cross correlation slice between x and y, index of the maximum absolute value
// and the maximum absolute value. You can specify number of samples to truncate both x and y slices
// or zero pad the two slices. The normalize flag will normalize both x and y slices to their own
// signal power. The resulting maximum absolute values will range from 0-1 for normalized, but not
// necessarily for non-normalized computations.
func xCorr(x []float64, y []float64, n int, normalize bool) ([]float64, int, float64) {
	// Negative lag means y is lagging behind x. Earliest timepoint is at index 0
	if minN := int(math.Max(float64(len(x)), float64(len(y)))); n < minN {
		n = minN
	}

	x = zeroPad(x, n)
	y = zeroPad(y, n)

	if normalize {
		x = zNormalize(x)
		y = zNormalize(y)
	}

	ft := fourier.NewFFT(n)

	cc := ft.Sequence(nil, mult(ft.Coefficients(nil, x), conj(ft.Coefficients(nil, y))))
	if normalize {
		floats.Scale(1.0/float64(n*n), cc)
	} else {
		floats.Scale(1.0/float64(n), cc)
	}

	mi := maxAbsIndex(cc)
	mv := cc[mi]

	if mi > n/2 {
		mi = mi - n
	}

	return cc, mi, mv
}

// xCorrWithX allows a precomputed FFT of X to be passed in for the purposes of batch
// execution and not repeatedly calculating FFT(x). Must pass in the fourier transform
// struct used to compute X.
func xCorrWithX(X []complex128, y []float64, n int, normalize bool) ([]float64, int, float64) {
	y = zeroPad(y, n)

	if normalize {
		y = zNormalize(y)
	}

	ft := fourier.NewFFT(n)
	cc := ft.Sequence(nil, mult(X, conj(ft.Coefficients(nil, y))))
	if normalize {
		floats.Scale(1.0/float64(n*n), cc)
	} else {
		floats.Scale(1.0/float64(n), cc)
	}

	mi := maxAbsIndex(cc)
	mv := cc[mi]

	if mi > n/2 {
		mi = mi - n
	}

	return cc, mi, mv
}
