package grid

import (
	"fmt"
	"math/bits"
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

func InitGrid(columnsSequence string) (*Grid, error) {
	grid := &Grid{}
	for _, columnRune := range columnsSequence {
		column, err := strconv.Atoi(string(columnRune))
		if err != nil {
			return nil, fmt.Errorf("invalid column character: %v", err)
		}
		column -= 1
		if column < 0 || column >= WIDTH || !grid.CanPlay(column) || grid.IsWinningMove(column) {
			return nil, fmt.Errorf("can't play at column %d", column+1)
		}
		grid.PlayColumn(column)
	}
	return grid, nil
}

func (grid *Grid) Key() uint64 {
	return grid.CurrentPosition + grid.Mask + BOTTOM
}

func (grid *Grid) MirrorKey() uint64 {
	return mirror(grid.CurrentPosition) + mirror(grid.Mask) + BOTTOM
}

func (grid *Grid) GetCanonicalKey() uint64 {
	return min(grid.Key(), grid.MirrorKey())
}

// mirror return the bitboard reversed horizontally
func mirror(bitboard uint64) uint64 {
	return ((bitboard & 0x7F) << 42) |
		((bitboard & (0x7F << 7)) << 28) |
		((bitboard & (0x7F << 14)) << 14) |
		(bitboard & (0x7F << 21)) |
		((bitboard & (0x7F << 28)) >> 14) |
		((bitboard & (0x7F << 35)) >> 28) |
		((bitboard & (0x7F << 42)) >> 42)
}

func FromKey(key uint64) *Grid {
	g := &Grid{}
	g.Mask = 0
	g.CurrentPosition = 0
	for i := range 7 {
		colMask := uint64(0x7F) << (i * 7)
		colBits := key & colMask

		if colBits > 0 {
			msb := uint64(1) << (63 - bits.LeadingZeros64(colBits))
			bottomBit := uint64(1) << (i * 7)
			columnMaskInGrid := (msb - bottomBit) & colMask
			g.Mask |= columnMaskInGrid
			g.CurrentPosition |= (colBits ^ msb)
		}
	}
	g.nbMoves = bits.OnesCount64(g.Mask)

	return g
}

func (grid *Grid) IsWinningMove(column int) bool {
	return (grid.winningPositionMask() & grid.possibleMask() & columnMask(column)) != 0
}

func (grid *Grid) CanPlay(column int) bool {
	return (grid.Mask & topMask(column)) == 0
}

func (grid *Grid) play(move uint64) {
	grid.CurrentPosition ^= grid.Mask
	grid.Mask |= move
	grid.nbMoves++
}

func (grid *Grid) PlayColumn(column int) {
	grid.play((grid.Mask + bottomMask(column)) & columnMask(column))
}

func (grid *Grid) canWinNext() bool {
	return (grid.winningPositionMask() & grid.possibleMask()) != 0
}

func (grid *Grid) isDraw() bool {
	return grid.nbMoves >= WIDTH*HEIGHT-2
}

func (grid *Grid) possibleNonLosingMoves() uint64 {
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

func (grid *Grid) moveScore(move uint64) int {
	return popCount(computeWinningPosition(grid.CurrentPosition|move, grid.Mask))
}

func popCount(move uint64) int {
	count := 0
	for move != 0 {
		move &= move - 1
		count++
	}
	return count
}
