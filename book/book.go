package book

import (
	"fmt"
	"log"
	"time"

	"github.com/Pietot/Gonnect-4/database"
	"github.com/Pietot/Gonnect-4/grid"
	"go.etcd.io/bbolt"
)

const NODE_THRESHOLD = 20_000_000

func CreateBook(maxDepth int) {
	db, err := bbolt.Open(database.DBName, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("Critical error: cannot open database."+
			"Check if another instance is running: %v", err)
	}
	defer db.Close()

	db.Update(func(tx *bbolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(database.BucketResults))
		tx.CreateBucketIfNotExists([]byte(database.BucketQueue))
		tx.CreateBucketIfNotExists([]byte(database.BucketPending))
		return nil
	})

	var isEmpty bool
	db.View(func(tx *bbolt.Tx) error {
		isEmpty = tx.Bucket([]byte(database.BucketQueue)).Stats().KeyN == 0
		return nil
	})

	if isEmpty {
		startGrid, _ := grid.InitGrid("")
		db.Update(func(tx *bbolt.Tx) error {
			database.AddToQueue(tx, grid.GetCanonicalKey(startGrid), 0)
			return nil
		})
	}

	for {
		key, depth, found := database.PopFromQueue(db)
		if !found || depth > maxDepth {
			fmt.Println("End of calculation.")
			break
		}

		g := grid.FromKey(key)

		var alreadyDone bool
		db.View(func(tx *bbolt.Tx) error {
			alreadyDone = database.IsAnalyzed(tx, key)
			return nil
		})

		if alreadyDone {
			continue
		}

		analysis, stats := g.Analyze()
		db.Update(func(tx *bbolt.Tx) error {
			if stats.NodeCount >= NODE_THRESHOLD {
				database.SaveResult(tx, key, analysis.Scores)
				fmt.Printf("[D:%d] %d saved (%d nodes)\n", depth, key, stats.NodeCount)
				for col := range 7 {
					if g.CanPlay(col) && !g.IsWinningMove(col) {
						child := *g
						child.PlayColumn(col)
						cKey := grid.GetCanonicalKey(&child)

						if !database.IsAnalyzed(tx, cKey) && !database.IsInQueue(tx, cKey) {
							database.AddToQueue(tx, cKey, depth+1)
						}
					} else {
						fmt.Printf("[D:%d] %d winning move in col %d, skipping childs from this node.\n", depth, key, col)
					}
				}
			} else {
				fmt.Printf("[D:%d] %d skipped (%d nodes). No children from this node will be added to the queue nor analyzed.\n", depth, key, stats.NodeCount)
			}
			return nil
		})
	}
}
