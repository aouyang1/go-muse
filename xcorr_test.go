package muse

import (
	"math"
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/floats"
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

func TestZNormalize(t *testing.T) {
	data := []struct {
		ts []float64
	}{
		{[]float64{0, 1, 2, 3, 4, 5}},
		{[]float64{0, 1, 2, 3, 4, 5, 6}},
		{[]float64{3, 4, 3, 4, 3}},
		{[]float64{99, 100, 101, 102, 103}},
	}

	for _, d := range data {
		zNormalize(d.ts)
		var ssum float64
		for i := 0; i < len(d.ts); i++ {
			ssum += d.ts[i] * d.ts[i]
		}
		if math.Abs(ssum-float64(len(d.ts)-1)) > 1E-8 {
			t.Errorf("Expected a squared sum of %d, but got %.3f for %v", len(d.ts), ssum, d.ts)
		}
	}

}

func TestZeroPad(t *testing.T) {
	dataset := []struct {
		x         []float64
		n         int
		expectedx []float64
	}{
		{[]float64{1, 2, 3, 4}, 6, []float64{0, 0, 1, 2, 3, 4}},
		{[]float64{1, 2, 3, 4}, 3, []float64{1, 2, 3, 4}},
		{[]float64{1, 2, 3, 4}, 4, []float64{1, 2, 3, 4}},
	}

	for _, d := range dataset {
		zpadx := zeroPad(d.x, d.n)
		if len(zpadx) != len(d.expectedx) {
			t.Fatalf("Expected length %d, but got length %d", len(d.expectedx), len(zpadx))
		}
		for i, v := range zpadx {
			if d.expectedx[i] != v {
				t.Fatalf("Expected value %v, but got %v", d.expectedx[i], v)
			}
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
			[]float64{1.00, -0.25, -0.25, -0.25, -0.25},
			0,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, 0, 0, 5},
			true,
			[]float64{-0.25, -0.25, -0.25, 1.00, -0.25},
			-2,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{5, 0, 0, 0, 0},
			true,
			[]float64{-0.25, -0.25, 1.00, -0.25, -0.25},
			2,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, -5, 0, 0},
			true,
			[]float64{-1.00, 0.25, 0.25, 0.25, 0.25},
			0,
			isNegative(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{-5, 0, 0, 0, 0},
			true,
			[]float64{0.25, 0.25, -1.00, 0.25, 0.25},
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
		ExpectedXCorr []float64
		ExpectedIdx   int
		ExpectedSign  func(float64) bool
	}{
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, 5, 0, 0},
			[]float64{1.00, -0.25, -0.25, -0.25, -0.25},
			0,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, 0, 0, 5},
			[]float64{-0.25, -0.25, -0.25, 1.00, -0.25},
			-2,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{5, 0, 0, 0, 0},
			[]float64{-0.25, -0.25, 1.00, -0.25, -0.25},
			2,
			isPositive(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{0, 0, -5, 0, 0},
			[]float64{-1.00, 0.25, 0.25, 0.25, 0.25},
			0,
			isNegative(),
		},
		{
			[]float64{0, 0, 2, 0, 0},
			[]float64{-5, 0, 0, 0, 0},
			[]float64{0.25, 0.25, -1.00, 0.25, 0.25},
			2,
			isNegative(),
		},
		{
			[]float64{0, 0, 2, 2, 0},
			[]float64{3, 3, 3, 3, 3},
			[]float64{0, 0, 0, 0, 0},
			0,
			func(x float64) bool { return x == 0 },
		},
	}

	for _, ds := range datasets {
		n := len(ds.X)
		ft := fourier.NewFFT(n)
		x, err := zNormalize(ds.X)
		if err != nil {
			t.Errorf("%+v\n", err)
			continue
		}
		floats.Scale(1/float64(len(x)-1), x)
		refFT := ft.Coefficients(nil, zeroPad(x, n))

		ftY := fourier.NewFFT(n)
		xcorr, mi, mv := xCorrWithX(refFT, ds.Y, ftY)

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

func BenchmarkZPad(b *testing.B) {
	x := []float64{1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		zeroPad(x, 10)
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

func BenchmarkXCorrWithX(b *testing.B) {
	x, y, n := setupXCorrData()

	ft := fourier.NewFFT(n)
	x, err := zNormalize(x)
	if err != nil {
		b.Fatalf("%+v\n", err)
	}
	X := ft.Coefficients(nil, zeroPad(x, n))

	ftY := fourier.NewFFT(n)
	for i := 0; i < b.N; i++ {
		xCorrWithX(X, y, ftY)

	}
}
