package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
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
			return ErrNotFound
		}
		data := b.Get([]byte(key))
		if data == nil {
			return ErrNotFound
		}
		if err := json.Unmarshal(data, dest); err != nil {
			return fmt.Errorf("failed to unmarshal value: %v", err)
		}
		return nil
	})
}

func (b *Bolt) List(ctx context.Context, bucket string, dest interface{}, limit int, from *string) (last *string, err error) {
	// make sure that the dest is a pointer
	destValuePtr := reflect.ValueOf(dest)
	if destValuePtr.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("dest must be a pointer")
	}

	// make sure that the dest is a pointer to a slice or an array, so that we can append to it
	destValue := destValuePtr.Elem()
	if destValue.Kind() != reflect.Slice && destValue.Kind() != reflect.Array {
		return nil, fmt.Errorf("dest must be a pointer to an array or a slice")
	}

	// get the type of the slice or array element
	elemValueType := destValue.Type().Elem()
	// get the value type in a case when dest is an array or slice of pointers
	if elemValueType.Kind() == reflect.Ptr {
		elemValueType = elemValueType.Elem()
	}
	if elemValueType.Kind() == reflect.Ptr {
		return nil, fmt.Errorf("dest value can obly have one level of indirection")
	}

	err = b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
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

			// initialize an empty array element, and get a pointer to it
			destElemValue := reflect.Zero(elemValueType)
			destElemValuePtr := reflect.New(elemValueType)
			destElemValuePtr.Elem().Set(destElemValue)

			if err := json.Unmarshal(v, destElemValuePtr.Interface()); err != nil {
				return fmt.Errorf("failed to unmarshal value: %v", err)
			}

			// append the element to the destination array of slice, depending on it's kind
			if destValue.Type().Elem().Kind() == reflect.Ptr {
				destValue.Set(reflect.Append(destValue, destElemValuePtr))
			} else {
				destValue.Set(reflect.Append(destValue, destElemValuePtr.Elem()))
			}

			i++
		}
		return nil
	})
	return
}
