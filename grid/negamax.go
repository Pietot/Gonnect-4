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

func (grid *Grid) isDraw() bool {
	return grid.nbMoves == 6*7
}

func (grid *Grid) deepCopy() *Grid {
	newGrid := make([][]int, len(grid.Grid))
	for i := range grid.Grid {
		newGrid[i] = make([]int, len(grid.Grid[i]))
		copy(newGrid[i], grid.Grid[i])
	}
	return &Grid{
		Grid:    newGrid,
		nbMoves: grid.nbMoves,
	}
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
