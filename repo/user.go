package repo

import (
	"encoding/json"
	"fmt"
	"github.com/aka-achu/watcher/model"
	"github.com/dgraph-io/badger/v2"
	"os"
)

// GetAdminDetails, retrieves the profile details of the admin.
func (db *Database) GetAdminDetails() (*model.User, error) {
	var user *model.User
	return user,
		db.Conn.View(
			func(tx *badger.Txn) error {
				if item, err := tx.Get([]byte(fmt.Sprintf("%s_%s", os.Getenv("USER_PREFIX"), "admin"))); err != nil {
					return err
				} else {
					return item.Value(func(v []byte) error {
						return json.Unmarshal(v, user)
					})
				}
			},
		)
}

// SaveUserDetails, updates admin profile details in the Database
func (db *Database) SaveUserDetails(user *model.User) error {
	if byteData, err := json.Marshal(user); err != nil {
		return err
	} else {
		return db.Conn.Update(
			func(tx *badger.Txn) error {
				return tx.Set(
					[]byte(fmt.Sprintf("%s_%s", os.Getenv("USER_PREFIX"), "admin")),
					byteData,
				)
			},
		)
	}
}
