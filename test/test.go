package test

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Pietot/Gonnect-4/grid"
)

func TestAnalyze() {
	lines, err := readPositionsFromFile("data/test_easy_end.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
	for _, line := range lines {
		position, expectedScore, expectedColumn := strings.Split(line, " ")[0], strings.Split(line, " ")[1], strings.TrimSpace(strings.Split(line, " ")[2])[1:]
		game, err := grid.InitGrid(position)
		if err != nil {
			fmt.Println("Error initializing grid:", err)
		}
		anal := game.Analyze()
		analScore, analColumn := getFoundValues(anal.Scores)
		expectedScoreInt, err := strconv.ParseInt(expectedScore, 10, 8)
		if err != nil {
			fmt.Println("Error parsing expected score:", err)
			continue
		}
		expectedColumnUint, err := strconv.ParseUint(expectedColumn, 10, 8)
		if err != nil {
			fmt.Println("Error parsing expected column:", err)
			continue
		}
		if analScore != int8(expectedScoreInt) || analColumn+1 != uint8(expectedColumnUint) {
			panic(fmt.Sprintf("Discrepancy found for position %s: expected score %s and best move %s, got score %d and best move %d\n", position, expectedScore, expectedColumn, analScore, analColumn))
		} else {
			fmt.Println("Position", position, "analyzed correctly")
		}
	}
}

func readPositionsFromFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	return lines, nil
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
