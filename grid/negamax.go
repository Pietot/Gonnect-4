package grid

import (
	"math"
	"time"

	"github.com/Pietot/Gonnect-4/evaluation"
	"github.com/Pietot/Gonnect-4/stats"
	"github.com/Pietot/Gonnect-4/utils"
)

var stack = 0

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
			RemainingMove: utils.IntPtr(stack),
		}
	}

	var bestEvaluation = &evaluation.Evaluation{Score: nil, BestMove: nil, RemainingMove: nil}

	for column := range 7 {
		copyGrid := grid.DeepCopy()
		droppedPiece, line := copyGrid.DropPiece(column, player)
		if droppedPiece && copyGrid.CheckWinFromIndex(player, line, column) {
			*nbPos++
			return &evaluation.Evaluation{
				Score:         utils.Float64Ptr(math.Inf(1)),
				BestMove:      &column,
				RemainingMove: utils.IntPtr(stack),
			}
		}
	}

	for column := range 7 {
		copyGrid := grid.DeepCopy()
		droppedPiece, _ := copyGrid.DropPiece(column, player)
		if droppedPiece {
			stack++
			childEvaluation := copyGrid.negamaxStats(utils.GetOpponent(player), nbPos).Negate()
			stack--
			newEvaluation := &evaluation.Evaluation{
				Score:         childEvaluation.Score,
				BestMove:      &column,
				RemainingMove: childEvaluation.RemainingMove,
			}
			if newEvaluation.IsBetterThan(bestEvaluation) {
				bestEvaluation = newEvaluation
			}
		}
	}

	return bestEvaluation
}
