package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/model"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type clusterProfileService struct {
	repo  model.ClusterProfileRepo
	store model.ObjectStore
}

func NewClusterProfileService(
	repo model.ClusterProfileRepo,
	store model.ObjectStore,
) *clusterProfileService {
	return &clusterProfileService{repo, store}
}

func (s *clusterProfileService) Create(
	ctx context.Context,
	profile *model.ClusterProfile,
) (
	*model.ClusterProfile,
	error,
) {
	profile.ID = uuid.New().String()
	profile.CreatedAt = time.Now().Unix()
	requestTraceID := ctx.Value("trace_id").(string)
	if err := s.repo.Create(ctx, profile); err != nil {
		logging.Error.Printf(" [DB] TraceID-%s Failed to create cluster profile. Error-%v",
			requestTraceID, err)
		return nil, err
	} else {
		return profile, nil
	}
}

func (s *clusterProfileService) FetchByID(
	ctx context.Context,
	profileID string,
) (
	*model.ClusterProfile,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	if profile, err := s.repo.FetchByID(ctx, profileID); err != nil {
		logging.Error.Printf(" [DB] TraceID-%s Failed to fetch the cluster profile. Error-%v",
			requestTraceID, err)
		return nil, err
	} else {
		return profile, nil
	}
}

func (s *clusterProfileService) FetchAll(
	ctx context.Context,
) (
	[]*model.ClusterProfile,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	if profiles, err := s.repo.FetchAll(ctx); err != nil {
		logging.Error.Printf(" [DB] TraceID-%s  Failed to fetch cluster profiles. Error-%v",
			requestTraceID, err)
		return nil, err
	} else {
		return profiles, nil
	}
}

func (s *clusterProfileService) UploadCertificate(
	ctx context.Context,
	r *http.Request,
) (
	map[string]string,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	// maxMemory 32MB
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		logging.Error.Printf(" [APP] TraceID-%s  Failed to parse multipart form. Error-%v",
			requestTraceID, err)
		return nil, err
	}

	tempCAFile, err := processFileUpload(ctx, r.MultipartForm.File["ca"][0], s.store)
	if err != nil {
		logging.Error.Printf(" [APP] TraceID-%s  Failed to process ca file upload. Error-%v",
			requestTraceID, err)
		return nil, err
	}
	tempCertFile, err := processFileUpload(ctx, r.MultipartForm.File["cert"][0], s.store)
	if err != nil {
		logging.Error.Printf(" [APP] TraceID-%s  Failed to process cert file upload. Error-%v",
			requestTraceID, err)
		return nil, err
	}
	tempKeyFile, err := processFileUpload(ctx, r.MultipartForm.File["key"][0], s.store)
	if err != nil {
		logging.Error.Printf(" [APP] TraceID-%s  Failed to process key file upload. Error-%v",
			requestTraceID, err)
		return nil, err
	}

	return map[string]string{
		"ca":   tempCAFile,
		"cert": tempCertFile,
		"key":  tempKeyFile,
	}, nil
}

func processFileUpload(
	ctx context.Context,
	fh *multipart.FileHeader,
	store model.ObjectStore,
) (
	string,
	error,
) {
	file, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	byteData, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	tempFileName := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(fh.Filename))
	tempFilePath := filepath.Join(os.TempDir(), tempFileName)
	err = ioutil.WriteFile(tempFilePath, byteData, 0666)
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFilePath)
	if err := store.Upload(ctx, os.Getenv("STORE_CERT_BUCKET"), tempFileName, tempFilePath); err != nil {
		return "", err
	}
	return tempFileName, nil
}
