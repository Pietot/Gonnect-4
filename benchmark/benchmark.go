package benchmark

import (
	"fmt"
	"os"
	"strings"

	"github.com/Pietot/Gonnect-4/grid"
)

var files = []string{
	"data/test_easy_end.txt",
	"data/test_easy_middle.txt",
	"data/test_easy_begin.txt",
	"data/test_medium_end.txt",
	"data/test_medium_middle.txt",
	"data/test_hard_begin.txt",
}

func Benchmark() {
	// warm up
	gameTest, _ := grid.InitGrid("533422")
	gameTest.Solve()

	for _, file := range files {
		lines, err := readPositionsFromFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
		}
		totalTimes := int64(0)
		nodeCounts := uint64(0)
		meanTimesPerNode := float64(0)
		nodesPerSecond := uint64(0)
		for _, line := range lines {
			position := strings.Split(line, " ")[0]
			game, err := grid.InitGrid(position)
			if err != nil {
				fmt.Println("Error initializing grid:", err)
			}
			_, stat := game.Solve()
			totalTimes += stat.TotalTimeNanoseconds
			nodeCounts += stat.NodeCount
			meanTimesPerNode += stat.MeanTimePerNode
			nodesPerSecond += stat.NodesPerSecond
		}

		fmt.Println("File:", file)
		fmt.Println("Mean total time (ns):", totalTimes/int64(len(lines)))
		fmt.Println("Mean node count:", nodeCounts/uint64(len(lines)))
		fmt.Println("Mean time per node (ns):", meanTimesPerNode/float64(len(lines)))
		fmt.Println("Mean nodes per second:", nodesPerSecond/uint64(len(lines)))
		fmt.Println()
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
