package boltdb

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/bahusvel/TunnelBeast/iptables"
	"github.com/boltdb/bolt"
)

var (
	db *bolt.DB

	BUCKETNAME        = "userRoutes"
	DBFILE            = "bolt.db"
	ErrBucketNotFound = errors.New("Error Bucket not found")
	ErrDoesNotExist   = errors.New("Error Input")
	ErrExists         = errors.New("Error Record Exists")
)

func Init() error {
	db, _ = bolt.Open(DBFILE, 0600, nil)
	return nil
}

func AddRecord(key string, value iptables.NATEntry) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BUCKETNAME))
		if err != nil {
			return ErrBucketNotFound
		}
		v := bucket.Get([]byte(key))
		if v != nil {
			return ErrExists
		}
		val, _ := json.Marshal(value)
		return bucket.Put([]byte(key), val)
	})
}

func DeleteRecord(key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}
		return bucket.Delete([]byte(key))
	})
}

func ListRecords(username string) (keys []string, values []iptables.NATEntry, err error) {
	return keys, values, db.View(func(tx *bolt.Tx) error {
		//bucket, err := tx.CreateBucketIfNotExists([]byte(BUCKETNAME))
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}

		return bucket.ForEach(func(k, v []byte) error {
			log.Println(string(k), string(v))
			if strings.Contains(string(k), username+"/") {
				key := string(k)
				keys = append(keys, key)
				var value iptables.NATEntry
				err := json.Unmarshal(v, &value)
				if err != nil {
					log.Println(err)
				}
				values = append(values, value)
			}
			return nil // Continue ForEach
		})
	})
}
