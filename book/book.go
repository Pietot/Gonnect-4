package book

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Pietot/Gonnect-4/database"
	"github.com/Pietot/Gonnect-4/evaluation"
	"github.com/Pietot/Gonnect-4/grid"
	"github.com/Pietot/Gonnect-4/progressbar"
	"github.com/Pietot/Gonnect-4/stats"
	"github.com/Pietot/Gonnect-4/transpositiontable"
	"github.com/dgraph-io/badger/v4"
)

const (
	JOB_QUEUE_SIZE = 5000
	NODE_THRESHOLD = 20_000_000
	BATCH_POP_SIZE = 500
)

type Job struct {
	Key   uint64
	Depth int
}

type Result struct {
	Key         uint64
	Depth       int
	Analysis    evaluation.Analysis
	Stats       stats.Stats
	Grid        *grid.Grid
	AlreadyDone bool
}

func CreateBook(maxDepth int) {
	dbName := "database/badger"
	database.GetDatabase(dbName)
	defer database.DB.Close()

	jobs := make(chan Job, JOB_QUEUE_SIZE)
	results := make(chan Result, JOB_QUEUE_SIZE)
	var activeJobs int64
	var wg sync.WaitGroup

	progress := progressbar.New()
	progress.InitQueueSize(database.CountQueue())

	numWorkers := runtime.NumCPU()
	fmt.Printf("Starting %d workers...\n\n", numWorkers)

	for i := range numWorkers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(jobs, results)
		}(i)
	}

	saveDone := make(chan bool)
	go func() {
		collector(results, &activeJobs, progress)
		saveDone <- true
	}()

	for {
		keys, depths, err := database.PopBatch(BATCH_POP_SIZE)
		if err != nil || len(keys) == 0 {
			if atomic.LoadInt64(&activeJobs) <= 0 {
				fmt.Println("\nQueue empty and workers inactive. End of calculation.")
				break
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}

		progress.RecordDequeued(int64(len(keys)))

		for i := range keys {
			if depths[i] > maxDepth {
				continue
			}
			progress.RecordDispatched(depths[i])
			atomic.AddInt64(&activeJobs, 1)
			jobs <- Job{Key: keys[i], Depth: depths[i]}
		}
	}

	close(jobs)
	wg.Wait()
	close(results)
	<-saveDone
}

func worker(jobs <-chan Job, results chan<- Result) {
	tt := transpositiontable.NewTranspositionTable()
	for job := range jobs {
		alreadyDone := false
		database.DB.View(func(txn *badger.Txn) error {
			alreadyDone = database.IsAnalyzed(txn, job.Key)
			return nil
		})

		if alreadyDone {
			results <- Result{Key: job.Key, AlreadyDone: true}
			continue
		}

		tt.Reset()
		g := grid.FromKey(job.Key)
		g.TransTable = tt
		analysis, stats := g.Analyze()

		results <- Result{
			Key:      job.Key,
			Depth:    job.Depth,
			Analysis: analysis,
			Stats:    stats,
			Grid:     g,
		}
	}
}

func collector(results <-chan Result, activeJobs *int64, progress *progressbar.Progress) {
	wb := database.DB.NewWriteBatch()
	defer func() { wb.Flush() }()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case res, ok := <-results:
			if !ok {
				return
			}
			if res.AlreadyDone {
				database.DB.Update(func(txn *badger.Txn) error {
					return database.RemovePending(txn, res.Key)
				})
				atomic.AddInt64(activeJobs, -1)
				continue
			}

			// If NodeCount is equal to 0, it means the position was in map.go. So it's already analyzed and above the threshold.
			if res.Stats.NodeCount >= NODE_THRESHOLD || res.Stats.NodeCount == 0 {
				var resKey [9]byte
				resKey[0] = database.PrefixResults[0]
				binary.BigEndian.PutUint64(resKey[1:], res.Key)

				scoreBytes := make([]byte, 7)
				for i, s := range res.Analysis.Scores {
					if s != nil {
						scoreBytes[i] = byte(*s)
					} else {
						scoreBytes[i] = 127
					}
				}
				wb.Set(resKey[:], scoreBytes)

				var childrenQueued int64
				database.DB.Update(func(txn *badger.Txn) error {
					database.RemovePending(txn, res.Key)
					for col := range 7 {
						if res.Grid.CanPlay(col) && !res.Grid.IsWinningMove(col) {
							child := *res.Grid
							child.PlayColumn(col)
							cKey := child.GetCanonicalKey()

							if !database.IsAnalyzed(txn, cKey) && !database.IsInQueue(txn, cKey) {
								if err := database.AddToQueue(txn, cKey, res.Depth+1); err == nil {
									childrenQueued++
								}
							}
						}
					}
					return nil
				})

				progress.RecordSaved(res.Depth)
				progress.RecordEnqueued(res.Depth, childrenQueued)
			} else {
				database.DB.Update(func(txn *badger.Txn) error {
					return database.RemovePending(txn, res.Key)
				})
				progress.RecordSkipped(res.Depth)
			}
			atomic.AddInt64(activeJobs, -1)

		case <-ticker.C:
			progress.Render()
			wb.Flush()
			wb = database.DB.NewWriteBatch()
		}
	}
}
