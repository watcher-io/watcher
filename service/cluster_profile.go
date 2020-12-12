package service

import (
	"context"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/google/uuid"
	"time"
)

type clusterProfileService struct{}

func NewClusterProfileService() *clusterProfileService {
	return &clusterProfileService{}
}

func (*clusterProfileService) Create(
	profile *model.ClusterProfile,
	r model.ClusterProfileRepo,
	ctx context.Context,
) (
	*model.ClusterProfile,
	error,
) {
	profile.ID = uuid.New().String()
	profile.CreatedAt = time.Now().Unix()
	requestTraceID := ctx.Value("trace_id").(string)
	if err := r.Create(profile, ctx); err != nil {
		logging.Error.Printf(" [DB] Failed to create cluster profile. Error-%v TraceID-%s",
			err, requestTraceID)
		return nil, err
	} else {
		logging.Info.Printf(" [DB] Cluster profile created. TraceID-%s",
			requestTraceID)
		return profile, nil
	}
}

func (*clusterProfileService) FetchByID(
	profileID string,
	r model.ClusterProfileRepo,
	ctx context.Context,
) (
	*model.ClusterProfile,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	if profile, err := r.FetchByID(profileID, ctx); err != nil {
		logging.Error.Printf(" [DB] Failed to fetch the cluster profile. Error-%v TraceID-%s",
			err, requestTraceID)
		return nil, err
	} else {
		logging.Info.Printf(" [DB] Cluster profile retrieved. TraceID-%s",
			requestTraceID)
		return profile, nil
	}
}

func (*clusterProfileService) FetchAll(
	r model.ClusterProfileRepo,
	ctx context.Context,
) (
	[]*model.ClusterProfile,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	if profiles, err := r.FetchAll(ctx); err != nil {
		logging.Error.Printf(" [DB] Failed to fetch cluster profiles. Error-%v TraceID-%s",
			err, requestTraceID)
		return nil, err
	} else {
		logging.Info.Printf(" [DB] Cluster profiles retrived. TraceID-%s",
			requestTraceID)
		return profiles, nil
	}
}
