package etcd

import (
	"context"
	"fmt"
	"github.com/aka-achu/watcher/model"
	"go.etcd.io/etcd/clientv3"
	"sync"
	"time"
)

type connection struct {
	client     *clientv3.Client
	lastAccess int64
}

type store struct {
	m map[string]*connection
	l sync.Mutex
}

func NewStore(
	maxTTL int,
) *store {
	s := &store{m: make(map[string]*connection)}
	go func() {
		for now := range time.Tick(time.Minute) {
			s.l.Lock()
			for clusterID, cluster := range s.m {
				if now.Unix()-cluster.lastAccess > int64(maxTTL) {
					_ = s.m[clusterID].client.Close()
					delete(s.m, clusterID)
				}
			}
			s.l.Unlock()
		}
	}()
	return s
}

func (s *store) Get(
	repo model.ClusterProfileRepo,
	profileID string,
	ctx context.Context,
) (
	*clientv3.Client,
	error,
) {
	s.l.Lock()
	defer s.l.Unlock()

	if cluster, ok := s.m[profileID]; ok {
		cluster.lastAccess = time.Now().Unix()
		return cluster.client, nil
	} else {
		clusterInfo, err := repo.FetchByID(profileID, ctx)
		if err != nil {
			fmt.Println("Failed to fetch the cluster profile", err)
			return nil, err
		}
		conn, err := connect(clusterInfo)
		if err != nil {
			return nil, err
		}
		s.m[profileID] = &connection{
			client:     conn,
			lastAccess: time.Now().Unix(),
		}
		return conn, nil
	}
}
