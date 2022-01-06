package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

type Bolt struct {
	db *bolt.DB
}

func MustNewBolt(filepath string) *Bolt {
	b, err := NewBolt(filepath)
	if err != nil {
		panic(err)
	}
	return b
}

func NewBolt(filepath string) (*Bolt, error) {
	db, err := bolt.Open(filepath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open bolt db: %v", err)
	}
	return &Bolt{
		db: db,
	}, nil
}

func (b *Bolt) Put(ctx context.Context, bucket, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			newBucket, err := tx.CreateBucket([]byte(bucket))
			if err != nil {
				return fmt.Errorf("failed to create bucket: %v", err)
			}
			log.Printf("[INFO] bolt: created bucket: '%s'", bucket)
			b = newBucket
		}
		if err := b.Put([]byte(key), data); err != nil {
			return fmt.Errorf("failed to put value: %v", err)
		}
		return nil
	})
}

func (b *Bolt) Get(ctx context.Context, bucket, key string, dest interface{}) error {
	return b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		data := b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("key not found")
		}
		if err := json.Unmarshal(data, dest); err != nil {
			return fmt.Errorf("failed to unmarshal value: %v", err)
		}
		return nil
	})
}

func (b *Bolt) List(ctx context.Context, bucket string, dest []interface{}, limit int, from *string) (last *string, err error) {
	if cap(dest) < limit {
		err = fmt.Errorf("destination slice is too small")
		return
	}
	err = b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		c := b.Cursor()
		var k, v []byte
		if from == nil {
			k, v = c.First()
		} else {
			k, v = c.Seek([]byte(*from))
		}

		for i := 0; k != nil && i < limit; k, v = c.Next() {
			kString := string(k)
			last = &kString
			if err := json.Unmarshal(v, dest[i]); err != nil {
				return fmt.Errorf("failed to unmarshal value: %v", err)
			}
			i++
		}
		return nil
	})
	return
}
