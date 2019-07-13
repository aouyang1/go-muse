package muse

import (
	"testing"
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
		xcorr, mi, mv := XCorr(ds.X, ds.Y, len(ds.X), ds.Normalize)

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

func TestXCorrBatch(t *testing.T) {

	datasets := []struct {
		X             []float64
		Y             [][]float64
		Normalize     bool
		ExpectedXCorr [][]float64
		ExpectedIdx   []int
		ExpectedSign  []func(float64) bool
	}{
		{
			[]float64{0, 0, 3, 0, 0},
			[][]float64{
				{0, 0, 5, 0, 0},
				{0, 0, 0, 0, 5},
				{5, 0, 0, 0, 0},
				{0, 0, -5, 0, 0},
				{-5, 0, 0, 0, 0},
				{3, 3, 3, 3, 3},
			},
			false,
			[][]float64{
				{15, 0, 0, 0, 0},
				{0, 0, 0, 15, 0},
				{0, 0, 15, 0, 0},
				{-15, 0, 0, 0, 0},
				{0, 0, -15, 0, 0},
				{9, 9, 9, 9, 9},
			},
			[]int{0, -2, 2, 0, 2, 0},
			[]func(float64) bool{
				isPositive(),
				isPositive(),
				isPositive(),
				isNegative(),
				isNegative(),
				isPositive(),
			},
		},
		{
			[]float64{0, 0, 3, 0, 0},
			[][]float64{
				{0, 0, 5, 0, 0},
				{0, 0, 0, 0, 5},
				{5, 0, 0, 0, 0},
				{0, 0, -5, 0, 0},
				{-5, 0, 0, 0, 0},
				{3, 3, 3, 3, 3},
			},
			true,
			[][]float64{
				{0.96, -0.24, -0.24, -0.24, -0.24},
				{-0.24, -0.24, -0.24, 0.96, -0.24},
				{-0.24, -0.24, 0.96, -0.24, -0.24},
				{-0.96, 0.24, 0.24, 0.24, 0.24},
				{0.24, 0.24, -0.96, 0.24, .24},
				{0, 0, 0, 0, 0},
			},
			[]int{0, -2, 2, 0, 2, 0},
			[]func(float64) bool{
				isPositive(),
				isPositive(),
				isPositive(),
				isNegative(),
				isNegative(),
				func(x float64) bool { return x == 0 },
			},
		},
	}

	for _, ds := range datasets {
		xcorr, mi, mv := XCorrBatch(ds.X, ds.Y, len(ds.X), ds.Normalize)

		switch {
		case len(xcorr) != len(ds.ExpectedXCorr):
			t.Fatalf("Expected %d cross correlation results, but got %d", len(ds.ExpectedXCorr), len(xcorr))
		case len(mi) != len(ds.ExpectedIdx):
			t.Fatalf("Expected %d index results, but got %d", len(ds.ExpectedIdx), len(mi))
		case len(mv) != len(ds.ExpectedSign):
			t.Fatalf("Expected %d max value results, but got %d", len(ds.ExpectedSign), len(mv))
		}

		for i := 0; i < len(xcorr); i++ {
			switch {
			case !prettyClose(xcorr[i], ds.ExpectedXCorr[i]):
				t.Fatalf("Expected cross correlation of %v, but got %v", ds.ExpectedXCorr[i], xcorr[i])
			case mi[i] != ds.ExpectedIdx[i]:
				t.Fatalf("Expected max index to be at %d, but found it at %d", ds.ExpectedIdx[i], mi[i])
			case !ds.ExpectedSign[i](mv[i]):
				t.Fatalf("Max value of, %f, sign evaluated to %t", mv[i], ds.ExpectedSign[i](mv[i]))
			}
		}
	}
}
