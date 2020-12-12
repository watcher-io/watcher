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
	repo model.ClusterProfileRepo,
	ctx context.Context,
) (
	*model.ClusterProfile,
	error,
) {
	profile.ID = uuid.New().String()
	profile.CreatedAt = time.Now().Unix()
	requestTraceID := ctx.Value("trace_id").(string)
	if err := repo.Create(profile, ctx); err != nil {
		logging.Error.Printf(" [DB] TraceID-%s Failed to create cluster profile. Error-%v",
			requestTraceID, err)
		return nil, err
	} else {
		return profile, nil
	}
}

func (*clusterProfileService) FetchByID(
	profileID string,
	repo model.ClusterProfileRepo,
	ctx context.Context,
) (
	*model.ClusterProfile,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	if profile, err := repo.FetchByID(profileID, ctx); err != nil {
		logging.Error.Printf(" [DB] TraceID-%s Failed to fetch the cluster profile. Error-%v",
			requestTraceID, err)
		return nil, err
	} else {
		return profile, nil
	}
}

func (*clusterProfileService) FetchAll(
	repo model.ClusterProfileRepo,
	ctx context.Context,
) (
	[]*model.ClusterProfile,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	if profiles, err := repo.FetchAll(ctx); err != nil {
		logging.Error.Printf(" [DB] TraceID-%s  Failed to fetch cluster profiles. Error-%v",
			requestTraceID, err)
		return nil, err
	} else {
		return profiles, nil
	}
}
