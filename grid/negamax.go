package grid

import (
	"time"

	"github.com/Pietot/Gonnect-4/evaluation"
	"github.com/Pietot/Gonnect-4/stats"
	"github.com/Pietot/Gonnect-4/transposition_table"
	"github.com/Pietot/Gonnect-4/utils"
)

var (
	columnOrder = [7]int{3, 4, 2, 5, 1, 6, 0}
	nodeCount   = uint64(0)
	trans_table = transposition_table.NewTranspositionTable()
)

func (grid *Grid) Solve() (*evaluation.Evaluation, *stats.Stats) {
	// Strong solver, use min = -1 and max = 1 for a weak solver
	minScore := int8(-(WIDTH*HEIGHT - grid.nbMoves) / 2)
	maxScore := int8((WIDTH*HEIGHT + 1 - grid.nbMoves) / 2)
	evaluation := &evaluation.Evaluation{}

	start := time.Now()

	if grid.canWinNext() {
		evaluation.Score = &maxScore
		evaluation.RemainingMoves = utils.Uint8Ptr(1)
	} else {
		for minScore < maxScore {
			middle := minScore + (maxScore-minScore)/2
			if middle <= 0 && minScore/2 < middle {
				middle = minScore / 2
			} else if middle >= 0 && maxScore/2 > middle {
				middle = maxScore / 2
			}
			result := grid.negamax(middle, middle+1)
			if result <= middle {
				maxScore = result
			} else {
				minScore = result
			}
		}

	}

	elapsed := time.Since(start)
	elapsedSeconds := elapsed.Seconds()
	nodesPerSecond := uint64(0)
	meanTimePerNode := 0.0
	if nodeCount > 0 {
		meanTimePerNode = (elapsed.Seconds() * 1_000_000) / float64(nodeCount)
	}
	if elapsedSeconds > 0 {
		nodesPerSecond = uint64(float64(nodeCount) / elapsedSeconds)
	}

	stats := &stats.Stats{
		TotalTimeMicroseconds: elapsed.Seconds() * 1_000_000,
		NodeCount:             nodeCount,
		MeanTimePerNode:       meanTimePerNode,
		NodesPerSecond:        nodesPerSecond,
	}

	evaluation.Score = &minScore
	evaluation.RemainingMoves = getRemainingMoves(minScore, grid.nbMoves)
	return evaluation, stats
}

func (grid *Grid) negamax(alpha int8, beta int8) int8 {
	nodeCount++

	nextMoves := grid.possibleNonLoosingMoves()
	if nextMoves == 0 {
		return -int8((WIDTH * HEIGHT) - grid.nbMoves/2)
	}

	if grid.isDraw() {
		return 0
	}

	min := int8(-(WIDTH*HEIGHT - 2 - grid.nbMoves) / 2)
	if alpha < min {
		alpha = min
		if alpha >= beta {
			return int8(alpha)
		}
	}

	max := int8((WIDTH*HEIGHT - 1 - grid.nbMoves) / 2)
	if value := trans_table.Get(grid.Key()); value > 0 {
		max = int8(int(value) + MIN_SCORE - 1)
	}
	if beta > max {
		beta = max
		if alpha >= beta {
			return int8(beta)
		}
	}

	for _, column := range columnOrder {
		if (nextMoves & columnMask(column)) != 0 {
			childGrid := *grid
			childGrid.play(column)
			childGridScore := -childGrid.negamax(-beta, -alpha)
			if childGridScore >= beta {
				return childGridScore
			}
			if childGridScore > alpha {
				alpha = childGridScore
			}
		}
	}

	trans_table.Put(grid.Key(), uint8(int(alpha)-MIN_SCORE+1))

	return alpha
}

func getRemainingMoves(score int8, nbMoves int) *uint8 {
	if score > 0 {
		return utils.Uint8Ptr(uint8((45-nbMoves)/2 - int(score)))
	}
	if score < 0 {
		return utils.Uint8Ptr(uint8((44-nbMoves)/2 + int(score)))
	}
	return utils.Uint8Ptr(uint8((WIDTH*HEIGHT - nbMoves) / 2))
}
