package grid

import (
	"fmt"
	"strconv"
)

const (
	WIDTH     = 7
	HEIGHT    = 6
	MIN_SCORE = -(WIDTH*HEIGHT)/2 + 3
	MAX_SCORE = (WIDTH*HEIGHT+1)/2 - 3
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
		if column < 0 || column >= WIDTH || !grid.CanPlay(column) || grid.IsWinningMove(column) {
			panic(fmt.Sprintf("Can't play at column %d", column+1))
		}
		grid.Play(column)
	}
	return grid
}

func (grid *Grid) CanPlay(column int) bool {
	return (grid.Mask & topMask(column)) == 0
}

func (grid *Grid) Play(column int) {
	grid.CurrentPosition ^= grid.Mask
	grid.Mask |= grid.Mask + bottomMask(column)
	grid.nbMoves++
}

func (grid *Grid) IsWinningMove(column int) bool {
	position := grid.CurrentPosition
	position |= (grid.Mask + bottomMask(column)) & columnMask(column)
	return CheckWin(position)
}

func (grid *Grid) IsDraw() bool {
	return grid.nbMoves == WIDTH*HEIGHT
}

func (grid *Grid) Key() uint64 {
	return grid.CurrentPosition + grid.Mask
}

func CheckWin(position uint64) bool {
	// Horizontal
	mask := position & (position >> (HEIGHT + 1))
	if mask&(mask>>(2*(HEIGHT+1))) != 0 {
		return true
	}

	// Vertical
	mask = position & (position >> 1)
	if mask&(mask>>2) != 0 {
		return true
	}

	// Diagonal 1 (\)
	mask = position & (position >> HEIGHT)
	if mask&(mask>>(2*HEIGHT)) != 0 {
		return true
	}

	// Anti-Diagonal (/)
	mask = position & (position >> (HEIGHT + 2))
	return mask&(mask>>(2*(HEIGHT+2))) != 0
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
