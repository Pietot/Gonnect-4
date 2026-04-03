package test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Pietot/Gonnect-4/config"
	"github.com/Pietot/Gonnect-4/grid"
	"github.com/Pietot/Gonnect-4/utils"
)

var files = []string{
	"../data/test_easy_end.json",
	"../data/test_easy_middle.json",
	"../data/test_easy_begin.json",
	"../data/test_medium_end.json",
	"../data/test_medium_middle.json",
	"../data/test_hard_begin.json",
}

// You don't need to test all positions from all files every time cause it takes long.
//
// You can only ensure the firsts files are correct and then stop the tests.
func TestAnalyze(t *testing.T) {
	for fileIndex, file := range files {
		f, err := os.Open(file)
		if err != nil {
			t.Fatalf("Error opening file %s: %v", file, err)
		}
		defer f.Close()

		var positions []utils.JSONPosition
		if err := json.NewDecoder(f).Decode(&positions); err != nil {
			t.Fatalf("Error decoding JSON file %s: %v", file, err)
		}

		for i, pos := range positions {
			game, err := grid.InitGrid(pos.Sequence)
			if err != nil {
				t.Fatalf("Error initializing grid (file %d, entry %d): %v", fileIndex, i, err)
			}

			analysis, _ := game.Analyze()
			// json.NewDecoder(f).Decode doesn't support null values for int8, so the analysis in the JSON file has 0 for the columns that are not filled, but in our code we use config.NIL_SCORE for that, so we need to convert the analysis from the JSON file before comparing it with the one from our code
			savedAnalysis := convert(pos.Analysis)

			if !scoresEqual(analysis.Scores, savedAnalysis) {
				t.Errorf(
					"[File %d, Entry %d] Scores mismatch for position %q\nExpected: %v\nGot:      %v",
					fileIndex,
					i,
					pos.Sequence,
					formatScores(savedAnalysis),
					formatScores(analysis.Scores),
				)
			} else {
				t.Logf("[File %d] Position %s analyzed correctly", fileIndex, pos.Sequence)
			}
		}
	}
}

func TestSolve(t *testing.T) {
	for fileIndex, file := range files {
		f, err := os.Open(file)
		if err != nil {
			t.Fatalf("Error opening file %s: %v", file, err)
		}
		defer f.Close()

		var positions []utils.JSONPosition
		if err := json.NewDecoder(f).Decode(&positions); err != nil {
			t.Fatalf("Error decoding JSON file %s: %v", file, err)
		}

		for i, pos := range positions {
			game, err := grid.InitGrid(pos.Sequence)
			if err != nil {
				t.Fatalf("Error initializing grid (file %d, entry %d): %v", fileIndex, i, err)
			}

			solved, _ := game.Solve()
			if solved.Score == config.NIL_SCORE {
				t.Errorf("[File %d, Entry %d] No solution found for position %q", fileIndex, i, pos.Sequence)
			} else if solved.Score != pos.Score {
				t.Errorf("[File %d, Entry %d] Score mismatch for position %q\nExpected: %d\nGot:      %d", fileIndex, i, pos.Sequence, pos.Score, solved.Score)
			} else {
				t.Logf("[File %d] Position %s solved correctly", fileIndex, pos.Sequence)
			}
		}
	}
}

func scoresEqual(a, b [7]int8) bool {
	for i := range 7 {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func formatScores(s [7]int8) string {
	out := make([]string, 7)
	for i, v := range s {
		if v == config.NIL_SCORE {
			out[i] = "N/A"
		} else {
			out[i] = fmt.Sprintf("%d", v)
		}
	}
	return "[" + strings.Join(out, ", ") + "]"
}

// I don't know how to name the function but it converts the analysis from the JSON file, with the possible null values, and convert it to config.NIL_SCORE
func convert(analysis [7]*int8) [7]int8 {
	var converted [7]int8
	for i, v := range analysis {
		if v == nil {
			converted[i] = config.NIL_SCORE
		} else {
			converted[i] = *v
		}
	}
	return converted
}
