package repository

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/watcher-io/watcher/logging"
	"os"
)

var conn *badger.DB

type Database struct {
	Conn *badger.DB
}

func NewDatabase() *Database {
	return &Database{Conn: conn}
}

func Initialize() {
	cfg := badger.DefaultOptions("data")
	cfg.CompactL0OnClose = true
	cfg.IndexCacheSize = 128
	cfg.EncryptionKey = []byte(os.Getenv("DB_ENCRYPTION_SECRET"))
	if db, err := badger.Open(cfg); err != nil {
		logging.Error.Fatalf(" [DB] Failed to connect with the db. Error-%v", err)
	} else {
		logging.Info.Printf(" [DB] Successfully established connection with the db.")
		conn = db
	}
}
