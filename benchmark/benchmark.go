package benchmark

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Pietot/Gonnect-4/book"
	"github.com/Pietot/Gonnect-4/grid"
	"github.com/Pietot/Gonnect-4/utils"
	"github.com/schollz/progressbar/v3"
)

var files = []string{
	"data/test_easy_end.txt",
	"data/test_easy_middle.txt",
	"data/test_easy_begin.txt",
	"data/test_medium_end.txt",
	"data/test_medium_middle.txt",
	"data/test_hard_begin.txt",
}

func BenchmarkAnalyze() {
	// warm up
	gameTest, _ := grid.InitGrid("533422")
	gameTest.Analyze()

	for _, file := range files {
		totalTimes := int64(0)
		nodeCounts := uint64(0)
		meanTimesPerNode := float64(0)
		nodesPerSecond := uint64(0)
		lines, err := utils.ReadPositionsFromFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
		}
		bar := progressbar.Default(int64(len(lines)))
		for _, line := range lines {
			position := strings.Split(line, " ")[0]
			game, err := grid.InitGrid(position)
			if err != nil {
				fmt.Println("Error initializing grid:", err)
			}
			_, stat := game.Analyze()
			totalTimes += stat.TotalTimeNanoseconds
			nodeCounts += stat.NodeCount
			meanTimesPerNode += stat.MeanTimePerNode
			nodesPerSecond += stat.NodesPerSecond
			bar.Add(1)
		}
		fmt.Println("File:                   ", file)
		fmt.Println("Mean total time (ns):   ", totalTimes/int64(len(lines)))
		fmt.Println("Mean node count:        ", nodeCounts/uint64(len(lines)))
		fmt.Println("Mean time per node (ns):", meanTimesPerNode/float64(len(lines)))
		fmt.Println("Mean nodes per second:  ", nodesPerSecond/uint64(len(lines)))
		fmt.Println()
	}
	fmt.Println()
}

func BenchmarkBookCreation() {
	dbName := "benchmark/badger_benchmark"
	os.Remove(dbName)

	start := time.Now()
	book.CreateBook(8, dbName)
	elapsed := time.Since(start)
	fmt.Printf("Book creation completed in %s\n", elapsed)
}
