package etcd

import (
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/repo"
	"go.etcd.io/etcd/clientv3"
	"sync"
	"time"
)

type Cluster struct {
	Connection *clientv3.Client
	lastAccess int64
}

type ConnectionStore struct {
	m map[string]*Cluster
	l sync.Mutex
}

func New(maxTTL int) (cs *ConnectionStore) {
	cs = &ConnectionStore{m: make(map[string]*Cluster)}
	go func() {
		for now := range time.Tick(time.Minute) {
			cs.l.Lock()
			for clusterID, cluster := range cs.m {
				if now.Unix()-cluster.lastAccess > int64(maxTTL) {
					_ = cs.m[clusterID].Connection.Close()
					delete(cs.m, clusterID)
				}
			}
			cs.l.Unlock()
		}
	}()
	return
}

func (cs *ConnectionStore) Get(db *repo.Database, clusterProfileID string) (*clientv3.Client, error) {
	cs.l.Lock()
	defer cs.l.Unlock()
	if cluster, ok := cs.m[clusterProfileID]; ok {
		cluster.lastAccess = time.Now().Unix()
		return cluster.Connection, nil
	} else {
		clusterInfo, err := db.GetClusterInfoByID(clusterProfileID)
		if err != nil {
			logging.Error.Printf(" [APP] Failed to fetch the cluster details for connection. Error-%v ClusterID-%s", err, clusterProfileID)
			return nil, err
		}
		conn, err := connect(clusterInfo)
		if err != nil {
			logging.Error.Printf(" [APP] Failed to connect to the requested cluster. Error-%v ClusterID-%s", err, clusterProfileID)
			return nil, err
		}
		logging.Info.Printf(" [APP] Successfully connected to the requested cluster. ClusterID-%s", clusterProfileID)
		cs.m[clusterProfileID] = &Cluster{
			Connection: conn,
			lastAccess: time.Now().Unix(),
		}
		return conn, nil
	}
}
