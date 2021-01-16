package muse

import (
	"container/heap"
	"math"
	"sync"
)

// Results tracks the top scores in sorted order given a specified maximum lag, top N
// and score threshold
type Results struct {
	sync.Mutex
	MaxLag    int
	TopN      int
	Threshold float64
	scores    Scores
}

// NewResults creates a new instance of results to track the top similar graphs
func NewResults(maxLag int, topN int, threshold float64) *Results {
	scores := make(Scores, 0, topN)

	// Build priority queue of size TopN so that we don't have to sort over the entire
	// score output
	heap.Init(&scores)

	return &Results{
		MaxLag:    maxLag,
		TopN:      topN,
		Threshold: threshold,
		scores:    scores,
	}
}

// passed checks if the input score satisfies the Results lag and threshold requirements
func (r *Results) passed(s Score) bool {
	return math.Abs(float64(s.Lag)) <= float64(r.MaxLag) &&
		math.Abs(float64(s.PercentScore)) >= r.Threshold
}

// Update records the input score
func (r *Results) Update(s Score) {
	if s.Labels == nil {
		// invalid score so don't update anything
		return
	}
	r.Lock()
	if r.passed(s) {
		if r.scores.Len() == r.TopN {
			if math.Abs(s.PercentScore) > math.Abs(r.scores[0].PercentScore) {
				heap.Pop(&r.scores)
				heap.Push(&r.scores, s)
			}
		} else {
			heap.Push(&r.scores, s)
		}
	}
	r.Unlock()
}

// Fetch returns the sorted scores in ascending order along with the average absolute percent score
func (r *Results) Fetch() (Scores, float64) {
	s := make(Scores, len(r.scores))
	var score Score
	var scoreSum float64
	numScores := len(r.scores)

	for i := numScores - 1; i >= 0; i-- {
		score = heap.Pop(&r.scores).(Score)
		scoreSum += math.Abs(score.PercentScore)
		s[i] = score
	}
	return s, scoreSum / float64(numScores)
}
