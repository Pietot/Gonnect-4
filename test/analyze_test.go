package test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/Pietot/Gonnect-4/grid"
	"github.com/Pietot/Gonnect-4/utils"
)

var files = []string{
	"../data/test_easy_end.txt",
	"../data/test_easy_middle.txt",
	"../data/test_easy_begin.txt",
	"../data/test_medium_end.txt",
	"../data/test_medium_middle.txt",
	"../data/test_hard_begin.txt",
}

// You don't need to test all positions from all files every time cause it takes long.
// You can only ensure the firsts files are correct and then stop the tests.
func TestAnalyze(t *testing.T) {
	for fileIndex, file := range files {
		lines, err := utils.ReadPositionsFromFile(file)
		if err != nil {
			t.Fatalf("Error reading file: %v", err)
		}
		for _, line := range lines {
			position, expectedScore, expectedColumn := strings.Split(line, " ")[0], strings.Split(line, " ")[1], strings.TrimSpace(strings.Split(line, " ")[2])[1:]
			game, err := grid.InitGrid(position)
			if err != nil {
				t.Fatalf("Error initializing grid: %v", err)
			}
			anal, _ := game.Analyze()
			analScore, analColumn := getFoundValues(anal.Scores)
			expectedScoreInt, err := strconv.ParseInt(expectedScore, 10, 8)
			if err != nil {
				t.Errorf("Error parsing expected score: %v", err)
				continue
			}
			expectedColumnUint, err := strconv.ParseUint(expectedColumn, 10, 8)
			if err != nil {
				t.Errorf("Error parsing expected column: %v", err)
				continue
			}
			if analScore != int8(expectedScoreInt) || analColumn+1 != uint8(expectedColumnUint) {
				t.Errorf("Discrepancy found for position %s: expected score %s and best move %s, got score %d and best move %d",
					position, expectedScore, expectedColumn, analScore, analColumn)
			} else {
				t.Logf("[File %d] Position %s analyzed correctly", fileIndex, position)
			}
		}
	}
}

func getFoundValues(scores [7]*int8) (bestScore int8, bestMove uint8) {
	bestScore = -128
	bestMove = 0
	for i, scorePtr := range scores {
		if scorePtr != nil && *scorePtr > bestScore {
			bestScore = *scorePtr
			bestMove = uint8(i)
		}
	}
	return bestScore, bestMove
}
