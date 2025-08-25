package evaluation

import (
	"fmt"
)

type Evaluation struct {
	Score          *int8
	RemainingMoves *uint8
}

func (e Evaluation) String() string {
	scoreStr := fmt.Sprintf("%d", *e.Score)
	remainingMoveStr := fmt.Sprintf("%d", *e.RemainingMoves)

	return "Score: " + scoreStr +
		", \nRemaining Moves: " + remainingMoveStr
}
