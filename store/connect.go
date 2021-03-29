package store

import (
	"crypto/tls"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/watcher-io/watcher/logging"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type objectStore struct {
	Conn *minio.Client
}

func NewObjectStore() *objectStore {
	return &objectStore{Conn: getConnection()}
}

func getConnection() *minio.Client {
	transport := &http.Transport{
		DisableKeepAlives:   false,
		MaxIdleConns:        5,
		MaxIdleConnsPerHost: 5,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}
	if os.Getenv("STORE_SSL") == "true" {
		cer, err := tls.LoadX509KeyPair(
			filepath.Join("cert", "store.crt"),
			filepath.Join("cert", "store.key"),
		)
		if err != nil {
			logging.Error.Fatalf(" [STORE] Failed to load the TLS certificates. Error-%v", err)
		}
		transport.TLSClientConfig = &tls.Config{Certificates: []tls.Certificate{cer}}
	}

	clientConnection, err := minio.New(
		os.Getenv("STORE_URI"),
		&minio.Options{
			Creds:        credentials.NewStaticV4(os.Getenv("STORE_ACCESS_KEY"), os.Getenv("STORE_SECRET_KEY"), ""),
			Secure:       os.Getenv("STORE_SSL") == "true",
			Transport:    transport,
			Region:       "",
			BucketLookup: 0,
			CustomMD5:    nil,
			CustomSHA256: nil,
		})
	if err != nil {
		logging.Error.Fatalf(" [STORE] Failed to connect to the object store. Error-%v", err)
		return nil
	} else {
		logging.Info.Printf(" [STORE] Successfully established connection with the object store.")
		return clientConnection
	}
}
