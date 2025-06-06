package evaluation

import (
	"fmt"
	"math"
)

type Evaluation struct {
	Score         *float64
	BestMove      *int
	RemainingMove *int
}

func (s *Evaluation) String() string {
	var scoreStr, bestMoveStr, remainingMoveStr string

	if s.Score != nil {
		scoreStr = fmt.Sprintf("%.2f", *s.Score)
	} else {
		scoreStr = "None"
	}

	if s.BestMove != nil {
		bestMoveStr = fmt.Sprintf("%d", *s.BestMove+1)
	} else {
		bestMoveStr = "None"
	}

	if s.RemainingMove != nil {
		remainingMoveStr = fmt.Sprintf("%d", *s.RemainingMove)
	} else {
		remainingMoveStr = "None"
	}

	return "Score: " + scoreStr +
		", \nBest Move: C" + bestMoveStr +
		", \nRemaining Moves: " + remainingMoveStr
}

func (eval1 *Evaluation) IsBetterThan(eval2 *Evaluation) bool {
	if eval1.Score == nil && eval2.Score == nil {
		return false
	}
	// A defined score is always better than a nil score
	if eval1.Score != nil && eval2.Score == nil {
		return true
	}
	if eval1.Score == nil && eval2.Score != nil {
		return false
	}

	// If both scores are defined and different, compare them
	if *eval1.Score != *eval2.Score {
		return *eval1.Score > *eval2.Score
	}

	// If scores are equal and victory is imminent, minimize the remaining moves
	if *eval1.Score == math.Inf(1) {
		return *eval1.RemainingMove < *eval2.RemainingMove
	}

	// If scores are equal and defeat is imminent, maximize the remaining moves
	if *eval1.Score == math.Inf(-1) {
		return *eval1.RemainingMove > *eval2.RemainingMove
	}

	// Useless case, never reached
	return false
}

func (e *Evaluation) Negate() *Evaluation {
	if e.Score == nil || *e.Score == 0.0 {
		return &Evaluation{
			Score:         e.Score,
			BestMove:      e.BestMove,
			RemainingMove: e.RemainingMove,
		}
	}
	neg := -(*e.Score)
	return &Evaluation{
		Score:         &neg,
		BestMove:      e.BestMove,
		RemainingMove: e.RemainingMove,
	}
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func IntPtr(i int) *int {
	return &i
}
