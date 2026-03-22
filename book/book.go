package book

import (
	"github.com/Pietot/Gonnect-4/database"
	"github.com/Pietot/Gonnect-4/grid"
	"github.com/Pietot/Gonnect-4/progressbar"
	"github.com/dgraph-io/badger/v4"
)

const (
	NODE_THRESHOLD = 20_000_000
	POPPED_KEY     = 1
)

func CreateBook(maxDepth int, dbName string) {
	pb := progressbar.NewProgressBar()
	db := database.GetDatabase(dbName)
	defer db.Close()

	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = database.PrefixPending
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			pb.TotalInQueue++
		}
		return nil
	})

	pb.Render()

	for {
		key, depth, found := database.PopFromQueue(db)
		pb.RemoveFromQueue()
		if !found {
			break
		}

		if depth > maxDepth {
			pb.Render()
			continue
		}

		if depth != pb.CurrentDepth {
			totalForDepth := database.CountKeysForDepth(db, depth) + POPPED_KEY
			pb.ResetDepth(depth, totalForDepth)
			pb.Render()
		}

		g := grid.FromKey(key)
		analysis, stats := g.Analyze()

		if stats.NodeCount >= NODE_THRESHOLD || isAlreadyInBook(stats.NodeCount) {
			err := db.Update(func(txn *badger.Txn) error {
				database.SaveResult(txn, key, analysis.Scores)
				pb.AddSaved()

				for col := range 7 {
					if g.CanPlay(col) && !g.IsWinningMove(col) {
						child := *g
						child.PlayColumn(col)
						cKey := child.GetCanonicalKey()

						if !database.IsAnalyzed(txn, cKey) && !database.IsInQueue(txn, cKey) {
							database.AddToQueue(txn, cKey, depth+1)
							pb.AddToQueue()
						}
					} else {
						pb.AddSkipped()
					}
				}
				return nil
			})
			if err != nil {
				panic(err)
			}
		} else {
			pb.AddSkipped()
		}

		pb.AddAnalyzed()
		pb.Render()
	}
}

func isAlreadyInBook(nodeCount uint64) bool {
	return nodeCount == 0
}
