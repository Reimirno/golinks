package bolt_mapper

import (
	"fmt"

	"github.com/boltdb/bolt"
)

func (b *BoltMapper) initializeBucket(name string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		return err
	})
}

func (b *BoltMapper) get(bucketName string, key string) ([]byte, error) {
	var value []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket not found: %s", bucketName)
		}
		value = b.Get([]byte(key)) // value is nil if key not found
		return nil
	})
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (b *BoltMapper) put(bucketName string, key string, value []byte) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket not found: %s", bucketName)
		}
		return b.Put([]byte(key), []byte(value))
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltMapper) delete(bucketName string, key string) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket not found: %s", bucketName)
		}
		return b.Delete([]byte(key))
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltMapper) foreach(bucketName string, action func(key string, value []byte) error) error {
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket not found: %s", bucketName)
		}
		return b.ForEach(func(k, v []byte) error {
			return action(string(k), v)
		})
	})
	return err
}
