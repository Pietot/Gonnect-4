package test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Pietot/Gonnect-4/grid"
)

func TestNegamax(t *testing.T) {
	currentPath, err := os.Getwd()
	check(err)
	files, err := os.ReadDir(currentPath + "/data")
	check(err)
	fmt.Printf("Current path: %s\n", currentPath)
	for _, file := range files {
		lines := readFile(currentPath + "/data/" + file.Name())
		for _, line := range lines {
			// split the line every space
			parts := strings.Fields(line)
			grid := createGrid(parts[0])
			expected := int([]rune(parts[1])[0] - '0')
			evaluation, _ := grid.Negamax(1)
			if *evaluation.Score > 0 && expected < 0 {
				t.Errorf("Unexpected positive evaluation for file: %s", file.Name())
			}
			if *evaluation.Score < 0 && expected > 0 {
				t.Errorf("Unexpected negative evaluation for file: %s", file.Name())
			}
			fmt.Printf("File: %s, Sequence: %s, Expected: %d, Evaluation: %f\n", file.Name(), parts[0], expected, *evaluation.Score)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(filename string) []string {
	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		check(err)
	}

	return lines
}

func createGrid(sequence string) *grid.Grid {
	player := 1
	grid := grid.InitGrid()
	for i := range sequence {
		col := int(sequence[i] - '0')
		grid.DropPiece(col-1, player)
		if player == 1 {
			player = 2
		} else {
			player = 1
		}
	}
	return grid
}
