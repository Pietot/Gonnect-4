package benchmark

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Pietot/Gonnect-4/book"
	"github.com/Pietot/Gonnect-4/database"
	"github.com/Pietot/Gonnect-4/grid"
	"go.etcd.io/bbolt"
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

	totalTimes := int64(0)
	nodeCounts := uint64(0)
	meanTimesPerNode := float64(0)
	nodesPerSecond := uint64(0)
	totalAnalyses := 0
	for _, file := range files {
		lines, err := readPositionsFromFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
		}
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
			totalAnalyses++
		}
	}
	fmt.Println("Mean total time (ns):", totalTimes/int64(totalAnalyses))
	fmt.Println("Mean node count:", nodeCounts/uint64(totalAnalyses))
	fmt.Println("Mean time per node (ns):", meanTimesPerNode/float64(totalAnalyses))
	fmt.Println("Mean nodes per second:", nodesPerSecond/uint64(totalAnalyses))
	fmt.Println()
}

func BenchmarkBookCreation() {
	// warm up
	gameTest, _ := grid.InitGrid("533422")
	gameTest.Analyze()

	os.Remove("benchmark/book_benchmark_d8.db")

	bookD8, err := bbolt.Open("benchmark/book_benchmark_d8.db", 0600, nil)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer bookD8.Close()

	database.DB = bookD8

	start := time.Now()
	book.CreateBook(8)
	elapsed := time.Since(start)
	fmt.Printf("Book creation completed in %s\n", elapsed)
}

func readPositionsFromFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	return lines, nil
}
