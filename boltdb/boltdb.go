package boltdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

type RecordValue struct {
	RecordName    string
	DestinationIP string
	ExternalPort  string
	InternalPort  string
}

var (
	db *bolt.DB

	BUCKETNAME          = "userRoutes"
	ErrBucketNotCreated = errors.New("Error Bucket not created")
	ErrBucketNotFound   = errors.New("Error Bucket not found")
	ErrExists           = errors.New("ERROR RECORD EXISTS")
	ErrNotExist         = errors.New("ERROR RECORD NOT EXIST")
)

func Init(Path string) error {
	var err error
	db, err = bolt.Open(Path, 0600, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(BUCKETNAME))
		if err != nil {
			log.Println(err)
			return ErrBucketNotCreated
		}
		return nil
	})
	return nil
}

func AddRecord(key string, value RecordValue) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}
		v := bucket.Get([]byte(key))
		if v != nil {
			return ErrExists
		}
		val, err := json.Marshal(value)
		if err != nil {
			log.Println(err)
			return err
		}
		return bucket.Put([]byte(key), val)
	})
}

func DeleteRecord(key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}
		v := bucket.Get([]byte(key))
		if v == nil {
			return ErrNotExist
		}
		return bucket.Delete([]byte(key))
	})
}

func ListRecords(username string) (values []RecordValue, err error) {
	return values, db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}

		cursor := bucket.Cursor()
		prefix := []byte(username)

		for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
			var value RecordValue
			err := json.Unmarshal(v, &value)
			if err != nil {
				log.Println(err)
				return err
			}
			values = append(values, value)
		}
		return nil
	})
}
