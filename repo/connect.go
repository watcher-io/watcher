package repo

import (
	"github.com/aka-achu/watcher/logging"
	"go.etcd.io/bbolt"
	"path/filepath"
	"time"
)

var Connection *bbolt.DB

func Initialize() {
	if db, err := bbolt.Open(filepath.Join("data", "watcher.db"),0666, &bbolt.Options{
		Timeout:         1 * time.Second,
	}); err != nil {
		logging.Error.Fatalf(" [DB] Failed to connect to the watcher.db. %v", err)
	} else {
		logging.Info.Printf(" [DB] Successfully established connection with the db")
		Connection = db
	}
}
