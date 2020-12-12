package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aka-achu/watcher/model"
	"github.com/dgraph-io/badger/v2"
)

type userRepo struct {
	conn *badger.DB
}

func NewUserRepo(db *badger.DB) *userRepo {
	return &userRepo{
		conn: db,
	}
}

func (r *userRepo) Create(user *model.User, ctx context.Context) error {
	if byteData, err := json.Marshal(user); err != nil {
		return err
	} else {
		return r.conn.Update(
			func(tx *badger.Txn) error {
				return tx.Set(
					[]byte(fmt.Sprintf("%s_%s", user.Prefix(), "admin")),
					byteData,
				)
			},
		)
	}
}

func (r *userRepo) Fetch(userName string, ctx context.Context) (*model.User, error) {
	var user model.User
	return &user,
		r.conn.View(
			func(tx *badger.Txn) error {
				if item, err := tx.Get([]byte(fmt.Sprintf("%s_%s", user.Prefix(), userName)));
					err != nil {
					return err
				} else {
					return item.Value(func(v []byte) error {
						return json.Unmarshal(v, &user)
					})
				}
			},
		)
}

