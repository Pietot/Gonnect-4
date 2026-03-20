package utils

import (
	"fmt"
	"os"
	"strings"
)

type JSONPosition struct {
	Sequence string  `json:"sequence"`
	Score    int8    `json:"score"`
	Analysis [7]int8 `json:"analysis"`
}

func FormatFloat(value float64) string {
	formatted := fmt.Sprintf("%.2f", value)

	parts := strings.Split(formatted, ".")
	parts[0] = addUnderscores(parts[0])

	return strings.Join(parts, ".")
}

func FormatUint64(value uint64) string {
	return addUnderscores(fmt.Sprintf("%d", value))
}

func GetTime(nanoseconds int64) string {
	if nanoseconds < 1_000 {
		return fmt.Sprintf("%d ns", nanoseconds)
	} else if nanoseconds < 1_000_000 {
		return fmt.Sprintf("%.2f µs", float64(nanoseconds)/1_000)
	} else if nanoseconds < 1_000_000_000 {
		return fmt.Sprintf("%.2f ms", float64(nanoseconds)/1_000_000)
	} else {
		return fmt.Sprintf("%.2f s", float64(nanoseconds)/1_000_000_000)
	}
}

func addUnderscores(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}

	var result strings.Builder
	mod := n % 3
	if mod == 0 {
		mod = 3
	}

	result.WriteString(s[:mod])
	for i := mod; i < n; i += 3 {
		result.WriteString("_")
		result.WriteString(s[i : i+3])
	}

	return result.String()
}

func GetBestScoreAndMove(scores [7]int8) (bestScore int8, bestMove uint8) {
	bestScore = -128
	bestMove = 0
	for i, score := range scores {
		if score != -128 && score > bestScore {
			bestScore = score
			bestMove = uint8(i)
		}
	}
	return bestScore, bestMove
}

func GetScores(book *map[uint64][7]int8, key uint64, mirrorKey uint64) (scores [7]int8, found bool) {
	if scores, found = (*book)[key]; found {
		return scores, true
	}
	if scores, found = (*book)[mirrorKey]; found {
		// Reverse the scores for the mirrored position
		for i := range 3 {
			scores[i], scores[6-i] = scores[6-i], scores[i]
		}
		return scores, true
	}
	return scores, false
}

func ReadPositionsFromFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	return lines, nil
}
