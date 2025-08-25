package evaluation

import (
	"fmt"
)

type Analyzation struct {
	Scores         [7]*int8
	RemainingMoves *uint8
	BestMove       *int
}

func (a Analyzation) String() string {
	var scoreStr, bestMoveStr, remainingMoveStr string

	for i := range 7 {
		if a.Scores[i] != nil {
			scoreStr += fmt.Sprintf("C%d: %d\n", i+1, *a.Scores[i])
			if a.BestMove == nil || *a.Scores[i] > *a.Scores[*a.BestMove] {
				bestMoveStr = fmt.Sprintf("C%d", i+1)
			}
		}
	}
	remainingMoveStr = fmt.Sprintf("%d", *a.RemainingMoves)

	return "Score: " + scoreStr +
		", \nBest Move: " + bestMoveStr +
		", \nRemaining Moves: " + remainingMoveStr
}
