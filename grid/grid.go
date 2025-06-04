package grid

import (
	"fmt"
)

type Grid struct {
	Grid [][]int
}

func InitGrid() *Grid {
	grid := &Grid{
		Grid: make([][]int, 6),
	}
	for i := range grid.Grid {
		grid.Grid[i] = make([]int, 7)
	}
	return grid
}

func PrintGrid(grid *Grid) {
	for _, row := range grid.Grid {
		for _, cell := range row {
			fmt.Printf("%d ", cell)
		}
		fmt.Println()
	}
}

func (grid *Grid) DropPiece(column int, player int) bool {
	for i := len(grid.Grid) - 1; i >= 0; i-- {
		if grid.Grid[i][column] == 0 {
			grid.Grid[i][column] = player
			return true
		}
	}
	return false
}

func (grid *Grid) CheckWin(player int, line int, column int) bool {
	// Horizontal check
	count := 0
	for index := range 7 {
		if grid.Grid[line][index] == player {
			count++
			if count == 4 {
				return true
			}
		} else {
			count = 0
		}
	}
	// Vertical check
	count = 0
	for index := range 6 {
		if grid.Grid[index][column] == player {
			count++
			if count == 4 {
				return true
			}
		} else {
			count = 0
		}
	}
	// Diagonal check
	count = 0
	for index := -3; index <= 3; index++ {
		if line+index >= 0 && line+index < 6 && column+index >= 0 && column+index < 7 {
			if grid.Grid[line+index][column+index] == player {
				count++
				if count == 4 {
					return true
				}
			} else {
				count = 0
			}
		}
	}
	// Anti-diagonal check
	count = 0
	for index := -3; index <= 3; index++ {
		if line+index >= 0 && line+index < 6 && column-index >= 0 && column-index < 7 {
			if grid.Grid[line+index][column-index] == player {
				count++
				if count == 4 {
					return true
				}
			} else {
				count = 0
			}
		}
	}
	return false
}
