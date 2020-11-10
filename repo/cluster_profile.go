package repo

import (
	"encoding/json"
	"errors"
	"github.com/aka-achu/watcher/model"
	"go.etcd.io/bbolt"
	"os"
)

// GetClusterProfiles, iterated the ${CLUSTER_PROFILE_BUCKET} bucket and
// fetches all the cluster profiles present in the bucket.
func (db *Database)GetClusterProfiles() ([]*model.ClusterProfile, error) {
	var clusterProfiles []*model.ClusterProfile
	return clusterProfiles, db.Conn.View(func(tx *bbolt.Tx) error {
		// Creating cursor object of the bucket for iteration
		c := tx.Bucket([]byte(os.Getenv("CLUSTER_PROFILE_BUCKET"))).Cursor()
		// Iterating the cursor object
		for clusterID, cluster := c.First(); clusterID != nil; clusterID, cluster = c.Next() {
			// If a valid key value pair is found then decode profile data into model.ClusterProfile object
			var clusterInfo model.ClusterProfile
			if err := json.Unmarshal(cluster, &clusterInfo); err != nil {
				return err
			}
			// Append the cluster profile info with the pre-declared array
			clusterProfiles = append(clusterProfiles, &clusterInfo)
		}
		return nil
	})
}

// CreateClusterProfile, creates a cluster profile inside ${CLUSTER_PROFILE_BUCKET} bucket, given
// a validated *model.ClusterProfile object
func (db *Database)CreateClusterProfile(cluster *model.ClusterProfile) error {
	if byteData, err := json.Marshal(cluster); err != nil {
		return err
	} else {
		return db.Conn.Update(
			func(tx *bbolt.Tx) error {
				return tx.Bucket([]byte(os.Getenv("CLUSTER_PROFILE_BUCKET"))).Put([]byte(cluster.ID), byteData)
			},
		)
	}
}

// GetClusterInfoByID, return a model.ClusterProfile object containing cluster details of having the requested
// id as the ClusterID
func (db *Database)GetClusterInfoByID(clusterID string) (*model.ClusterProfile, error) {
	var cluster model.ClusterProfile
	return &cluster,
		db.Conn.View(
			func(tx *bbolt.Tx) error {
				byteData := tx.Bucket([]byte(os.Getenv("CLUSTER_PROFILE_BUCKET"))).Get([]byte(clusterID))
				if len(byteData) == 0 {
					return errors.New("requested cluster does not exist in the Database")
				} else {
					return json.Unmarshal(
						byteData,
						&cluster,
					)
				}
			},
		)
}
