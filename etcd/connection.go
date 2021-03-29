package etcd

import (
	"context"
	"github.com/watcher-io/watcher/model"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
	"os"
	"time"
)

func connect(
	ctx context.Context,
	cfg *model.ClusterProfile,
	store model.ObjectStore,
) (
	*clientv3.Client,
	error,
) {

	if cfg.TLS {
		caFilePath, err := store.Fetch(ctx, os.Getenv("STORE_CERT_BUCKET"), cfg.CAFile)
		if err != nil {
			return nil, err
		}
		certFilePath, err := store.Fetch(ctx, os.Getenv("STORE_CERT_BUCKET"), cfg.CertFile)
		if err != nil {
			return nil, err
		}
		keyFilePath, err := store.Fetch(ctx, os.Getenv("STORE_CERT_BUCKET"), cfg.KeyFile)
		if err != nil {
			return nil, err
		}

		//defer os.Remove(caFilePath)
		//defer os.Remove(certFilePath)
		//defer os.Remove(keyFilePath)

		tlsConfig, err := transport.TLSInfo{
			TrustedCAFile: caFilePath,
			CertFile:      certFilePath,
			KeyFile:       keyFilePath,
		}.ClientConfig()
		if err != nil {
			return nil, err
		}

		return clientv3.New(clientv3.Config{
			DialTimeout: 2 * time.Second,
			Endpoints:   cfg.Endpoints,
			TLS:         tlsConfig,
			Username:    cfg.Username,
			Password:    cfg.Password,
		})

	} else {
		return clientv3.New(clientv3.Config{
			DialTimeout: 2 * time.Second,
			Endpoints:   cfg.Endpoints,
			TLS:         nil,
			Username:    cfg.Username,
			Password:    cfg.Password,
		})
	}
}
