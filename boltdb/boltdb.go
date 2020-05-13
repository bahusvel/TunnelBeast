package boltdb

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

type Favorite struct {
	FavoriteName  string
	DestinationIP string
	ExternalPort  string
	InternalPort  string
}

type User struct {
	Password  string //hashed
	Favorites map[string]Favorite
}

var (
	db *bolt.DB

	BUCKETNAME          = "user"
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

func AddFavorite(username string, favorite Favorite) error {
	var user User
	user, err := getUser(username)
	if err != nil && err != ErrNotExists {
		//ErrNotExists happens if LDAP auth used
		log.Println(err)
		return err
	}

	if _, ok := user.Favorites[favorite.FavoriteName]; ok {
		return ErrExists
	}

	if user.Favorites == nil {
		user.Favorites = make(map[string]Favorite)
	}

	user.Favorites[favorite.FavoriteName] = favorite

	return saveUser(username, user)
}

func DeleteFavorite(username string, favoritename string) error {
	var user User
	user, err := getUser(username)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, ok := user.Favorites[favoritename]; !ok {
		return ErrNotExists
	}

	delete(user.Favorites, favoritename)
	return saveUser(username, user)
}

func GetFavorite(username string, favoritename string) (value Favorite, err error) {
	return value, db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}
		v := bucket.Get([]byte(username))
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

func ListFavorites(username string) (values []Favorite, err error) {
	var user User
	user, err = getUser(username)
	if err != nil {
		log.Println(err)
		return
	}

	if user.Favorites == nil {
		return
	}

	for _, v := range user.Favorites {
		values = append(values, v)
	}

	return
}

func AddUser(username string, password string) error {
	var user User
	rawHash := sha256.Sum256([]byte(password))
	hashpwd := hex.EncodeToString(rawHash[:])
	user.Password = hashpwd

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}

		v := bucket.Get([]byte(username))
		if v != nil {
			return ErrExists
		}

		value, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
			return err
		}

		return bucket.Put([]byte(username), value)
	})
}

func Authenticate(username string, password string) bool {
	var user User
	user, err := getUser(username)
	if err != nil {
		return false
	}

	rawHash := sha256.Sum256([]byte(password))
	hashpwd := hex.EncodeToString(rawHash[:])

	if hashpwd != user.Password {
		return false
	}

	return true
}

func UpdateUserPassword(username string, password string) error {
	var user User
	user, err := getUser(username)
	if err != nil {
		log.Println(err)
		return err
	}

	user.Password = password
	err = saveUser(username, user)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// for admin only
func DeleteUser(username string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}

		v := bucket.Get([]byte(username))
		if v == nil {
			return ErrNotExists
		}

		return bucket.Delete([]byte(username))
	})
}

func ListUsers() (users []string, err error) {
	return users, db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}

		cursor := bucket.Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			users = append(users, string(k))
		}
		return nil
	})
}

func getUser(username string) (user User, err error) {
	return user, db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}
		v := bucket.Get([]byte(username))
		if v == nil {
			return ErrNotExists
		}

		err := json.Unmarshal(v, &user)
		if err != nil {
			log.Println(err)
			return err
		}

		if user.Favorites == nil {
			user.Favorites = make(map[string]Favorite)
		}

		return nil
	})
}

func saveUser(username string, user User) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if bucket == nil {
			return ErrBucketNotFound
		}

		value, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
			return err
		}

		if user.Favorites == nil {
			user.Favorites = make(map[string]Favorite)
		}

		return bucket.Put([]byte(username), value)
	})
}
