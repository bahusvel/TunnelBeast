package boltdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

type Favourite struct {
	Name          string
	DestinationIP string
	ExternalPort  string
	InternalPort  string
}

type UserInfo struct {
	Password string //hashed
}

var (
	db *bolt.DB

	BUCKETNAME          = "userRoutes"
	ErrBucketNotCreated = errors.New("Error Bucket not created")
	ErrBucketNotFound   = errors.New("Error Bucket not found")
	ErrExists           = errors.New("ERROR EXISTS")
	ErrNotExists        = errors.New("ERROR NOT EXISTS")
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

func AddFavourite(key string, value Favourite) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
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

func DeleteFavourite(key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}
		return bucket.Delete([]byte(key))
	})
}

func GetFavourite(key string) (value Favourite, err error) {
	return value, db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}
		v := bucket.Get([]byte(key))
		if v == nil {
			return ErrNotExists
		}
		err = json.Unmarshal(v, &value)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})
}

func ListFavourites(username string) (values []Favourite, err error) {
	return values, db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}

		cursor := bucket.Cursor()
		prefix := []byte(username + "/favorite/")

		for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
			var value Favourite
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

func AddUser(username string, password string) error {

	return nil
}

func UpdateUser(username string, password string) error {
	return nil
}

func DeleteUser(username string, password string) error {
	return nil
}

func ListUsers() (users []string, err error) {
	return users, db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}

		cursor := bucket.Cursor()
		prefix := []byte("users/")

		for k, _ := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = cursor.Next() {
			log.Println(string(k), ":", len(prefix))
			user := string(k)[len(prefix):]
			if err != nil {
				log.Println(err)
				return err
			}

			users = append(users, user)
		}
		return nil
	})
}
