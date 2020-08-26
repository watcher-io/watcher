package repo

import (
	"github.com/aka-achu/watcher/logging"
	"go.etcd.io/bbolt"
	"path/filepath"
	"time"
)

type Repo struct {
	Conn *bbolt.DB
}

// Declaring a global db Connection variable
var DB Repo

// repo.Initialize, initialized the database connection
func Initialize() {
	if db, err := bbolt.Open(filepath.Join("data", "watcher.db"), 0666, &bbolt.Options{
		Timeout: 1 * time.Second,
	}); err != nil {
		logging.Error.Fatalf(" [DB] Failed to connect to the watcher.db. Error-%v", err)
	} else {
		logging.Info.Printf(" [DB] Successfully established connection with the db")
		DB = Repo{
			Conn: db,
		}
	}
}
