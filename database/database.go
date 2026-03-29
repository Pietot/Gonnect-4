package database

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"

	"github.com/dgraph-io/badger/v4"
)

var (
	PrefixResults = []byte("R:")
	PrefixQueue   = []byte("Q:")
	PrefixPending = []byte("P:")
)

const KEY_EMPTY_POSITION uint64 = 4432676798593

func makeKey(prefix []byte, key []byte) []byte {
	k := make([]byte, len(prefix)+len(key))
	copy(k, prefix)
	copy(k[len(prefix):], key)
	return k
}

func GetDatabase(dbName string) *badger.DB {
	opts := badger.DefaultOptions(dbName)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{Prefix: PrefixQueue})
		defer it.Close()

		it.Rewind()
		if !it.Valid() {
			return AddToQueue(txn, KEY_EMPTY_POSITION, 0)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func uint64ToBytes(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func AddToQueue(txn *badger.Txn, key uint64, depth int) error {
	// Key for the queue : [1 byte for depth] + [8 bytes for the key]
	qk := make([]byte, 9)
	qk[0] = byte(depth)
	binary.BigEndian.PutUint64(qk[1:], key)

	err := txn.Set(makeKey(PrefixQueue, qk), []byte{byte(depth)})
	if err != nil {
		return err
	}

	return txn.Set(makeKey(PrefixPending, uint64ToBytes(key)), []byte{byte(depth)})
}

func PopFromQueue(db *badger.DB) (key uint64, depth int, found bool) {
	db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = PrefixQueue
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		it.Rewind()
		if !it.Valid() {
			return nil
		}

		item := it.Item()
		k := item.KeyCopy(nil)

		payload := k[len(PrefixQueue):]
		depth = int(payload[0])
		key = binary.BigEndian.Uint64(payload[1:])
		found = true

		txn.Delete(k)
		txn.Delete(makeKey(PrefixPending, uint64ToBytes(key)))

		return nil
	})
	return
}

func IsAnalyzed(txn *badger.Txn, key uint64) bool {
	_, err := txn.Get(makeKey(PrefixResults, uint64ToBytes(key)))
	return err == nil
}

func IsInQueue(txn *badger.Txn, key uint64) bool {
	_, err := txn.Get(makeKey(PrefixPending, uint64ToBytes(key)))
	return err == nil
}

func SaveResult(txn *badger.Txn, key uint64, scores [7]int8) error {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(scores)
	if err != nil {
		return err
	}
	return txn.Set(makeKey(PrefixResults, uint64ToBytes(key)), buf.Bytes())
}

func CountKeysForDepth(db *badger.DB, depth int) int {
	count := 0
	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = PrefixQueue
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := append(PrefixQueue, byte(depth))
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			count++
		}
		return nil
	})
	return count
}
