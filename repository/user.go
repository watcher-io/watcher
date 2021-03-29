package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/watcher-io/watcher/model"
)

type userRepo struct {
	conn *badger.DB
}

func NewUserRepo(db *badger.DB) *userRepo {
	return &userRepo{
		conn: db,
	}
}

func (r *userRepo) Create(
	ctx context.Context,
	user *model.User,
) error {
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

func (r *userRepo) Fetch(
	ctx context.Context,
	userName string,
) (
	*model.User,
	error,
) {
	var user model.User
	return &user,
		r.conn.View(
			func(tx *badger.Txn) error {
				if item, err := tx.Get([]byte(fmt.Sprintf("%s_%s", user.Prefix(), userName))); err != nil {
					return err
				} else {
					return item.Value(func(v []byte) error {
						return json.Unmarshal(v, &user)
					})
				}
			},
		)
}
