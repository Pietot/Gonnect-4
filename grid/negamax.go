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

func (grid *Grid) GetScore() int8 {
	// Strong solver, use min = -1 and max = 1 for a weak solver
	minScore := int8(-(WIDTH*HEIGHT - grid.nbMoves) / 2)
	maxScore := int8((WIDTH*HEIGHT + 1 - grid.nbMoves) / 2)

	if grid.canWinNext() {
		return maxScore
	}
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
	return minScore
}

func (grid *Grid) Solve() (evaluation.Evaluation, stats.Stats) {
	start := time.Now()

	score := grid.GetScore()

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

	stats := stats.Stats{
		TotalTimeMicroseconds: elapsed.Seconds() * 1_000_000,
		NodeCount:             nodeCount,
		MeanTimePerNode:       meanTimePerNode,
		NodesPerSecond:        nodesPerSecond,
	}

	nodeCount = 0
	
	return evaluation.Evaluation{
		Score:          &score,
		RemainingMoves: GetRemainingMoves(score, grid.nbMoves),
	}, stats
}

func (grid *Grid) Analyze() evaluation.Analyzation {
	scores := evaluation.Analyzation{}
	bestMove := uint8(0)
	maxScore := int8(-128)
	for column := range 7 {
		if grid.canPlay(column) {
			var score int8
			if grid.IsWinningMove(column) {
				score = int8((WIDTH*HEIGHT + 1 - grid.nbMoves) / 2)
			} else {
				childGrid := *grid
				childGrid.playColumn(column)
				score = -childGrid.GetScore()
			}
			scores.Scores[column] = &score
			if score > maxScore {
				maxScore = score
				bestMove = uint8(column)
			}
		}
	}

	bestRemainingMoves := GetRemainingMoves(maxScore, grid.nbMoves)

	scores.RemainingMoves = bestRemainingMoves
	scores.BestMove = &bestMove
	return scores
}

func (grid *Grid) negamax(alpha int8, beta int8) int8 {
	nodeCount++

	nextMoves := grid.possibleNonLosingMoves()
	if nextMoves == 0 {
		return -int8(((WIDTH * HEIGHT) - grid.nbMoves) / 2)
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

	moves := newMoveSorter()
	for i := WIDTH - 1; i >= 0; i-- {
		column := columnOrder[i]
		if move := nextMoves & columnMask(column); move != 0 {
			moves.addMove(move, grid.moveScore(move))
		}
	}

	for nextMove := moves.getNextMove(); nextMove != 0; nextMove = moves.getNextMove() {
		childGrid := *grid
		childGrid.play(nextMove)
		childGridScore := -childGrid.negamax(-beta, -alpha)
		if childGridScore >= beta {
			return childGridScore
		}
		if childGridScore > alpha {
			alpha = childGridScore
		}
	}

	trans_table.Put(grid.Key(), uint8(int(alpha)-MIN_SCORE+1))

	return alpha
}

func GetRemainingMoves(score int8, nbMoves int) *uint8 {
	if score > 0 {
		return utils.Uint8Ptr(uint8((45-nbMoves)/2 - int(score)))
	}
	if score < 0 {
		return utils.Uint8Ptr(uint8((44-nbMoves)/2 + int(score)))
	}
	return utils.Uint8Ptr(uint8((WIDTH*HEIGHT - nbMoves) / 2))
}
