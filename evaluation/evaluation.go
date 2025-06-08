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
		bestMoveStr = fmt.Sprintf("C%d", *s.BestMove+1)
	} else {
		bestMoveStr = "None"
	}

	if s.RemainingMove != nil {
		remainingMoveStr = fmt.Sprintf("%d", *s.RemainingMove)
	} else {
		remainingMoveStr = "None"
	}

	return "Score: " + scoreStr +
		", \nBest Move: " + bestMoveStr +
		", \nRemaining Moves: " + remainingMoveStr
}

func (eval1 *Evaluation) IsBetterThan(eval2 *Evaluation) bool {
	// If one score is nil, the other is better
	if eval1.Score == nil {
		return false
	}
	if eval2.Score == nil {
		return true
	}

	// If scores differ, compare them
	if *eval1.Score != *eval2.Score {
		return *eval1.Score > *eval2.Score
	}

	// If scores are equal, compare remaining moves (nil = worse)
	if eval1.RemainingMove == nil {
		return false
	}
	if eval2.RemainingMove == nil {
		return true
	}

	// If winning, fewer remaining moves is better
	if *eval1.Score == math.Inf(1) {
		return *eval1.RemainingMove < *eval2.RemainingMove
	}
	// If losing, more remaining moves is better
	if *eval1.Score == math.Inf(-1) {
		return *eval1.RemainingMove > *eval2.RemainingMove
	}

	// Never reach here if both scores are non-nil and equal
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
