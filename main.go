package main

import (
	"fmt"

	"github.com/Pietot/Gonnect-4/grid"
)

func main() {
	testGrid := grid.InitGrid()
	testGrid.Grid = [][]int{
		{0, 2, 2, 2, 1, 0, 0},
		{0, 1, 1, 1, 2, 2, 2},
		{0, 1, 2, 2, 1, 1, 1},
		{1, 2, 1, 1, 2, 2, 2},
		{2, 2, 1, 2, 1, 2, 1},
		{1, 2, 1, 1, 2, 1, 1},
	}
	fmt.Println(testGrid.Negamax(2))
}
