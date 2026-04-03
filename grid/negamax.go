//go:generate go run ../cmd/export/main.go
package grid

import (
	"time"

	"github.com/Pietot/Gonnect-4/config"
	"github.com/Pietot/Gonnect-4/database"
	"github.com/Pietot/Gonnect-4/evaluation"
	"github.com/Pietot/Gonnect-4/stats"
	"github.com/Pietot/Gonnect-4/transpositiontable"
	"github.com/Pietot/Gonnect-4/utils"
)

var (
	columnOrder = [7]int{3, 4, 2, 5, 1, 6, 0}
	nodeCount   = uint64(0)
	transTable  = transpositiontable.NewTranspositionTable()
)

// Gets the score of the current position using a negamax algorithm with alpha-beta pruning and a transposition table.
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

// Solves the current position by returning its score, the remaining moves until a forced win and some statistics about the solving process.
func (grid *Grid) Solve() (evaluation.Evaluation, stats.Stats) {
	if config.IsBookEnabled {
		scores, found := utils.GetScores(&database.ExportedBook, grid.Key(), grid.MirrorKey())
		if found {
			score, _ := utils.GetBestScoreAndMove(scores)
			return evaluation.Evaluation{
					Score:          score,
					RemainingMoves: GetRemainingMoves(score, grid.nbMoves),
				}, stats.Stats{
					TotalTimeNanoseconds: 0,
					NodeCount:            0,
					MeanTimePerNode:      0,
					NodesPerSecond:       0,
				}
		}
	}

	start := time.Now()
	score := grid.GetScore()
	elapsed := time.Since(start)

	elapsedNanoseconds := elapsed.Nanoseconds()
	nodesPerSecond := uint64(0)
	meanTimePerNode := 0.0
	if nodeCount > 0 {
		meanTimePerNode = float64(elapsedNanoseconds) / float64(nodeCount)
	}
	if elapsedNanoseconds > 0 {
		nodesPerSecond = uint64(float64(nodeCount) * 1_000_000_000 / float64(elapsedNanoseconds))
	}

	stats := stats.Stats{
		TotalTimeNanoseconds: elapsedNanoseconds,
		NodeCount:            nodeCount,
		MeanTimePerNode:      meanTimePerNode,
		NodesPerSecond:       nodesPerSecond,
	}

	nodeCount = 0

	return evaluation.Evaluation{
		Score:          score,
		RemainingMoves: GetRemainingMoves(score, grid.nbMoves),
	}, stats
}

// Analyzes the current position by returning the score for all possible moves, the remaining moves until a forced win and the best move to play, as well as some statistics about the analyzing process.
func (grid *Grid) Analyze() (evaluation.Analysis, stats.Stats) {
	scores := evaluation.Analysis{}
	bestMove := uint8(0)
	maxScore := int8(-128)

	if config.IsBookEnabled {
		sc, found := utils.GetScores(&database.ExportedBook, grid.Key(), grid.MirrorKey())
		if found {
			scores.Scores = sc
			_, bestMove = utils.GetBestScoreAndMove(sc)
			scores.BestMove = bestMove
			scores.RemainingMoves = GetRemainingMoves(sc[bestMove], grid.nbMoves)
			return scores, stats.Stats{
				TotalTimeNanoseconds: 0,
				NodeCount:            0,
				MeanTimePerNode:      0,
				NodesPerSecond:       0,
			}
		}
	}

	start := time.Now()
	for column := range 7 {
		var score int8
		if grid.CanPlay(column) {
			if grid.IsWinningMove(column) {
				score = int8((WIDTH*HEIGHT + 1 - grid.nbMoves) / 2)
			} else {
				childGrid := *grid
				childGrid.PlayColumn(column)
				score = -childGrid.GetScore()
			}
			scores.Scores[column] = score
			if score > maxScore {
				maxScore = score
				bestMove = uint8(column)
			}
		} else {
			scores.Scores[column] = config.NIL_SCORE
		}
	}

	elapsed := time.Since(start)
	elapsedNanoseconds := elapsed.Nanoseconds()
	nodesPerSecond := uint64(0)
	meanTimePerNode := 0.0
	if nodeCount > 0 {
		meanTimePerNode = float64(elapsedNanoseconds) / float64(nodeCount)
	}
	if elapsedNanoseconds > 0 {
		nodesPerSecond = uint64(float64(nodeCount) * 1_000_000_000 / float64(elapsedNanoseconds))
	}

	stats := stats.Stats{
		TotalTimeNanoseconds: elapsedNanoseconds,
		NodeCount:            nodeCount,
		MeanTimePerNode:      meanTimePerNode,
		NodesPerSecond:       nodesPerSecond,
	}

	nodeCount = 0

	bestRemainingMoves := GetRemainingMoves(maxScore, grid.nbMoves)

	scores.RemainingMoves = bestRemainingMoves
	scores.BestMove = bestMove
	return scores, stats
}

// Negamax algorithm with alpha-beta pruning transposition table and move ordering.
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
	if beta > max {
		beta = max
		if alpha >= beta {
			return int8(beta)
		}
	}

	key := grid.Key()
	if value := transTable.Get(key); value > 0 {
		if int(value) > MAX_SCORE-MIN_SCORE+1 {
			min := int8(int(value) + 2*MIN_SCORE - MAX_SCORE - 2)
			if alpha < min {
				alpha = min
				if alpha >= beta {
					return int8(alpha)
				}
			}
		} else {
			max := int8(int(value) + MIN_SCORE - 1)
			if beta > max {
				beta = max
				if alpha >= beta {
					return int8(beta)
				}
			}
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
			transTable.Put(key, uint8(int(childGridScore)+MAX_SCORE-2*MIN_SCORE+2))
			return childGridScore
		}
		if childGridScore > alpha {
			alpha = childGridScore
		}
	}

	transTable.Put(key, uint8(int(alpha)-MIN_SCORE+1))

	return alpha
}

// Returns the number of remaining moves until a forced win for the current position, given its score and the number of moves played so far.
func GetRemainingMoves(score int8, nbMoves int) uint8 {
	if score > 0 {
		return uint8((45-nbMoves)/2 - int(score))
	}
	if score < 0 {
		return uint8((44-nbMoves)/2 + int(score))
	}
	return uint8((WIDTH*HEIGHT - nbMoves) / 2)
}
