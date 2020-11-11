package repo

import (
	"encoding/json"
	"fmt"
	"github.com/aka-achu/watcher/model"
	"github.com/dgraph-io/badger/v2"
	"os"
)

// GetClusterProfiles, iterated the ${CLUSTER_PROFILE_BUCKET} bucket and
// fetches all the cluster profiles present in the bucket.
func (db *Database) GetClusterProfiles() ([]*model.ClusterProfile, error) {
	var clusterProfiles []*model.ClusterProfile
	return clusterProfiles, db.Conn.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(os.Getenv("CLUSTER_PREFIX"))
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			var clusterInfo *model.ClusterProfile
			if err := item.Value(func(v []byte) error {
				return json.Unmarshal(v, clusterInfo)
			}); err != nil {
				return err
			} else {
				clusterProfiles = append(clusterProfiles, clusterInfo)
			}
		}
		return nil
	})
}

// CreateClusterProfile, creates a cluster profile inside ${CLUSTER_PROFILE_BUCKET} bucket, given
// a validated *model.ClusterProfile object
func (db *Database) CreateClusterProfile(cluster *model.ClusterProfile) error {
	if byteData, err := json.Marshal(cluster); err != nil {
		return err
	} else {
		return db.Conn.Update(
			func(tx *badger.Txn) error {
				return tx.Set(
					[]byte(fmt.Sprintf("%s_%s", os.Getenv("CLUSTER_PREFIX"), cluster.ID)),
					byteData,
				)
			},
		)
	}
}

// GetClusterInfoByID, return a model.ClusterProfile object containing cluster details of having the requested
// id as the ClusterID
func (db *Database) GetClusterInfoByID(clusterID string) (*model.ClusterProfile, error) {
	var cluster *model.ClusterProfile
	return cluster,
		db.Conn.View(
			func(tx *badger.Txn) error {
				if item, err := tx.Get([]byte(fmt.Sprintf("%s_%s", os.Getenv("CLUSTER_PREFIX"), clusterID))); err != nil {
					return err
				} else {
					return item.Value(func(v []byte) error {
						return json.Unmarshal(v,cluster)
					})
				}
			},
		)
}
