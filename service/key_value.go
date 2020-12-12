package service

import (
	"context"
	"github.com/aka-achu/watcher/etcd"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
)

type kvService struct{}

func NewKVService() *kvService {
	return &kvService{}
}

func (*kvService) Put(
	profileID string,
	kv *model.PutKVRequest,
	r model.ClusterProfileRepo,
	ctx context.Context,
) (
	*model.PutKVResponse,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	conn, err := etcd.Store.Get(r, profileID, ctx)
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
