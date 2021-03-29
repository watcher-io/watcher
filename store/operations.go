package store

import (
	"context"
	"github.com/minio/minio-go/v7"
	"os"
	"path/filepath"
)

func (s *objectStore) Upload(
	ctx context.Context,
	bucketName string,
	objectName string,
	filePath string,
) error {
	var putObjectOptions minio.PutObjectOptions
	_, err := s.Conn.FPutObject(
		ctx,
		bucketName,
		objectName,
		filePath,
		putObjectOptions,
	)
	return err
}

func (s *objectStore) Fetch(
	ctx context.Context,
	bucketName string,
	objectName string,
) (
	string,
	error,
) {
	filePath := filepath.Join(os.TempDir(), objectName)
	var getObjectOptions minio.GetObjectOptions
	return filePath, s.Conn.FGetObject(
		ctx,
		bucketName,
		objectName,
		filePath,
		getObjectOptions,
	)

}
