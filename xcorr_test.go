package muse

import (
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/fourier"
)

func isPositive() func(float64) bool {
	return func(x float64) bool { return x > 0 }
}

func isNegative() func(float64) bool {
	return func(x float64) bool { return x < 0 }
}

func TestNextPowOf2(t *testing.T) {
	data := []struct {
		val      float64
		expected int
	}{
		{1.0, 1},
		{1.5, 2},
		{4.5, 8},
		{15.9, 16},
		{-5, 0},
		{0, 0},
	}

	for _, d := range data {
		if val := nextPowOf2(d.val); val != d.expected {
			t.Errorf("Expected %d, but got %d", d.expected, val)
		}
	}
}

func TestXCorr(t *testing.T) {

	datasets := []struct {
		X             []float64
		Y             []float64
		Normalize     bool
		ExpectedXCorr []float64
		ExpectedIdx   int
		ExpectedSign  func(float64) bool
	}{
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, 5, 0, 0},
			false,
			[]float64{10, 0, 0, 0, 0},
			0,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, 0, 0, 5},
			false,
			[]float64{0, 0, 0, 10, 0},
			-2,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{5, 0, 0, 0, 0},
			false,
			[]float64{0, 0, 10, 0, 0},
			2,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, -5, 0, 0},
			false,
			[]float64{-10, 0, 0, 0, 0},
			0,
			isNegative(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{-5, 0, 0, 0, 0},
			false,
			[]float64{0, 0, -10, 0, 0},
			2,
			isNegative(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, 5, 0, 0},
			true,
			[]float64{0.96, -0.24, -0.24, -0.24, -0.24},
			0,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, 0, 0, 5},
			true,
			[]float64{-0.24, -0.24, -0.24, 0.96, -0.24},
			-2,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{5, 0, 0, 0, 0},
			true,
			[]float64{-0.24, -0.24, 0.96, -0.24, -0.24},
			2,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, -5, 0, 0},
			true,
			[]float64{-0.96, 0.24, 0.24, 0.24, 0.24},
			0,
			isNegative(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{-5, 0, 0, 0, 0},
			true,
			[]float64{0.24, 0.24, -0.96, 0.24, 0.24},
			2,
			isNegative(),
		},
		{
			[]float64{0, 0, 2, 2, 0},
			[]float64{3, 3, 3, 3, 3},
			true,
			[]float64{0, 0, 0, 0, 0},
			0,
			func(x float64) bool { return x == 0 },
		},
	}

	for _, ds := range datasets {
		xcorr, mi, mv := xCorr(ds.X, ds.Y, len(ds.X), ds.Normalize)

		if !prettyClose(xcorr, ds.ExpectedXCorr) {
			t.Errorf("Expected cross correlation of %v, but got %v", ds.ExpectedXCorr, xcorr)
		}

		if mi != ds.ExpectedIdx {
			t.Errorf("Expected max index to be at %d, but found it at %d", ds.ExpectedIdx, mi)
		}

		if !ds.ExpectedSign(mv) {
			t.Errorf("Max value of, %f, sign evaluated to %t", mv, ds.ExpectedSign(mv))
		}

	}
}

func TestXCorrWithX(t *testing.T) {

	datasets := []struct {
		X             []float64
		Y             []float64
		Normalize     bool
		ExpectedXCorr []float64
		ExpectedIdx   int
		ExpectedSign  func(float64) bool
	}{
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, 5, 0, 0},
			false,
			[]float64{10, 0, 0, 0, 0},
			0,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, 0, 0, 5},
			false,
			[]float64{0, 0, 0, 10, 0},
			-2,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{5, 0, 0, 0, 0},
			false,
			[]float64{0, 0, 10, 0, 0},
			2,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, -5, 0, 0},
			false,
			[]float64{-10, 0, 0, 0, 0},
			0,
			isNegative(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{-5, 0, 0, 0, 0},
			false,
			[]float64{0, 0, -10, 0, 0},
			2,
			isNegative(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, 5, 0, 0},
			true,
			[]float64{0.96, -0.24, -0.24, -0.24, -0.24},
			0,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, 0, 0, 5},
			true,
			[]float64{-0.24, -0.24, -0.24, 0.96, -0.24},
			-2,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{5, 0, 0, 0, 0},
			true,
			[]float64{-0.24, -0.24, 0.96, -0.24, -0.24},
			2,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, -5, 0, 0},
			true,
			[]float64{-0.96, 0.24, 0.24, 0.24, 0.24},
			0,
			isNegative(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{-5, 0, 0, 0, 0},
			true,
			[]float64{0.24, 0.24, -0.96, 0.24, 0.24},
			2,
			isNegative(),
		},
		{
			[]float64{0, 0, 2, 2, 0},
			[]float64{3, 3, 3, 3, 3},
			true,
			[]float64{0, 0, 0, 0, 0},
			0,
			func(x float64) bool { return x == 0 },
		},
	}

	for _, ds := range datasets {
		n := len(ds.X)
		ft := fourier.NewFFT(n)
		ref := zeroPad(ds.X, n)
		if ds.Normalize {
			ref = zNormalize(ref)
		}
		refFT := ft.Coefficients(nil, ref)

		xcorr, mi, mv := xCorrWithX(refFT, ds.Y, n, ds.Normalize)

		if !prettyClose(xcorr, ds.ExpectedXCorr) {
			t.Errorf("Expected cross correlation of %v, but got %v", ds.ExpectedXCorr, xcorr)
		}

		if mi != ds.ExpectedIdx {
			t.Errorf("Expected max index to be at %d, but found it at %d", ds.ExpectedIdx, mi)
		}

		if !ds.ExpectedSign(mv) {
			t.Errorf("Max value of, %f, sign evaluated to %t", mv, ds.ExpectedSign(mv))
		}

	}
}
func BenchmarkZNormalize(b *testing.B) {
	x := []float64{1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		zNormalize(x)
	}
}

func setupXCorrData() ([]float64, []float64, int) {
	n := 16385
	x := make([]float64, n)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = rand.Float64()
		y[i] = rand.Float64()
	}

	return x, y, nextPowOf2(float64(n))
}

func BenchmarkXCorr(b *testing.B) {
	x, y, n := setupXCorrData()

	for i := 0; i < b.N; i++ {
		xCorr(x, y, n, false)
	}
}

func BenchmarkXCorrNormalize(b *testing.B) {
	x, y, n := setupXCorrData()

	for i := 0; i < b.N; i++ {
		xCorr(x, y, n, true)
	}
}

func BenchmarkXCorrWithXNormalize(b *testing.B) {
	x, y, n := setupXCorrData()

	ft := fourier.NewFFT(n)
	X := ft.Coefficients(nil, zNormalize(zeroPad(x, n)))

	for i := 0; i < b.N; i++ {
		xCorrWithX(X, y, n, true)
	}
}
