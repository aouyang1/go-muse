package muse

// Scores is a slice of individual Score
type Scores []Score

// Score keeps track of the cross correlation score and the related series
type Score struct {
	Labels       *Labels `json:"labels"`
	Lag          int     `json:"lag"`
	PercentScore float64 `json:"percentScore"`
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
