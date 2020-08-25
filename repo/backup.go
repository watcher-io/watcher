package repo

import (
	"go.etcd.io/bbolt"
	"net/http"
	"strconv"
)

// BackupRepo, writes db snapshot to a response writer
func BackupRepo(w http.ResponseWriter) error {
	return Connection.View(func(tx *bbolt.Tx) error {
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})
}
