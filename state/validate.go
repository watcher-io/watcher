package state

import (
	"github.com/aka-achu/watcher/logging"
	"go.etcd.io/bbolt"
	"os"
	"path/filepath"
	"time"
)

// Validate, validates the state of the application.
func Validate() {
	// Checking for the existence of the data directory
	// If the data dir does not exist and create the data directory
	// and initializes the buckets for "watcher.db".
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		if err := os.Mkdir("data", 0755); err != nil {
			logging.Error.Fatalf("[APP] Failed to create data directory. %v", err)
		} else {
			logging.Info.Printf(" [APP] Successfully created the data directory.")
		}
	}

	// Initializing the database buckets
	if err := InitBuckets(); err != nil {
		logging.Error.Fatalf(" [DB] Failed to initialize the buckets. %v", err)
	} else {
		logging.Info.Printf(" [DB] Successfully initialzied the buckets.")
	}
}

// InitBuckets, creates the required buckets if the buckets are not present
// in the "watcher.db" file.
func InitBuckets() error {
	if db, err := bbolt.Open(filepath.Join("data", "watcher.db"), 0666, &bbolt.Options{
		Timeout: 1 * time.Second,
	}); err != nil {
		return err
	} else {
		defer db.Close()

		// Creating a manual db transaction for creation of buckets
		if tx, err := db.Begin(true); err != nil {
			return err
		} else {

			// Creating the required buckets in the same transaction
			if _, err := tx.CreateBucketIfNotExists([]byte(os.Getenv("USER_PROFILE_BUCKET"))); err != nil {
				_ = tx.Rollback()
				return err
			}
			if _, err := tx.CreateBucketIfNotExists([]byte(os.Getenv("CLUSTER_PROFILE_BUCKET"))); err != nil {
				_ = tx.Rollback()
				return err
			}
			if _, err := tx.CreateBucketIfNotExists([]byte(os.Getenv("CERTIFICATE_BUCKET"))); err != nil {
				_ = tx.Rollback()
				return err
			}

			// Committing the changes of the current transaction
			if err := tx.Commit(); err != nil {
				return err
			}
		}
		return nil
	}
}
