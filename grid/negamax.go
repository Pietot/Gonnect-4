package grid

import (
	"math"

	"github.com/Pietot/Gonnect-4/evaluation"
	"github.com/Pietot/Gonnect-4/utils"
)

const COMPUTER = 1

var stack = 0

func (grid *Grid) Negamax(player int) *evaluation.Evaluation {
	if grid.IsDraw() {
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
