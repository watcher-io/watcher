package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aka-achu/watcher/model"
	"go.etcd.io/bbolt"
	"os"
)

func GetClusterProfiles() ([]*model.Cluster, error) {
	var clusterProfiles []*model.Cluster
	return clusterProfiles, Connection.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket([]byte(os.Getenv("CLUSTER_PROFILE_BUCKET"))).Cursor()
		for clusterID, cluster := c.First(); clusterID != nil; clusterID, cluster = c.Next() {
			var clusterInfo model.Cluster
			fmt.Println("Came here")
			fmt.Printf("%s:%s", clusterID, cluster)
			if err := json.Unmarshal(cluster, &clusterInfo); err != nil {
				return err
			}
			clusterProfiles = append(clusterProfiles, &clusterInfo)
		}
		return nil
	})
}

func CreateClusterProfile(cluster *model.Cluster) error {
	if byteData, err := json.Marshal(cluster); err != nil {
		return err
	} else {
		return Connection.Update(
			func(tx *bbolt.Tx) error {
				return tx.Bucket([]byte(os.Getenv("CLUSTER_PROFILE_BUCKET"))).Put([]byte(cluster.ID), byteData)
			},
		)
	}
}

func GetClusterInfoByID(clusterID string) (*model.Cluster, error) {
	var cluster model.Cluster
	return &cluster,
		Connection.View(
			func(tx *bbolt.Tx) error {
				byteData := tx.Bucket([]byte(os.Getenv("CLUSTER_PROFILE_BUCKET"))).Get([]byte(clusterID))
				if len(byteData) == 0 {
					return errors.New("requested cluster does not exist in the database")
				} else {
					return json.Unmarshal(
						byteData,
						&cluster,
					)
				}
			},
		)
}
