package repo

import (
	"encoding/json"
	"github.com/aka-achu/watcher/model"
	"go.etcd.io/bbolt"
	"os"
)

// GetUserDetails, retrieves the profile details of the admin.
func GetUserDetails() (*model.User, error) {
	var user model.User
	return &user,
		Connection.View(
			func(tx *bbolt.Tx) error {
				byteData := tx.Bucket([]byte(os.Getenv("USER_PROFILE_BUCKET"))).Get([]byte("admin"))
				if len(byteData) == 0 {
					return nil
				} else {
					return json.Unmarshal(
						byteData,
						&user,
					)
				}
			},
		)
}

// SaveUserDetails, updates admin profile details in the database
func SaveUserDetails(user *model.User) error {
	if byteData, err := json.Marshal(user); err != nil {
		return err
	} else {
		return Connection.Update(
			func(tx *bbolt.Tx) error {
				return tx.Bucket([]byte(os.Getenv("USER_PROFILE_BUCKET"))).Put([]byte("admin"), byteData)
			},
		)
	}
}
