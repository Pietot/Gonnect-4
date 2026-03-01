package database

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"

	"go.etcd.io/bbolt"
)

var (
	BucketResults = "Results"
	BucketQueue   = "Queue"
	BucketPending = "Pending"
	DB            *bbolt.DB
)

const KEY_EMPTY_POSITION uint64 = 4432676798593

func GetDatabase(dbName string) *bbolt.DB {
	db, err := bbolt.Open(dbName, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bbolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(BucketResults))
		q, _ := tx.CreateBucketIfNotExists([]byte(BucketQueue))
		tx.CreateBucketIfNotExists([]byte(BucketPending))
		// Add initial position to queue if empty
		if q.Stats().KeyN == 0 {
			AddToQueue(tx, KEY_EMPTY_POSITION, 0)
		}
		return nil
	})

	return db
}

func Uint64ToBytes(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func AddToQueue(tx *bbolt.Tx, key uint64, depth int) {
	q := tx.Bucket([]byte(BucketQueue))
	p := tx.Bucket([]byte(BucketPending))

	ck := make([]byte, 9)
	ck[0] = byte(depth)
	binary.BigEndian.PutUint64(ck[1:], key)

	q.Put(ck, []byte{byte(depth)})
	p.Put(Uint64ToBytes(key), []byte{byte(depth)})
}

func PopFromQueue(db *bbolt.DB) (key uint64, depth int, found bool) {
	db.Update(func(tx *bbolt.Tx) error {
		q := tx.Bucket([]byte(BucketQueue))
		p := tx.Bucket([]byte(BucketPending))

		c := q.Cursor()
		k, v := c.First()
		if k == nil {
			return nil
		}

		key = binary.BigEndian.Uint64(k[1:])
		depth = int(v[0])
		found = true

		q.Delete(k)
		p.Delete(Uint64ToBytes(key))
		return nil
	})
	return
}

func IsAnalyzed(tx *bbolt.Tx, key uint64) bool {
	return tx.Bucket([]byte(BucketResults)).Get(Uint64ToBytes(key)) != nil
}

func IsInQueue(tx *bbolt.Tx, key uint64) bool {
	return tx.Bucket([]byte(BucketPending)).Get(Uint64ToBytes(key)) != nil
}

func SaveResult(tx *bbolt.Tx, key uint64, scores [7]*int8) {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(scores)
	tx.Bucket([]byte(BucketResults)).Put(Uint64ToBytes(key), buf.Bytes())
}
