package service

import (
	"context"
	"github.com/aka-achu/watcher/etcd"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
)

type dashboardService struct{}

func NewDashboardService() *dashboardService {
	return &dashboardService{}
}

func (*dashboardService) ViewCluster(
	profileID string,
	repo model.ClusterProfileRepo,
	ctx context.Context,
) (
	*model.Cluster,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	conn, err := etcd.Store.Get(repo, profileID, ctx)
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
