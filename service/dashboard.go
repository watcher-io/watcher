package service

import (
	"context"
	"github.com/watcher-io/watcher/etcd"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/model"
)

type dashboardService struct {
	repo  model.ClusterProfileRepo
	store model.ObjectStore
}

func NewDashboardService(
	repo model.ClusterProfileRepo,
	store model.ObjectStore,
) *dashboardService {
	return &dashboardService{repo, store}
}

func (s *dashboardService) ViewCluster(
	ctx context.Context,
	profileID string,
) (
	*model.Cluster,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	conn, err := etcd.Store.Get(s.repo, profileID, s.store, ctx)
	if err != nil {
		logging.Error.Printf(" [APP] TraceID-%s Failed to establish connection with the cluster. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return nil, err
	}
	clusterState, err := etcd.FetchMember(ctx, conn)
	if err != nil {
		logging.Error.Printf(" [APP] TraceID-%s Failed to fetch cluster info. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return nil, err
	}
	return clusterState, nil
}
