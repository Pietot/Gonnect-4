package grid

import (
	"fmt"
	"strconv"
)

const (
	WIDTH      int    = 7
	HEIGHT     int    = 6
	MIN_SCORE  int    = -(WIDTH*HEIGHT)/2 + 3
	MAX_SCORE  int    = (WIDTH*HEIGHT+1)/2 - 3
	BOTTOM     uint64 = 0b0000001000000100000010000001000000100000010000001
	BOARD_MASK uint64 = BOTTOM * ((1 << HEIGHT) - 1)
)

type Grid struct {
	CurrentPosition uint64
	Mask            uint64
	nbMoves         int
}

func InitGrid(columnsSequence string) *Grid {
	grid := &Grid{}
	for _, columnRune := range columnsSequence {
		column, err := strconv.Atoi(string(columnRune))
		if err != nil {
			panic(fmt.Sprintf("Invalid column character: %v", err))
		}
		column -= 1
		if column < 0 || column >= WIDTH || !grid.canPlay(column) || grid.IsWinningMove(column) {
			panic(fmt.Sprintf("Can't play at column %d", column+1))
		}
		grid.play(column)
	}
	return grid
}

func (grid *Grid) Key() uint64 {
	return grid.CurrentPosition + grid.Mask + BOTTOM
}

func (grid *Grid) canPlay(column int) bool {
	return (grid.Mask & topMask(column)) == 0
}

func (grid *Grid) play(column int) {
	grid.CurrentPosition ^= grid.Mask
	grid.Mask |= grid.Mask + bottomMask(column)
	grid.nbMoves++
}

func (grid *Grid) IsWinningMove(column int) bool {
	return (grid.winningPositionMask() & grid.possibleMask() & columnMask(column)) != 0
}

func (grid *Grid) canWinNext() bool {
	return (grid.winningPositionMask() & grid.possibleMask()) != 0
}

func (grid *Grid) isDraw() bool {
	return grid.nbMoves >= WIDTH*HEIGHT-2
}

func (grid *Grid) findNextWinningMove() *int {
	winningMask := grid.winningPositionMask()
	possibleMask := grid.possibleMask()
	winningMoves := winningMask & possibleMask

	for _, column := range columnOrder {
		columnMask := columnMask(column)
		if (winningMoves & columnMask) != 0 {
			return &column
		}
	}

	panic("No winning move found, but expected one")
}

func (grid *Grid) getRandomColumn() *int {
	for _, column := range columnOrder {
		if grid.canPlay(column) {
			return &column
		}
	}
	panic("No playable column found")
}

func (grid *Grid) possibleNonLoosingMoves() uint64 {
	possible_mask := grid.possibleMask()
	opponent_win := grid.opponentWinningPositionMask()
	forced_moves := possible_mask & opponent_win
	if forced_moves != 0 {
		if (forced_moves & (forced_moves - 1)) != 0 {
			return 0
		}
		possible_mask = forced_moves
	}
	return possible_mask & ^(opponent_win >> 1)
}

func (grid *Grid) possibleMask() uint64 {
	return (grid.Mask + BOTTOM) & BOARD_MASK
}

func (grid *Grid) winningPositionMask() uint64 {
	return computeWinningPosition(grid.CurrentPosition, grid.Mask)
}

func (grid *Grid) opponentWinningPositionMask() uint64 {
	return computeWinningPosition(grid.CurrentPosition^grid.Mask, grid.Mask)
}

func topMask(column int) uint64 {
	return (uint64(1) << (HEIGHT - 1)) << (column * (HEIGHT + 1))
}

func bottomMask(column int) uint64 {
	return uint64(1) << (column * (HEIGHT + 1))
}

func columnMask(column int) uint64 {
	return ((uint64(1) << HEIGHT) - 1) << (column * (HEIGHT + 1))
}

func computeWinningPosition(position uint64, mask uint64) uint64 {
	// Vertical
	r := (position << 1) & (position << 2) & (position << 3)

	uint_height := uint64(HEIGHT)

	// Horizontal
	p := (position << (uint_height + 1)) & (position << (2 * (uint_height + 1)))
	r |= p & (position << (3 * (uint_height + 1)))
	r |= p & (position >> (uint_height + 1))
	p >>= (3 * (uint_height + 1))
	r |= p & (position << (uint_height + 1))
	r |= p & (position >> (3 * (uint_height + 1)))

	// Diagonal (\)
	p = (position << uint_height) & (position << (2 * uint_height))
	r |= p & (position << (3 * uint_height))
	r |= p & (position >> uint_height)
	p >>= (3 * uint_height)
	r |= p & (position << uint_height)
	r |= p & (position >> (3 * uint_height))

	// Anti-Diagonal (/)
	p = (position << (uint_height + 2)) & (position << (2 * (uint_height + 2)))
	r |= p & (position << (3 * (uint_height + 2)))
	r |= p & (position >> (uint_height + 2))
	p >>= (3 * (uint_height + 2))
	r |= p & (position << (uint_height + 2))
	r |= p & (position >> (3 * (uint_height + 2)))

	return r & (BOARD_MASK ^ mask)
}
