package grid

import (
	"math"
	"time"

	"github.com/Pietot/Gonnect-4/evaluation"
	"github.com/Pietot/Gonnect-4/stats"
	"github.com/Pietot/Gonnect-4/utils"
)

var movePlayed = 0
var columnOrder = [7]int{3, 4, 2, 5, 1, 6, 0}
var nodeCount = int64(0)

func (grid *Grid) Negamax() (*evaluation.Evaluation, *stats.Stats) {
	start := time.Now()

	// Strong solver, use alpha = -1 and beta = 1 for a weak solver
	result := grid.negamaxStats(math.Inf(-1), math.Inf(1))

	elapsed := time.Since(start)
	elapsedSeconds := elapsed.Seconds()
	nbPosPerSec := 0.0
	meanTimePerPos := 0.0
	if nodeCount > 0 {
		meanTimePerPos = (elapsed.Seconds() * 1_000_000) / float64(nodeCount)
	}
	if elapsedSeconds > 0 {
		nbPosPerSec = float64(nodeCount) / elapsedSeconds
	}

	stats := &stats.Stats{
		TotalTimeMicroseconds: elapsed.Seconds() * 1_000_000,
		NodeCount:             nodeCount,
		MeanTimePerPosition:   meanTimePerPos,
		PositionsPerSecond:    nbPosPerSec,
	}

	return result, stats
}

func (grid *Grid) negamaxStats(alpha float64, beta float64) *evaluation.Evaluation {
	nodeCount++
	if grid.IsDraw() {
		return &evaluation.Evaluation{
			Score:         utils.Float64Ptr(0.0),
			BestMove:      nil,
			RemainingMove: utils.IntPtr(movePlayed),
		}
	}

	for column := range 7 {
		if grid.CanPlay(column) && grid.IsWinningMove(column) {
			return &evaluation.Evaluation{
				Score:         utils.Float64Ptr(float64(int(WIDTH*HEIGHT+1-grid.nbMoves) / 2)),
				BestMove:      &column,
				RemainingMove: utils.IntPtr(movePlayed + 1),
			}
		}
	}

	max := float64((WIDTH*HEIGHT - 1 - grid.nbMoves) / 2)

	if beta > max {
		beta = max
		if alpha >= beta {
			return &evaluation.Evaluation{
				Score:         utils.Float64Ptr(beta),
				BestMove:      nil,
				RemainingMove: utils.IntPtr(movePlayed + 1),
			}
		}
	}

	var bestMove *int
	var bestScore = alpha
	var bestRemainingMove *int

	for _, column := range columnOrder {
		if grid.CanPlay(column) {
			childGrid := *grid
			childGrid.Play(column)
			movePlayed++
			childEvaluation := childGrid.negamaxStats(
				-beta,
				-alpha).Negate()
			movePlayed--
			if *childEvaluation.Score >= beta {
				return &evaluation.Evaluation{
					Score:         childEvaluation.Score,
					BestMove:      &column,
					RemainingMove: childEvaluation.RemainingMove,
				}
			}
			if *childEvaluation.Score > bestScore || bestMove == nil {
				bestScore = *childEvaluation.Score
				bestMove = &column
				bestRemainingMove = childEvaluation.RemainingMove
				alpha = bestScore
			}
		}
	}

	return &evaluation.Evaluation{
		Score:         utils.Float64Ptr(bestScore),
		BestMove:      bestMove,
		RemainingMove: bestRemainingMove,
	}
}
