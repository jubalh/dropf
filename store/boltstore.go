// The functions in filestore are used to store key value value pairs
// of the files.
// A file will be saved in a directory set in the configuration file,
// and will be available for users under a certain link.
// We store the connection between that link to the actual path of the
// file in a Bolt database.
package store

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
	"github.com/satori/go.uuid"
)

type BoltStore struct {
	Name string
}

// updateItem updates an item in the Bolt database.
func updateItem(db bolt.DB, key string, value string) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("default"))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(key), []byte(value))
		return err
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// getItem retrieves the value in the Bolt database by its key.
func getItem(db bolt.DB, key string) (value string, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("default"))
		if bucket == nil {
			return fmt.Errorf("Bucket q not found")
		}

		value = string(bucket.Get([]byte(key)))

		return nil
	})

	return
}

// removeItem removes an item by its key from the Bolt database.
func removeItem(db bolt.DB, key string) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("default"))
		if bucket == nil {
			return fmt.Errorf("Bucket q not found")
		}

		err := bucket.Delete([]byte(key))
		return err
	})
	return
}

func (BoltStore) Store(file *File) (uuid.UUID, error) {
}
