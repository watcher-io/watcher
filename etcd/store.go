package etcd

import (
	"context"
	"fmt"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/model"
	"go.etcd.io/etcd/clientv3"
	"os"
	"path/filepath"
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
	ctx context.Context,
	repo model.ClusterProfileRepo,
	maxTTL int,
) *store {
	s := &store{m: make(map[string]*connection)}
	go func() {
		for now := range time.Tick(time.Minute) {
			s.l.Lock()
			for clusterID, cluster := range s.m {
				if now.Unix()-cluster.lastAccess > int64(maxTTL) {
					_ = s.m[clusterID].client.Close()
					//todo clear the certificate files associated with the expired connection
					if profile, err := repo.FetchByID(ctx, clusterID); err != nil {
						logging.Error.Fatalf(" [ETCD] Failed to fetch profile details while flushing the certificates. Error-%v",
							err)
					} else {
						fmt.Println("Flushing the files")
						_ = os.Remove(filepath.Join(os.TempDir(), profile.CertFile))
						_ = os.Remove(filepath.Join(os.TempDir(), profile.CAFile))
						_ = os.Remove(filepath.Join(os.TempDir(), profile.KeyFile))
					}
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
	store model.ObjectStore,
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
		clusterInfo, err := repo.FetchByID(ctx, profileID)
		if err != nil {
			return nil, err
		}
		conn, err := connect(ctx, clusterInfo, store)
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
