package repo

import (
	"go.etcd.io/bbolt"
	"net/http"
	"strconv"
)

// BackupRepo, writes db snapshot to a response writer
func (db *Repo) BackupRepo(w http.ResponseWriter) error {
	return db.Conn.View(func(tx *bbolt.Tx) error {
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})
}
