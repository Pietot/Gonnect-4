package grid

import (
	"math"
	"time"

	"github.com/Pietot/Gonnect-4/evaluation"
	"github.com/Pietot/Gonnect-4/stats"
	"github.com/Pietot/Gonnect-4/utils"
)

var movePlayed = 0

func (grid *Grid) Negamax(player int) (*evaluation.Evaluation, *stats.Stats) {
	start := time.Now()
	nbPos := int64(0)

	// Strong solver, use alpha = -1 and beta = 1 for a weaker solver
	result := grid.negamaxStats(player, &nbPos, math.Inf(-1), math.Inf(1))

	elapsed := time.Since(start)
	elapsedSeconds := elapsed.Seconds()
	meanNbPos := float64(nbPos)
	nbPosPerSec := 0.0
	meanTimePerPos := 0.0
	if meanNbPos > 0 {
		meanTimePerPos = (elapsed.Seconds() * 1000) / meanNbPos
	}
	if elapsedSeconds > 0 {
		nbPosPerSec = meanNbPos / elapsedSeconds
	}

	stats := &stats.Stats{
		TotalTimeMs:         elapsed.Seconds() * 1000,
		NumPositions:        nbPos,
		MeanTimePerPosition: meanTimePerPos,
		PositionsPerSecond:  nbPosPerSec,
	}

	return result, stats
}

func (grid *Grid) negamaxStats(player int, nbPos *int64, alpha float64, beta float64) *evaluation.Evaluation {
	if grid.IsDraw() {
		*nbPos++
		return &evaluation.Evaluation{
			Score:         utils.Float64Ptr(0.0),
			BestMove:      nil,
			RemainingMove: utils.IntPtr(movePlayed),
		}
	}

	var bestEvaluation = &evaluation.Evaluation{Score: nil, BestMove: nil, RemainingMove: nil}

	copyGrid := grid.DeepCopy()
	for column := range 7 {
		droppedPiece, line := copyGrid.DropPiece(column, player)
		if !droppedPiece {
			continue
		}
		*nbPos++
		movePlayed++
		if copyGrid.CheckWinFromIndex(player, line, column) {
			movePlayed--
			return &evaluation.Evaluation{
				Score:         utils.Float64Ptr(math.Inf(1)),
				BestMove:      &column,
				RemainingMove: utils.IntPtr(movePlayed + 1),
			}
		}
		movePlayed--
		copyGrid.Grid[line][column] = 0
	}

	max := float64((6*7 - 1 - grid.nbMoves)) / 2.0

	if beta > max {
		beta = max
		if alpha >= beta {
			return &evaluation.Evaluation{
				Score:         utils.Float64Ptr(beta),
				BestMove:      nil,
				RemainingMove: nil,
			}
		}
	}

	copyGrid = grid.DeepCopy()
	for column := range 7 {
		droppedPiece, line := copyGrid.DropPiece(column, player)
		if droppedPiece {
			movePlayed++
			childEvaluation := copyGrid.negamaxStats(
				utils.GetOpponent(player),
				nbPos,
				alpha,
				beta).Negate()
			movePlayed--
			*nbPos++
			newEvaluation := &evaluation.Evaluation{
				Score:         childEvaluation.Score,
				BestMove:      &column,
				RemainingMove: childEvaluation.RemainingMove,
			}
			if newEvaluation.IsBetterThan(bestEvaluation) {
				bestEvaluation = newEvaluation
			}
			copyGrid.Grid[line][column] = 0
		}
		if bestEvaluation.Score != nil && *bestEvaluation.Score >= beta {
			return bestEvaluation
		}
		if bestEvaluation.Score != nil && *bestEvaluation.Score > alpha {
			alpha = *bestEvaluation.Score
		}
	}

	return &evaluation.Evaluation{
		Score:         utils.Float64Ptr(alpha),
		BestMove:      bestEvaluation.BestMove,
		RemainingMove: bestEvaluation.RemainingMove,
	}
}
