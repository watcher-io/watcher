package repo

import (
	"encoding/json"
	"github.com/aka-achu/watcher/model"
	"go.etcd.io/bbolt"
	"os"
)

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
						tx.Bucket([]byte(os.Getenv("USER_PROFILE_BUCKET"))).Get([]byte("admin")),
						&user,
					)
				}
			},
		)
}

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
