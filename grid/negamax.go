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
	nodeCount   = int64(0)
	trans_table = transposition_table.NewTranspositionTable()
)

func (grid *Grid) Solve() (*evaluation.Evaluation, *stats.Stats) {
	// Strong solver, use min = -1 and max = 1 for a weak solver
	minScore := int8(-(WIDTH*HEIGHT - grid.nbMoves) / 2)
	maxScore := int8((WIDTH*HEIGHT + 1 - grid.nbMoves) / 2)
	min := &evaluation.Evaluation{Score: &minScore}
	max := &evaluation.Evaluation{Score: &maxScore}

	start := time.Now()

	if grid.CanWinNext() {
		min = &evaluation.Evaluation{
			Score:          utils.Int8Ptr(maxScore),
			BestMove:       grid.FindNextWinningMove(),
			RemainingMoves: utils.Uint8Ptr(1),
		}
	} else {
		for *min.Score < *max.Score {
			middle := int8(*max.Score+*min.Score) / 2
			if middle <= 0 && *min.Score/2 < middle {
				middle = *min.Score / 2
			} else if middle >= 0 && *max.Score/2 > middle {
				middle = *max.Score / 2
			}
			result := grid.negamax(middle, middle+1)
			if *result.Score <= middle {
				max = result
			} else {
				min = result
			}
		}

	}

	elapsed := time.Since(start)
	elapsedSeconds := elapsed.Seconds()
	nodesPerSecond := int64(0)
	meanTimePerNode := 0.0
	if nodeCount > 0 {
		meanTimePerNode = (elapsed.Seconds() * 1_000_000) / float64(nodeCount)
	}
	if elapsedSeconds > 0 {
		nodesPerSecond = int64(float64(nodeCount) / elapsedSeconds)
	}

	stats := &stats.Stats{
		TotalTimeMicroseconds: elapsed.Seconds() * 1_000_000,
		NodeCount:             nodeCount,
		MeanTimePerNode:       meanTimePerNode,
		NodesPerSecond:        nodesPerSecond,
	}

	if min.RemainingMoves == nil {
		min = max
	}

	return min, stats
}

func (grid *Grid) negamax(alpha int8, beta int8) *evaluation.Evaluation {
	nodeCount++

	nextMoves := grid.possibleNonLoosingMoves()
	if nextMoves == 0 {
		return &evaluation.Evaluation{
			Score:          utils.Int8Ptr(int8(-(WIDTH*HEIGHT - grid.nbMoves) / 2)),
			BestMove:       grid.getRandomColumn(),
			RemainingMoves: utils.Uint8Ptr(2),
		}
	}

	if grid.IsDraw() {
		return &evaluation.Evaluation{
			Score:          utils.Int8Ptr(0),
			BestMove:       nil,
			RemainingMoves: utils.Uint8Ptr(1),
		}
	}

	min := int8(-(WIDTH*HEIGHT - 2 - grid.nbMoves) / 2)
	if alpha < min {
		alpha = min
		if alpha >= beta {
			return &evaluation.Evaluation{
				Score:          utils.Int8Ptr(alpha),
				BestMove:       nil,
				RemainingMoves: utils.Uint8Ptr(0),
			}
		}
	}

	max := int8((WIDTH*HEIGHT - 1 - grid.nbMoves) / 2)
	value, remaining, found := trans_table.Get(grid.Key())
	if found {
		max = int8(int(value) + MIN_SCORE - 1)
	}
	if beta > max {
		beta = max
		if alpha >= beta {
			return &evaluation.Evaluation{
				Score:          utils.Int8Ptr(beta),
				BestMove:       nil,
				RemainingMoves: utils.Uint8Ptr(remaining),
			}
		}
	}

	var bestMove *int
	var bestScore = alpha
	var bestRemainingMoves *uint8

	for _, column := range columnOrder {
		if (nextMoves & columnMask(column)) != 0 {
			childGrid := *grid
			childGrid.Play(column)
			childEvaluation := childGrid.negamax(-beta, -alpha).Negate()
			if *childEvaluation.Score >= beta {
				return &evaluation.Evaluation{
					Score:          childEvaluation.Score,
					BestMove:       &column,
					RemainingMoves: utils.Uint8Ptr(*childEvaluation.RemainingMoves + 1),
				}
			}
			if *childEvaluation.Score > bestScore || bestMove == nil {
				bestScore = *childEvaluation.Score
				bestMove = &column
				bestRemainingMoves = utils.Uint8Ptr(*childEvaluation.RemainingMoves + 1)
				alpha = bestScore
			}
		}
	}

	trans_table.Put(grid.Key(), uint8(int(alpha)-MIN_SCORE+1), uint8(*bestRemainingMoves))

	return &evaluation.Evaluation{
		Score:          utils.Int8Ptr(bestScore),
		BestMove:       bestMove,
		RemainingMoves: bestRemainingMoves,
	}
}
