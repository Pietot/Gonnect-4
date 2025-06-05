package grid

import (
	"fmt"
	"math"
)

type evaluation struct {
	score         *float64
	bestMove      *int
	remainingMove *int
}

const COMPUTER = 1

var stack = 0

func (s *evaluation) String() string {
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

func (grid *Grid) Negamax(player int) *evaluation {
	if grid.IsDraw() {
		score := 0.0
		remainingMoves := stack
		return &evaluation{&score, nil, &remainingMoves}
	}

	var bestEvaluation = &evaluation{nil, nil, nil}

	for column := range 7 {
		copyGrid := grid.DeepCopy()
		droppedPiece, line := copyGrid.DropPiece(column, player)
		if droppedPiece && copyGrid.CheckWinFromIndex(player, line, column) {
			remainingMoves := stack
			score := math.Inf(1)
			return &evaluation{&score, &column, &remainingMoves}
		}
	}

	for column := range 7 {
		copyGrid := grid.DeepCopy()
		droppedPiece, _ := copyGrid.DropPiece(column, player)
		if droppedPiece {
			stack++
			childEvaluation := copyGrid.Negamax(getOpponent(player)).Negate()
			stack--
			newEvaluation := &evaluation{
				score:         childEvaluation.score,
				bestMove:      &column,
				remainingMove: childEvaluation.remainingMove,
			}
			if newEvaluation.isBetterThan(bestEvaluation) {
				bestEvaluation = newEvaluation
			}
		}
	}

	return bestEvaluation
}

func getOpponent(player int) int {
	if player == COMPUTER {
		return 2
	}
	return 1
}
