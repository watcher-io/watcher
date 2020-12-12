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
	r model.ClusterProfileRepo,
	ctx context.Context,
) (
	*model.Cluster,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	conn, err := etcd.Store.Get(r, profileID, ctx)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to establish connection with the cluster. Error-%v TraceID-%s ClusterProfileID-%s",
			err, requestTraceID, profileID)
		return nil, err
	}

	clusterState, err := etcd.FetchMember(ctx, conn)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to fetch cluster info. Error-%v TraceID-%s ClusterProfileID-%s",
			err, requestTraceID, profileID)
		return nil, err
	}
	return clusterState, nil
}
