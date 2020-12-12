package repo

import (
	"github.com/aka-achu/watcher/logging"
	"github.com/dgraph-io/badger/v2"
	"os"
)

type Database struct {
	Conn *badger.DB
}

func NewDatabase() *Database {
	return &Database{Conn: getConnection()}
}

func getConnection() *badger.DB {
	if db, err := badger.Open(
		badger.DefaultOptions("data").
			WithEncryptionKey([]byte(os.Getenv("DB_ENCRYPTION_SECRET"))).WithTruncate(true),
	); err != nil {
		logging.Error.Fatalf(" [DB] Failed to connect to the watcher.db. Error-%v", err)
		return nil
	} else {
		logging.Info.Printf(" [DB] Successfully established connection with the db")
		return db
	}
}
