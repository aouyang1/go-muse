package muse

import (
	"testing"

	"gonum.org/v1/gonum/fourier"
)

func BenchmarkXCorr(b *testing.B) {
	x := []float64{1, 2, 3, 4}
	y := []float64{1, 2, 3, 4}
	n := 512

	for i := 0; i < b.N; i++ {
		XCorr(x, y, n, false)
	}
}

func BenchmarkXCorrNormalize(b *testing.B) {
	x := []float64{1, 2, 3, 4}
	y := []float64{1, 2, 3, 4}
	n := 512

	for i := 0; i < b.N; i++ {
		XCorr(x, y, n, true)
	}
}

func BenchmarkXCorrWithXNormalize(b *testing.B) {
	x := []float64{1, 2, 3, 4}
	y := []float64{1, 2, 3, 4}
	n := 512

	x = ZeroPad(x, n)
	x = ZNormalize(x)

	ft := fourier.NewFFT(n)
	X := ft.Coefficients(nil, x)

	for i := 0; i < b.N; i++ {
		XCorrWithX(X, y, n, true)
	}
}

func BenchmarkXCorrBatchNormalizex1(b *testing.B) {
	x := []float64{1, 2, 3, 4}
	y := []float64{1, 2, 3, 4}
	n := 512

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
	x := []float64{1, 2, 3, 4}
	y := []float64{1, 2, 3, 4}
	n := 512

	NumSeries := 10

	var multiY [][]float64
	for rep := 0; rep < NumSeries; rep++ {
		multiY = append(multiY, y)
	}

	for i := 0; i < b.N; i++ {
		XCorrBatch(x, multiY, n, true)
	}
}
