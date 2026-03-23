package benchmark

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Pietot/Gonnect-4/book"
	"github.com/Pietot/Gonnect-4/config"
	"github.com/Pietot/Gonnect-4/grid"
	"github.com/Pietot/Gonnect-4/stats"
	"github.com/Pietot/Gonnect-4/utils"
	"github.com/schollz/progressbar/v3"
)

var files = []string{
	"data/test_easy_end.json",
	"data/test_easy_middle.json",
	"data/test_easy_begin.json",
	"data/test_medium_end.json",
	"data/test_medium_middle.json",
	"data/test_hard_begin.json",
}

type benchmarkFunc func(*grid.Grid) (any, stats.Stats)

func runBenchmark(name string, benchFunc benchmarkFunc) {
	// warm up
	gameTest, _ := grid.InitGrid("533422")
	benchFunc(gameTest)

	start := time.Now()

	for _, file := range files {
		var totalTimes int64
		var nodeCounts uint64
		var meanTimesPerNode float64
		var nodesPerSecond uint64

		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("Error opening file %s: %v", file, err)
		}
		defer f.Close()

		var positions []utils.JSONPosition
		if err := json.NewDecoder(f).Decode(&positions); err != nil {
			log.Fatalf("Error decoding JSON file %s: %v", file, err)
		}

		bar := progressbar.Default(int64(len(positions)))
		for _, pos := range positions {
			game, err := grid.InitGrid(pos.Sequence)
			if err != nil {
				log.Printf("Error initializing grid for position %q: %v", pos.Sequence, err)
				continue
			}
			_, stat := benchFunc(game)
			totalTimes += stat.TotalTimeNanoseconds
			nodeCounts += stat.NodeCount
			meanTimesPerNode += stat.MeanTimePerNode
			nodesPerSecond += stat.NodesPerSecond
			bar.Add(1)
		}
		fmt.Printf("%s - File:                %s\n", name, file)
		fmt.Printf("Mean total time (ns):     %d\n", totalTimes/int64(len(positions)))
		fmt.Printf("Mean node count:          %d\n", nodeCounts/uint64(len(positions)))
		fmt.Printf("Mean time per node (ns):  %.2f\n", meanTimesPerNode/float64(len(positions)))
		fmt.Printf("Mean nodes per second:    %d\n\n", nodesPerSecond/uint64(len(positions)))
	}
	fmt.Printf("%s benchmark completed in %s\n", name, time.Since(start))
	fmt.Println()
}

func BenchmarkSolve() {
	runBenchmark("Solve", func(g *grid.Grid) (any, stats.Stats) {
		return g.Solve()
	})
}

func BenchmarkAnalyze() {
	runBenchmark("Analyze", func(g *grid.Grid) (any, stats.Stats) {
		return g.Analyze()
	})
}

func BenchmarkBookCreation() {
	os.Remove(config.BENCHMARK_DB_PATH)

	start := time.Now()
	book.CreateBook(8, config.BENCHMARK_DB_PATH)
	elapsed := time.Since(start)
	fmt.Printf("Book creation completed in %s\n", elapsed)
}
