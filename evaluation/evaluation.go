package evaluation

import (
	"fmt"
)

type Evaluation struct {
	Score          *int8
	BestMove       *int
	RemainingMoves *uint8
}

func (s *Evaluation) String() string {
	var scoreStr, bestMoveStr, remainingMoveStr string

	if s.Score != nil {
		scoreStr = fmt.Sprintf("%.d", *s.Score)
	} else {
		scoreStr = "None"
	}

	if s.BestMove != nil {
		bestMoveStr = fmt.Sprintf("C%d", *s.BestMove+1)
	} else {
		bestMoveStr = "None"
	}

	if s.RemainingMoves != nil {
		remainingMoveStr = fmt.Sprintf("%d", *s.RemainingMoves)
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
			Score:          e.Score,
			BestMove:       e.BestMove,
			RemainingMoves: e.RemainingMoves,
		}
	}
	neg := -(*e.Score)
	return &Evaluation{
		Score:          &neg,
		BestMove:       e.BestMove,
		RemainingMoves: e.RemainingMoves,
	}
}
