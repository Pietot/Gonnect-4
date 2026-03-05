package database

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v4"
)

var (
	// Prefixes pour simuler des buckets
	PrefixResults = []byte{0x01}
	PrefixQueue   = []byte{0x02}
	PrefixPending = []byte{0x03}
	DB            *badger.DB
)

const KEY_EMPTY_POSITION uint64 = 4432676798593

func GetDatabase(path string) {
	opts := badger.DefaultOptions(path).
		WithLoggingLevel(badger.WARNING).
		WithIndexCacheSize(512 << 20). // 512MB de cache d'index pour accélérer les IsAnalyzed
		WithSyncWrites(false)          // On privilégie la vitesse pour le pre-compute

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	DB = db

	// Initialisation : ajout de la première position si la DB est neuve
	err = DB.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		it.Seek(PrefixQueue)
		if !it.ValidForPrefix(PrefixQueue) {
			return AddToQueue(txn, KEY_EMPTY_POSITION, 0)
		}
		return nil
	})
	if err != nil {
		log.Fatal("Erreur init DB:", err)
	}
}

func GetScores(key uint64, mirroredKey uint64) (scores [7]*int8, found bool) {
	err := DB.View(func(txn *badger.Txn) error {
		// 1. Chercher la clé originale
		k := make([]byte, 9)
		k[0] = PrefixResults[0]
		binary.BigEndian.PutUint64(k[1:], key)

		item, err := txn.Get(k)
		if err == nil {
			err = item.Value(func(val []byte) error {
				scores = bytesToScores(val)
				found = true
				return nil
			})
			return err
		}

		// 2. Si non trouvé, chercher la clé miroir
		binary.BigEndian.PutUint64(k[1:], mirroredKey)
		item, err = txn.Get(k)
		if err == nil {
			err = item.Value(func(val []byte) error {
				rawScores := bytesToScores(val)
				// Inverser les scores pour la position miroir (col 0 devient col 6, etc.)
				for i := 0; i < 4; i++ {
					scores[i], scores[6-i] = rawScores[6-i], rawScores[i]
				}
				found = true
				return nil
			})
			return err
		}

		return nil // Pas trouvé
	})

	if err != nil && err != badger.ErrKeyNotFound {
		fmt.Printf("Erreur lecture scores: %v\n", err)
	}
	return
}

// Helper interne pour décoder les 7 octets vers [7]*int8
func bytesToScores(b []byte) (scores [7]*int8) {
	if len(b) < 7 {
		return
	}
	for i := range 7 {
		if b[i] == 127 { // 127 est notre marqueur pour 'nil' défini dans le collector
			scores[i] = nil
		} else {
			val := int8(b[i])
			scores[i] = &val
		}
	}
	return
}

// Uint64ToBytes sans allocation inutile
func Uint64ToBytes(v uint64, buf []byte) {
	binary.BigEndian.PutUint64(buf, v)
}

// AddToQueue prépare les clés pour Queue (q + depth + key) et Pending (p + key)
func AddToQueue(txn *badger.Txn, key uint64, depth int) error {
	qk := make([]byte, 10) // 1 (prefix) + 1 (depth) + 8 (key)
	qk[0] = PrefixQueue[0]
	qk[1] = byte(depth)
	binary.BigEndian.PutUint64(qk[2:], key)

	pk := make([]byte, 9) // 1 (prefix) + 8 (key)
	pk[0] = PrefixPending[0]
	binary.BigEndian.PutUint64(pk[1:], key)

	if err := txn.Set(qk, []byte{byte(depth)}); err != nil {
		return err
	}
	return txn.Set(pk, []byte{byte(depth)})
}

// PopBatch extrait plusieurs éléments d'un coup pour amortir le coût de la transaction
func PopBatch(batchSize int) ([]uint64, []int, error) {
	var keys []uint64
	var depths []int

	err := DB.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(PrefixQueue); it.ValidForPrefix(PrefixQueue) && len(keys) < batchSize; it.Next() {
			item := it.Item()
			k := item.Key()

			d := int(k[1])
			v := binary.BigEndian.Uint64(k[2:])

			keys = append(keys, v)
			depths = append(depths, d)

			// Supprimer uniquement de la Queue; PrefixPending reste jusqu'à ce que le collecteur finisse
			if err := txn.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})
	return keys, depths, err
}

func IsAnalyzed(txn *badger.Txn, key uint64) bool {
	var k [9]byte
	k[0] = PrefixResults[0]
	binary.BigEndian.PutUint64(k[1:], key)
	_, err := txn.Get(k[:])
	return err == nil
}

func CountQueue() int64 {
	var count int64
	DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(PrefixQueue); it.ValidForPrefix(PrefixQueue); it.Next() {
			count++
		}
		return nil
	})
	return count
}

func RemovePending(txn *badger.Txn, key uint64) error {
	var pk [9]byte
	pk[0] = PrefixPending[0]
	binary.BigEndian.PutUint64(pk[1:], key)
	return txn.Delete(pk[:])
}

func IsInQueue(txn *badger.Txn, key uint64) bool {
	var k [9]byte
	k[0] = PrefixPending[0]
	binary.BigEndian.PutUint64(k[1:], key)
	_, err := txn.Get(k[:])
	return err == nil
}
