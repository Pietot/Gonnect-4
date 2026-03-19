package evaluation

import (
	"fmt"

	"github.com/Pietot/Gonnect-4/config"
)

type Analysis struct {
	Scores         [7]int8
	RemainingMoves uint8
	BestMove       uint8
}

func (a Analysis) String() string {
	scores := make([]string, len(a.Scores))
	for i, score := range a.Scores {
		if score == config.NIL_SCORE {
			scores[i] = "N/A "
		} else {
			scores[i] = fmt.Sprintf("%d", score)
		}
	}
	return "Scores: " + fmt.Sprintf("%v", scores) +
		", \nBest Move: " + fmt.Sprintf("C%d", a.BestMove+1) +
		", \nRemaining Moves: " + fmt.Sprintf("%d", a.RemainingMoves)
}
