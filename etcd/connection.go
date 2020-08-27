package etcd

import (
	"github.com/aka-achu/watcher/model"
	"go.etcd.io/etcd/clientv3"
)

func connect(cfg *model.Cluster) (*clientv3.Client, error) {
	if !cfg.TLS {
		return clientv3.New(clientv3.Config{
			Endpoints:            cfg.Endpoints,
			TLS:                  nil,
			Username:             cfg.Username,
			Password:             cfg.Password,
		})
	} else {
		return nil,nil
	}

}
