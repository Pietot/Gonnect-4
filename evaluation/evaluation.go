package evaluation

import (
	"fmt"
)

type Evaluation struct {
	Score          *int8
	BestMove       *uint8
	RemainingMoves *uint8
}

func (s *Evaluation) String() string {
	var scoreStr, bestMoveStr, remainingMoveStr string

	if s.Score != nil {
		scoreStr = fmt.Sprintf("%d", *s.Score)
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
