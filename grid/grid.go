package grid

type Grid struct {
	Grid    [][]int
	nbMoves int
}

func InitGrid() *Grid {
	grid := &Grid{
		Grid: make([][]int, 6),
	}
	for i := range grid.Grid {
		grid.Grid[i] = make([]int, 7)
	}
	grid.nbMoves = 0
	return grid
}

func (grid *Grid) String() string {
	var output string
	for _, row := range grid.Grid {
		for _, cell := range row {
			switch cell {
			case 1:
				// Color red for player 1
				output += "\033[31m@\033[0m "
			case 2:
				// Color blue for player 2
				output += "\033[34m@\033[0m "
			default:
				// Empty cell
				output += "  "
			}
		}
		output += "\n"
	}
	output += "-------------\n"
	output += "1 2 3 4 5 6 7\n"
	return output
}

func (grid *Grid) DropPiece(column int, player int) (bool, int) {
	for i := len(grid.Grid) - 1; i >= 0; i-- {
		if grid.Grid[i][column] == 0 {
			grid.Grid[i][column] = player
			return true, i
		}
	}
	return false, -1
}

func (grid *Grid) CheckWinFromIndex(player int, line int, column int) bool {
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

func (grid *Grid) IsDraw() bool {
	// For debugging purposes when I directly initialize the grid
	if grid.nbMoves == 0 {
		count := 0
		for _, row := range grid.Grid {
			for _, cell := range row {
				if cell != 0 {
					count++
				}
			}
		}
		return count == 6*7
	}
	return grid.nbMoves == 6*7
}

func (grid *Grid) DeepCopy() *Grid {
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
