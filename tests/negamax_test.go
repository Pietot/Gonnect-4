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
			grid := grid.InitGrid(parts[0])
			expected := int8([]rune(parts[1])[0] - '0')
			score := grid.GetScore()
			if score != expected {
				t.Errorf("Unexpected score for file: %s and sequence: %s, Expected: %d, Got: %d", file.Name(), parts[0], expected, score)
			}
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
