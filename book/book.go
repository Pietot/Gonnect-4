package book

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/Pietot/Gonnect-4/database"
	"github.com/Pietot/Gonnect-4/evaluation"
	"github.com/Pietot/Gonnect-4/grid"
	"github.com/Pietot/Gonnect-4/stats"
	"go.etcd.io/bbolt"
)

const (
	WORKER_COUNT   = 0
	JOB_QUEUE_SIZE = 1000
	NODE_THRESHOLD = 20_000_000
)

type Job struct {
	Key   uint64
	Depth int
}

type Result struct {
	Key      uint64
	Depth    int
	Analysis evaluation.Analysis
	Stats    stats.Stats
	Grid     *grid.Grid
	WorkerID int
}

func CreateBook(maxDepth int) {
	jobs := make(chan Job, JOB_QUEUE_SIZE)
	results := make(chan Result, JOB_QUEUE_SIZE)

	var wg sync.WaitGroup

	numWorkers := WORKER_COUNT
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	fmt.Printf("Starting %d workers...\n", numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, jobs, results)
		}(i)
	}

	saveDone := make(chan bool)
	go func() {
		collector(results)
		saveDone <- true
	}()

	for {
		// PopFromQueue is a write transaction (Delete), so it's blocking.
		// Possible optimization: Read in batches of 100 keys to reduce transaction overhead.
		key, depth, found := database.PopFromQueue(database.DB)
		if !found {
			// Nothing left in the queue? We wait a bit, because the Collector might
			// be adding new children.
			// If the results channel is also empty, it's really finished.
			if len(results) == 0 && len(jobs) == 0 {
				fmt.Println("Queue empty and workers inactive. End of calculation.")
				break
			}
			time.Sleep(1 * time.Second)
			continue
		}

		if depth > maxDepth {
			fmt.Printf("\033[33m[D:%d] %d reached max depth. End of calculation.\033[0m\n", depth, key)
			continue
		}

		jobs <- Job{Key: key, Depth: depth}
	}

	close(jobs)
	wg.Wait()
	close(results)
	<-saveDone
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
	for job := range jobs {
		var alreadyDone bool
		database.DB.View(func(tx *bbolt.Tx) error {
			alreadyDone = database.IsAnalyzed(tx, job.Key)
			return nil
		})

		if alreadyDone {
			continue
		}

		g := grid.FromKey(job.Key)
		analysis, stats := g.Analyze()

		results <- Result{
			Key:      job.Key,
			Depth:    job.Depth,
			Analysis: analysis,
			Stats:    stats,
			Grid:     g,
			WorkerID: id,
		}

	}
}

func collector(results <-chan Result) {
	// To optimize BoltDB, we can batch writes,
	// but for simplicity, we keep one transaction per result here.
	// Ideally: accumulate X results or wait T time before committing.

	for res := range results {
		database.DB.Update(func(tx *bbolt.Tx) error {

			if res.Stats.NodeCount >= NODE_THRESHOLD {
				database.SaveResult(tx, res.Key, res.Analysis.Scores)
				fmt.Printf("\033[32m[D:%d-W:%d] %d saved (%d nodes)\033[0m\n", res.Depth, res.WorkerID, res.Key, res.Stats.NodeCount)

				for col := range 7 {
					if res.Grid.CanPlay(col) && !res.Grid.IsWinningMove(col) {
						child := *res.Grid
						child.PlayColumn(col)
						cKey := grid.GetCanonicalKey(&child)
						if !database.IsAnalyzed(tx, cKey) && !database.IsInQueue(tx, cKey) {
							database.AddToQueue(tx, cKey, res.Depth+1)
						}
					} else {
						fmt.Printf("\033[33m[D:%d-W:%d] %d winning move in col %d, skipping childs from this node.\033[0m\n", res.Depth, res.WorkerID, res.Key, col)
					}
				}

			} else {
				fmt.Printf("\033[31m[D:%d-W:%d] %d skipped (%d nodes).No children from this node will be added to the queue nor analyzed.\033[0m\n", res.Depth, res.WorkerID, res.Key, res.Stats.NodeCount)
			}
			return nil
		})
	}
}
