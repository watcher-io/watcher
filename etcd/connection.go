package etcd

import (
	"fmt"
	"github.com/aka-achu/watcher/model"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func connect(cfg *model.Cluster) (*clientv3.Client, error) {
	fmt.Println(cfg)
	if !cfg.TLS {
		return clientv3.New(clientv3.Config{
			DialTimeout: 2 * time.Second,
			Endpoints:   cfg.Endpoints,
			TLS:         nil,
			Username:    cfg.Username,
			Password:    cfg.Password,
		})
	} else {
		return nil, nil
	}
}
