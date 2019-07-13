package muse

import (
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/fourier"
)

func BenchmarkZNormalize(b *testing.B) {
	x := []float64{1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		ZNormalize(x)
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

func BenchmarkFFT(b *testing.B) {
	x, _, _ := setupXCorrData()
	ft := fourier.NewFFT(len(x))
	for i := 0; i < b.N; i++ {
		ft.Coefficients(nil, x)
	}
}

func BenchmarkIFFT(b *testing.B) {
	x, _, _ := setupXCorrData()
	ft := fourier.NewFFT(len(x))
	X := ft.Coefficients(nil, x)
	for i := 0; i < b.N; i++ {
		ft.Sequence(nil, X)
	}
}
func BenchmarkXCorr(b *testing.B) {
	x, y, n := setupXCorrData()

	for i := 0; i < b.N; i++ {
		XCorr(x, y, n, false)
	}
}

func BenchmarkXCorrNormalize(b *testing.B) {
	x, y, n := setupXCorrData()

	for i := 0; i < b.N; i++ {
		XCorr(x, y, n, true)
	}
}

func BenchmarkXCorrWithXNormalize(b *testing.B) {
	x, y, n := setupXCorrData()

	ft := fourier.NewFFT(n)
	X := ft.Coefficients(nil, ZNormalize(ZeroPad(x, n)))

	for i := 0; i < b.N; i++ {
		XCorrWithX(X, y, n, true)
	}
}

func BenchmarkXCorrBatchNormalizex1(b *testing.B) {
	x, y, n := setupXCorrData()
	NumSeries := 1

	var multiY [][]float64
	for rep := 0; rep < NumSeries; rep++ {
		multiY = append(multiY, y)
	}

	for i := 0; i < b.N; i++ {
		XCorrBatch(x, multiY, n, true)
	}
}

func BenchmarkXCorrBatchNormalizex10(b *testing.B) {
	x, y, n := setupXCorrData()
	NumSeries := 10

	var multiY [][]float64
	for rep := 0; rep < NumSeries; rep++ {
		multiY = append(multiY, y)
	}

	for i := 0; i < b.N; i++ {
		XCorrBatch(x, multiY, n, true)
	}
}
