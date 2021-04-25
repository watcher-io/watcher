package service

import (
	"context"
	"github.com/watcher-io/watcher/etcd"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/model"
)

type kvService struct {
	repo  model.ClusterProfileRepo
	store model.ObjectStore
}

func NewKVService(
	repo model.ClusterProfileRepo,
	store model.ObjectStore,
) *kvService {
	return &kvService{repo, store}
}

func (s *kvService) Put(
	ctx context.Context,
	profileID string,
	kv *model.PutKVRequest,
) (
	*model.PutKVResponse,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	conn, err := etcd.Store.Get(s.repo, profileID, s.store, ctx)
	if err != nil {
		logging.Error.Printf(" [APP]  TraceID-%s Failed to establish connection with the cluster. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return nil, err
	}
	putResponse, err := etcd.PutKV(ctx, conn, kv)
	if err != nil {
		logging.Error.Printf(" [APP] TraceID-%s Failed to put the key-value pair in the cluster. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return nil, err
	} else {
		logging.Info.Printf(" [APP] TraceID-%s Successfully stored the kv pair. ClusterProfileID-%s",
			requestTraceID, profileID)
		return putResponse, nil
	}
}

func (s *kvService) Get(
	ctx context.Context,
	profileID string,
	kv *model.GetKVRequest,
) (
	*model.GetKVResponse,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	conn, err := etcd.Store.Get(s.repo, profileID, s.store, ctx)
	if err != nil {
		logging.Error.Printf(" [APP]  TraceID-%s Failed to establish connection with the cluster. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return nil, err
	}
	getResponse, err := etcd.GetKV(ctx, conn, kv)
	if err != nil {
		logging.Error.Printf(" [APP] TraceID-%s Failed to get the key-value from the cluster. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return nil, err
	} else {
		logging.Info.Printf(" [APP] TraceID-%s Successfully fetched the kv pair(s). ClusterProfileID-%s",
			requestTraceID, profileID)
		return getResponse, nil
	}
}

func (s *kvService) Delete(
	ctx context.Context,
	profileID string,
	kv *model.DeleteKVRequest,
) (
	*model.DeleteKVResponse,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	conn, err := etcd.Store.Get(s.repo, profileID, s.store, ctx)
	if err != nil {
		logging.Error.Printf(" [APP]  TraceID-%s Failed to establish connection with the cluster. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return nil, err
	}
	deleteResponse, err := etcd.DeleteKV(ctx, conn, kv)
	if err != nil {
		logging.Error.Printf(" [APP] TraceID-%s Failed to delete the key-value from the cluster. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return nil, err
	} else {
		logging.Info.Printf(" [APP] TraceID-%s Successfully delete requested kv pair(s). ClusterProfileID-%s",
			requestTraceID, profileID)
		return deleteResponse, nil
	}
}
