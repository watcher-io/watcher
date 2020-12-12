package service

import (
	"context"
	"github.com/aka-achu/watcher/etcd"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/validator"
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
	if err := validator.Validate.Struct(kv); err != nil {
		logging.Error.Printf(" [APP] Failed to validate the request body for required fields. Error-%v TraceID-%s",
			err, requestTraceID)
		return nil, err
	}
	conn, err := etcd.Store.Get(r, profileID, ctx)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to establish connection with the cluster. Error-%v TraceID-%s ClusterProfileID-%s",
			err, requestTraceID, profileID)
		return nil, err
	}

	putResponse, err := etcd.PutKV(ctx, conn, kv)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to put the key-value pair in the cluster. Error-%v TraceID-%s ClusterProfileID-%s",
			err, requestTraceID, profileID)
		return nil, err
	} else {
		logging.Info.Printf(" [APP] Successfully stored the kv pair. TraceID-%s ClusterProfileID-%s",
			requestTraceID, profileID)
		return putResponse, nil
	}
}
