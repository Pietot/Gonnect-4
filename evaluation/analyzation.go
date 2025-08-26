package evaluation

import (
	"fmt"
)

type Analyzation struct {
	Scores         [7]*int8
	RemainingMoves *uint8
	BestMove       *uint8
}

func (a Analyzation) String() string {
	scores := make([]string, len(a.Scores))
	for i, ptr := range a.Scores {
		if ptr != nil {
			scores[i] = fmt.Sprintf("%d", *ptr)
		} else {
			scores[i] = "N/A"
		}
	}
	return "Scores: " + fmt.Sprintf("%v", scores) +
		", \nBest Move: " + fmt.Sprintf("C%d", *a.BestMove+1) +
		", \nRemaining Moves: " + fmt.Sprintf("%d", *a.RemainingMoves)
}
