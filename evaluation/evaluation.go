package evaluation

import (
	"fmt"
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
