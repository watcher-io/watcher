package repo

import (
	"github.com/aka-achu/watcher/logging"
	"go.etcd.io/bbolt"
	"path/filepath"
	"time"
)

type Database struct {
	Conn *bbolt.DB
}

// NewRepo, initialized a db object containing a Database connection
func NewDatabase() *Database {
	return &Database{Conn: getConnection()}
}

// getConnection, establishes a Database connection
func getConnection() *bbolt.DB {
	if db, err := bbolt.Open(
		filepath.Join("data", "watcher.db"),
		0666,
		&bbolt.Options{
			Timeout: 1 * time.Second,
		}); err != nil {
		logging.Error.Fatalf(" [DB] Failed to connect to the watcher.db. Error-%v", err)
		return nil
	} else {
		logging.Info.Printf(" [DB] Successfully established connection with the db")
		return db
	}
}
