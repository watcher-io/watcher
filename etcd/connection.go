package etcd

import (
	"fmt"
	"github.com/aka-achu/watcher/model"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func connect(cfg *model.ClusterProfile) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		DialTimeout: 2 * time.Second,
		Endpoints:   cfg.Endpoints,
		TLS:         nil,
		Username:    cfg.Username,
		Password:    cfg.Password,
	})
}

func C(cfg *model.ClusterProfile) {
	fmt.Println(connect(cfg))
}
