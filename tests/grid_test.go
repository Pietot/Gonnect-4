package test

import (
	"testing"

	"github.com/Pietot/Gonnect-4/grid"
)

func TestConnect4(t *testing.T) {
	test_game := initGame()
	for i, game := range test_game {
		if !game.Grid.CheckWinFromIndex(game.player, game.line, game.column) {
			t.Errorf("Expected player %d to win in game %d", game.player, i)
		}
	}
}

func TestGridInitialization(t *testing.T) {
	test_grid := grid.InitGrid()
	if len(test_grid.Grid) != 6 || len(test_grid.Grid[0]) != 7 {
		t.Errorf("Expected grid to be 6 rows and 7 columns, got %d rows and %d columns", len(test_grid.Grid), len(test_grid.Grid[0]))
	}
	// Check if all cells are initialized to 0
	for i := range test_grid.Grid {
		for j := range test_grid.Grid[i] {
			if test_grid.Grid[i][j] != 0 {
				t.Errorf("Expected cell (%d, %d) to be 0, got %d", i, j, test_grid.Grid[i][j])
			}
		}
	}
}

type game struct {
	Grid   *grid.Grid
	player int
	line   int
	column int
}

func initGrid(custom_grid [][]int) *grid.Grid {
	return &grid.Grid{
		Grid: custom_grid,
	}
}

func initGame() []game {
	return []game{
		{Grid: initGrid([][]int{
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{2, 0, 0, 0, 0, 0, 0},
			{2, 0, 0, 0, 0, 0, 0},
			{2, 1, 1, 1, 1, 0, 0},
		}), player: 1, line: 5, column: 1},
		{Grid: initGrid([][]int{
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{2, 0, 0, 0, 0, 0, 0},
			{2, 0, 0, 0, 0, 0, 0},
			{2, 1, 1, 1, 1, 0, 0},
		}), player: 1, line: 5, column: 4},
		{Grid: initGrid([][]int{
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{1, 0, 0, 0, 0, 0, 0},
			{1, 0, 0, 0, 0, 0, 0},
			{1, 0, 0, 0, 0, 0, 0},
			{1, 0, 2, 2, 2, 0, 0},
		}), player: 1, line: 2, column: 0},
		{Grid: initGrid([][]int{
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 0},
			{0, 0, 0, 0, 1, 2, 0},
			{0, 0, 0, 1, 2, 1, 0},
			{0, 0, 1, 2, 2, 2, 0},
		}), player: 1, line: 2, column: 5},
		{Grid: initGrid([][]int{
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 0},
			{0, 0, 0, 0, 1, 2, 0},
			{0, 0, 0, 1, 2, 2, 0},
			{0, 0, 1, 2, 2, 2, 0},
		}), player: 1, line: 5, column: 2},
		{Grid: initGrid([][]int{
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0, 0, 0},
			{0, 0, 2, 1, 0, 0, 0},
			{0, 0, 2, 2, 1, 0, 0},
			{0, 0, 2, 2, 2, 1, 0},
		}), player: 1, line: 2, column: 2},
		{Grid: initGrid([][]int{
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0, 0, 0},
			{0, 0, 2, 1, 0, 0, 0},
			{0, 0, 2, 2, 1, 0, 0},
			{0, 0, 2, 2, 2, 1, 0},
		}), player: 1, line: 5, column: 5},
	}
}
