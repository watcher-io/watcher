package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/watcher-io/watcher/model"
	"os"
)

type clusterProfileRepo struct {
	conn *badger.DB
}

func NewClusterProfileRepo(db *badger.DB) *clusterProfileRepo {
	return &clusterProfileRepo{
		conn: db,
	}
}

func (r *clusterProfileRepo) FetchAll(
	ctx context.Context,
) (
	[]*model.ClusterProfile,
	error,
) {
	var clusterProfiles []*model.ClusterProfile
	return clusterProfiles, r.conn.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(os.Getenv("CLUSTER_PREFIX"))
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			var clusterInfo model.ClusterProfile
			if err := item.Value(func(v []byte) error {
				return json.Unmarshal(v, &clusterInfo)
			}); err != nil {
				return err
			} else {
				clusterProfiles = append(clusterProfiles, &clusterInfo)
			}
		}
		return nil
	})
}

func (r *clusterProfileRepo) Create(
	ctx context.Context,
	cluster *model.ClusterProfile,
) error {
	if byteData, err := json.Marshal(cluster); err != nil {
		return err
	} else {
		return r.conn.Update(
			func(tx *badger.Txn) error {
				return tx.Set(
					[]byte(fmt.Sprintf("%s_%s", os.Getenv("CLUSTER_PREFIX"), cluster.ID)),
					byteData,
				)
			},
		)
	}
}

func (r *clusterProfileRepo) FetchByID(
	ctx context.Context,
	clusterID string,
) (
	*model.ClusterProfile,
	error,
) {
	var cluster model.ClusterProfile
	return &cluster,
		r.conn.View(
			func(tx *badger.Txn) error {
				if item, err := tx.Get(
					[]byte(fmt.Sprintf("%s_%s", os.Getenv("CLUSTER_PREFIX"),
						clusterID,
					))); err != nil {
					return err
				} else {
					return item.Value(func(v []byte) error {
						return json.Unmarshal(v, &cluster)
					})
				}
			},
		)
}
