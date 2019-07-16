package muse

import (
	"container/heap"
	"math"

	"gonum.org/v1/gonum/fourier"
)

var (
	// Concurrency sets the amount of parallelism while running cross correlations
	Concurrency = 1
)

// Muse is the primary struct to setup and run a z-normalized cross correlation between a
// reference series against each individual comparison series while tracking the resulting scores
type Muse struct {
	Reference  *Series
	Comparison *Group
	Results    *Results
}

// New creates a new Muse instance with a set reference timeseries, a
// comparison group of timeseries, and results
func New(ref *Series, comp *Group, results *Results) *Muse {
	return &Muse{Reference: ref, Comparison: comp, Results: results}
}

// scoreSingle calculates the highest score for a single set of label values given
// a reference time series
func (m *Muse) scoreSingle(idx int, refFT []complex128, labelValues Labels, n int, sem chan struct{}, graphScores []chan Score) {
	var compScore Score
	var maxVal float64
	var lag int

	maxScore := Score{}

	compGraphs := m.Comparison.FilterByLabelValues(labelValues)
	// for each time series, store the time series with highest relationship
	// with the reference time series
	for _, compTs := range compGraphs {
		// calculates the cross correlation lag and value between the reference and
		// comparison time series. boolean value specifies that we are normalizing
		// the the time series so that the power of of the reference and comparison
		// is equivalent. output value will range between 0 and 1 due to normalizing
		_, lag, maxVal = XCorrWithX(refFT, compTs.Values(), n, true)
		compScore = Score{
			Labels:       compTs.Labels(),
			Lag:          lag,
			PercentScore: int(math.Abs(maxVal*100) + 0.5),
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
func (m *Muse) Run(groupByLabels []string) {
	// Find the next power 2 that's at least twice as long as the the number of values
	// in the reference time series
	var N int
	if m.Reference.Length() > m.Comparison.Length() {
		N = nextPowOf2(float64(m.Reference.Length()))
	} else {
		N = nextPowOf2(float64(m.Comparison.Length()))
	}
	if N < 2*m.Reference.Length() {
		N = nextPowOf2(float64(N + 1))
	}

	ft := fourier.NewFFT(N)
	ref := ZeroPad(m.Reference.Values(), N)
	ref = ZNormalize(ref)
	refFT := ft.Coefficients(nil, ref)

	labelValuesSet := m.Comparison.IndexLabelValues(groupByLabels)

	// Slice of score channels will handle the output of the concurrent cross correlation
	// comparison
	graphScores := make([]chan Score, len(labelValuesSet))
	for i := range graphScores {
		graphScores[i] = make(chan Score)
	}

	// Sem channel is used to rate limit the number of concurrent go routines for cross
	// correlation comparison
	var sem = make(chan struct{}, Concurrency)
	var graphIdx int

	// Iterate over all the comparison graphs and determines the highest score a graph has
	// compared to the reference time series and stores into the slice of score channels
	for _, lv := range labelValuesSet {
		select {
		case sem <- struct{}{}:
			go m.scoreSingle(graphIdx, refFT, lv, N, sem, graphScores)
			graphIdx++
		}
	}

	// Build priority queue of size TopN so that we don't have to sort over the entire
	// score output
	heap.Init(&m.Results.scores)

	var s Score
	for _, scoreCh := range graphScores {
		s = <-scoreCh
		if m.Results.passed(s) {
			if m.Results.scores.Len() == m.Results.TopN {
				if s.PercentScore > m.Results.scores[0].PercentScore {
					heap.Pop(&m.Results.scores)
					heap.Push(&m.Results.scores, s)
				}
			} else {
				heap.Push(&m.Results.scores, s)
			}
		}
	}
}

// Results tracks the top scores in sorted order given a specified maximum lag, top N
// and score threshold
type Results struct {
	MaxLag    int
	TopN      int
	Threshold float64
	scores    Scores
}

// NewResults creates a new instance of results to track the top similar graphs
func NewResults(maxLag int, topN int, threshold float64) *Results {
	return &Results{
		MaxLag:    maxLag,
		TopN:      topN,
		Threshold: threshold,
		scores:    make([]Score, 0, topN),
	}
}

// passed checks if the input score satisfies the Results lag and threshold requirements
func (r *Results) passed(s Score) bool {
	return math.Abs(float64(s.Lag)) <= float64(r.MaxLag) &&
		math.Abs(float64(s.PercentScore)) >= r.Threshold*100
}

// Fetch returns the sorted scores in ascending order along with the average percent score
func (r *Results) Fetch() (Scores, float64) {
	s := make(Scores, 0, len(r.scores))
	var score Score
	var scoreSum int
	numScores := len(r.scores)

	for len(r.scores) > 0 {
		score = heap.Pop(&r.scores).(Score)
		scoreSum += score.PercentScore
		s = append(s, score)
	}
	return s, float64(scoreSum) / float64(numScores)
}

// Scores is a slice of individual Score
type Scores []Score

// Score keeps track of the cross correlation score and the related series
type Score struct {
	Labels       Labels `json:"labels"`
	Lag          int    `json:"lag"`
	PercentScore int    `json:"percentScore"`
}

func (s Scores) Len() int {
	return len(s)
}

func (s Scores) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Scores) Less(i, j int) bool {
	return s[i].PercentScore < s[j].PercentScore
}

// Push implements the function in the heap interface
func (s *Scores) Push(x interface{}) {
	*s = append(*s, x.(Score))
}

// Pop implements the function in the heap interface
func (s *Scores) Pop() interface{} {
	x := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return x
}
