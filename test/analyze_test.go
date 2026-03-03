package test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Pietot/Gonnect-4/grid"
)

var files = []string{
	"../data/test_easy_end.json",
	"../data/test_easy_middle.json",
	"../data/test_easy_begin.json",
	"../data/test_medium_end.json",
	"../data/test_medium_middle.json",
	"../data/test_hard_begin.json",
}

type JSONPosition struct {
	Sequence string   `json:"sequence"`
	Score    int8     `json:"score"`
	Analysis [7]*int8 `json:"analysis"`
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

		var positions []JSONPosition
		if err := json.NewDecoder(f).Decode(&positions); err != nil {
			t.Fatalf("Error decoding JSON file %s: %v", file, err)
		}

		for i, pos := range positions {
			game, err := grid.InitGrid(pos.Sequence)
			if err != nil {
				t.Fatalf("Error initializing grid (file %d, entry %d): %v", fileIndex, i, err)
			}

			anal, _ := game.Analyze()

			if !scoresEqual(anal.Scores, pos.Analysis) {
				t.Errorf(
					"[File %d, Entry %d] Scores mismatch for position %q\nExpected: %v\nGot:      %v",
					fileIndex,
					i,
					pos.Sequence,
					formatScores(pos.Analysis),
					formatScores(anal.Scores),
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

		var positions []JSONPosition
		if err := json.NewDecoder(f).Decode(&positions); err != nil {
			t.Fatalf("Error decoding JSON file %s: %v", file, err)
		}

		for i, pos := range positions {
			game, err := grid.InitGrid(pos.Sequence)
			if err != nil {
				t.Fatalf("Error initializing grid (file %d, entry %d): %v", fileIndex, i, err)
			}

			solved, _ := game.Solve()
			if solved.Score == nil {
				t.Errorf("[File %d, Entry %d] No solution found for position %q", fileIndex, i, pos.Sequence)
			} else if *solved.Score != pos.Score {
				t.Errorf("[File %d, Entry %d] Score mismatch for position %q\nExpected: %d\nGot:      %d", fileIndex, i, pos.Sequence, pos.Score, *solved.Score)
			} else {
				t.Logf("[File %d] Position %s solved correctly", fileIndex, pos.Sequence)
			}
		}
	}
}

func scoresEqual(a, b [7]*int8) bool {
	for i := range 7 {
		if a[i] == nil && b[i] == nil {
			continue
		}
		if a[i] == nil || b[i] == nil {
			return false
		}
		if *a[i] != *b[i] {
			return false
		}
	}
	return true
}

func formatScores(s [7]*int8) string {
	out := make([]string, 7)
	for i, v := range s {
		if v == nil {
			out[i] = "null"
		} else {
			out[i] = fmt.Sprintf("%d", *v)
		}
	}
	return "[" + strings.Join(out, ", ") + "]"
}
