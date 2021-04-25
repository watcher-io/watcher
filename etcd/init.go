package etcd

import (
	"context"
	"github.com/watcher-io/watcher/repository"
)

var Store *store

func Initialize() {
	Store = NewStore(
		context.Background(),
		repository.NewClusterProfileRepo(repository.NewDatabase().Conn),
		600,
	)
}
