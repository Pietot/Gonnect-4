package grid

import (
	"math"
	"time"

	"github.com/Pietot/Gonnect-4/evaluation"
	"github.com/Pietot/Gonnect-4/stats"
	"github.com/Pietot/Gonnect-4/utils"
)

var movedPlayed = 0

func (grid *Grid) Negamax(player int) (*evaluation.Evaluation, *stats.Stats) {
	start := time.Now()
	nbPos := int64(0)

	result := grid.negamaxStats(player, &nbPos)

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

func (grid *Grid) negamaxStats(player int, nbPos *int64) *evaluation.Evaluation {
	if grid.IsDraw() {
		*nbPos++
		return &evaluation.Evaluation{
			Score:         utils.Float64Ptr(0.0),
			BestMove:      nil,
			RemainingMove: utils.IntPtr(movedPlayed),
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
		movedPlayed++
		if copyGrid.CheckWinFromIndex(player, line, column) {
			return &evaluation.Evaluation{
				Score:         utils.Float64Ptr(math.Inf(1)),
				BestMove:      &column,
				RemainingMove: utils.IntPtr(movedPlayed),
			}
		}
		movedPlayed--
		copyGrid.Grid[line][column] = 0
	}

	copyGrid = grid.DeepCopy()
	for column := range 7 {
		droppedPiece, line := copyGrid.DropPiece(column, player)
		if droppedPiece {
			movedPlayed++
			childEvaluation := copyGrid.negamaxStats(utils.GetOpponent(player), nbPos).Negate()
			movedPlayed--
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
	}

	return bestEvaluation
}
