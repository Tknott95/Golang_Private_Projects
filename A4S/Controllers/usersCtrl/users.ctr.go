package users

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"
	"golang.org/x/crypto/bcrypt"
)

var (
	// DB is the reference to our DB, which contains our user data.
	DB = newDB()

	// ErrUserAlreadyExists is the error thrown when a user attempts to create
	// a new user in the DB with a duplicate username.
	ErrUserAlreadyExists = errors.New("users: username already exists")

	// ErrUserNotFound is the error thrown when a user can't be found in the
	// database.
	ErrUserNotFound = errors.New("users: user not found")
)

// Store is a reference to our BoltDB instance that contains two seperate
// internal stores: a user store, and a session store.
type Store struct {
	DB       *bolt.DB
	Users    string
	Sessions string
}

// newDB is a convenience method to initalize our DB.
func newDB() *Store {
	// Create or open the database
	db, err := bolt.Open("users.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	// Create the Users bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Users"))
		if err != nil {
			return err
		}
		return nil
	})

	// Create the Sessions bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Sessions"))
		if err != nil {
			return err
		}
		return nil
	})

	return &Store{
		DB:       db,
		Users:    "Users",
		Sessions: "Sessions",
	}
}

// NewUser accepts a username and password and creates a new user in our DB
// from it.
func NewUser(username string, password string) error {
	err := exists(username)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return DB.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Users))
		return b.Put([]byte(username), hashedPassword)
	})
}

// AuthenticateUser accepts a username and password, and checks that the given
// password matches the hashed password. It returns nil on success, and an
// error on failure.
func AuthenticateUser(username string, password string) error {
	var hashedPassword []byte
	DB.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Users))
		hashedPassword = b.Get([]byte(username))
		return nil
	})
	if hashedPassword == nil {
		return ErrUserNotFound
	}
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// OverrideOldPassword overrides the old password with the new password. For
// use when resetting passwords.
func OverrideOldPassword(username string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return DB.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Users))
		return b.Put([]byte(username), hashedPassword)
	})
}

// exists is an internal utility function for ensuring the usernames are
// unique.
func exists(username string) error {
	var result []byte
	DB.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Users))
		result = b.Get([]byte(username))
		return nil
	})
	if result != nil {
		return ErrUserAlreadyExists
	}
	return nil
}
