package test

import (
	"testing"

	"github.com/Pietot/Gonnect-4/grid"
)

func TestConnect4(t *testing.T) {
	test_games := initGames()
	for i, game := range test_games {
		if !game.Grid.IsWinningMove(game.column) {
			t.Errorf("Expected player %d to win in game %d", game.player, i)
		}
	}
}

type game struct {
	Grid   *grid.Grid
	player int
	column int
}

func initGames() []game {
	return []game{
		{
			Grid:   grid.InitGrid("213141"),
			player: 1,
			column: 4,
		},
		{
			Grid:   grid.InitGrid("131415"),
			player: 1,
			column: 0,
		},
		{
			Grid:   grid.InitGrid("3446666575"),
			player: 1,
			column: 4,
		},
		{
			Grid:   grid.InitGrid("7554644333"),
			player: 1,
			column: 2,
		},
	}
}
