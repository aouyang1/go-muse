package muse

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/fourier"
	"gonum.org/v1/gonum/stat"
)

var (
	errStdDevZero = errors.New("Standard deviation of zero")
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

// mult multiplies two slices element by element saving in the dst slice
func mult(dst, src []complex128) {
	if len(dst) != len(src) {
		panic(fmt.Sprintf("Non equivalent length of slices, x: %d, y: %d", len(dst), len(src)))
	}
	for i, v := range dst {
		dst[i] = v * src[i]
	}
}

// conj changes the input into the complex conjugate of a slice of complex values
func conj(x []complex128) {
	for i, v := range x {
		x[i] = cmplx.Conj(v)
	}
}

// zeroPad re-slices the input array to a size n with leading zeroes
func zeroPad(x []float64, n int) []float64 {
	if n < len(x) {
		return x
	}

	xpad := make([]float64, n)
	for i := 0; i < len(x); i++ {
		xpad[n-len(x)+i] = x[i]
	}
	return xpad
}

// zNormalize removes the mean and divides each value by the standard
// deviation of the resulting series
func zNormalize(x []float64) ([]float64, error) {
	n := float64(len(x))
	floats.AddConst(-floats.Sum(x)/n, x)

	stdX := stat.StdDev(x, nil)

	if stdX == 0 {
		return nil, errStdDevZero
	}
	floats.Scale(1/stdX, x)
	return x, nil
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

	if normalize {
		var err error
		x, err = zNormalize(x)
		if err != nil {
			if err.Error() == errStdDevZero.Error() {
				return nil, 0, 0
			}
			// Unknown error from zNormalize
			log.Printf("%+v\n", err)
			return nil, 0, 0
		}
		y, err = zNormalize(y)
		if err != nil {
			if err.Error() == errStdDevZero.Error() {
				return nil, 0, 0
			}
			// Unknown error from zNormalize
			log.Printf("%+v\n", err)
			return nil, 0, 0
		}
	}
	x = zeroPad(x, n)
	y = zeroPad(y, n)

	ft := fourier.NewFFT(n)

	X := ft.Coefficients(nil, x)
	Y := ft.Coefficients(nil, y)
	conj(Y)
	mult(X, Y)
	cc := ft.Sequence(nil, X)
	if normalize {
		floats.Scale(1.0/float64(n*(n-1)), cc)
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
func xCorrWithX(X []complex128, y []float64, ft *fourier.FFT) ([]float64, int, float64) {
	var err error

	n := ft.Len()
	y, err = zNormalize(y)
	if err != nil {
		if err.Error() == errStdDevZero.Error() {
			return nil, 0, 0
		}
		// Unknown error from zNormalize
		log.Printf("%+v\n", err)
		return nil, 0, 0
	}
	y = zeroPad(y, n)

	C := ft.Coefficients(nil, y)
	conj(C)
	mult(C, X)
	cc := ft.Sequence(nil, C)
	floats.Scale(1.0/float64(n), cc)

	mi := maxAbsIndex(cc)
	mv := cc[mi]

	if mi > n/2 {
		mi = mi - n
	}

	return cc, mi, mv
}
