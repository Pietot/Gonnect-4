package grid

import (
	"fmt"
	"math"
)

type score struct {
	score         *float64
	bestMove      *int
	remainingMove *int
}

func (s *score) String() string {
	var scoreStr, bestMoveStr, remainingMoveStr string

	if s.score != nil {
		scoreStr = fmt.Sprintf("%.2f", *s.score)
	} else {
		scoreStr = "None"
	}

	if s.bestMove != nil {
		bestMoveStr = fmt.Sprintf("%d", *s.bestMove+1)
	} else {
		bestMoveStr = "None"
	}

	if s.remainingMove != nil {
		remainingMoveStr = fmt.Sprintf("%d", *s.remainingMove)
	} else {
		remainingMoveStr = "None"
	}

	return "Score: " + scoreStr +
		", \nBest Move: C" + bestMoveStr +
		", \nRemaining Moves: " + remainingMoveStr
}

func (eval1 *evaluation) isBetterThan(eval2 *evaluation) bool {
	// A defined score is always better than a nil score
	if eval2.score == nil {
		return true
	}
	if eval1.score == nil {
		return false
	}

	// If both scores are defined and different, compare them
	if *eval1.score != *eval2.score {
		return *eval1.score > *eval2.score
	}

	// If scores are equal and victory is imminent, minimize the remaining moves
	if *eval1.score == math.Inf(1) {
		return *eval1.remainingMove < *eval2.remainingMove
	}

	// If scores are equal and defeat is imminent, maximize the remaining moves
	if *eval1.score == math.Inf(-1) {
		return *eval1.remainingMove > *eval2.remainingMove
	}

	// Useless case
	return false
}

func (grid *Grid) Negamax(depth int, maximizingPlayer int) *score {
	if grid.isDraw() {
		val := 0.0
		return &score{&val, nil, nil}
	}
	copyGrid := grid.deepCopy()
	for column := range 7 {
		canPlay, line := copyGrid.DropPiece(column, maximizingPlayer)
		if canPlay && copyGrid.CheckWinFromIndex(maximizingPlayer, line, column) {
			val := math.Inf(1)
			remainingMove := 1
			return &score{&val, &column, &remainingMove}
		}
	}
	return &score{}
}
