package service

import (
	"context"
	"github.com/watcher-io/watcher/etcd"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/model"
)

type maintenanceService struct {
	repo  model.ClusterProfileRepo
	store model.ObjectStore
}

func NewMaintenanceService(
	repo model.ClusterProfileRepo,
	store model.ObjectStore,
) *maintenanceService {
	return &maintenanceService{repo, store}
}

func (s *maintenanceService) DisarmAlarm(
	ctx context.Context,
	profileID string,
	alarm *model.DisarmAlarmRequest,
) (
	*model.DisarmAlarmResponse,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	conn, err := etcd.Store.Get(s.repo, profileID, s.store, ctx)
	if err != nil {
		logging.Error.Printf(" [APP]  TraceID-%s Failed to establish connection with the cluster. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return nil, err
	}
	disarmAlarmResponse, err := etcd.DisarmAlarm(ctx, conn, alarm)
	if err != nil {
		logging.Error.Printf(" [APP] TraceID-%s Failed to disarm the alarm for the requested node. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return nil, err
	} else {
		logging.Info.Printf(" [APP] TraceID-%s Successfully disarmed the alarm for the requested node. ClusterProfileID-%s",
			requestTraceID, profileID)
		return disarmAlarmResponse, nil
	}
}

func (s *maintenanceService) Defragment(
	ctx context.Context,
	profileID string,
	defragment *model.DefragmentRequest,
) (
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	conn, err := etcd.Store.Get(s.repo, profileID, s.store, ctx)
	if err != nil {
		logging.Error.Printf(" [APP]  TraceID-%s Failed to establish connection with the cluster. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return err
	}
	err = etcd.Defragment(ctx, conn, defragment)
	if err != nil {
		logging.Error.Printf(" [APP] TraceID-%s Failed to de-fragment the requested node. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return err
	} else {
		logging.Info.Printf(" [APP] TraceID-%s Successfully de-fragmented the requested node. ClusterProfileID-%s",
			requestTraceID, profileID)
		return nil
	}
}

func (s *maintenanceService) Snapshot(
	ctx context.Context,
	profileID string,
) (
	string,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	conn, err := etcd.Store.Get(s.repo, profileID, s.store, ctx)
	if err != nil {
		logging.Error.Printf(" [APP]  TraceID-%s Failed to establish connection with the cluster. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return "", err
	}
	snapshotFile, err := etcd.Snapshot(ctx, conn)
	if err != nil {
		logging.Error.Printf(" [APP] TraceID-%s Failed to take snapshot of the cluster. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return "", err
	} else {
		logging.Info.Printf(" [APP] TraceID-%s Successfully took a snapshot of the cluster. ClusterProfileID-%s",
			requestTraceID, profileID)
		return snapshotFile, nil
	}
}

func (s *maintenanceService) MoveLeader(
	ctx context.Context,
	profileID string,
	leader *model.MoveLeaderRequest,
) (
	*model.MoveLeaderResponse,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	conn, err := etcd.Store.Get(s.repo, profileID, s.store, ctx)
	if err != nil {
		logging.Error.Printf(" [APP]  TraceID-%s Failed to establish connection with the cluster. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return nil, err
	}
	moveLeaderResponse, err := etcd.MoveLeader(ctx, conn, leader)
	if err != nil {
		logging.Error.Printf(" [APP] TraceID-%s Failed transfer the  cluster leader. Error-%v ClusterProfileID-%s",
			requestTraceID, err, profileID)
		return nil, err
	} else {
		logging.Info.Printf(" [APP] TraceID-%s Successfully transferred the cluster leadership. ClusterProfileID-%s",
			requestTraceID, profileID)
		return moveLeaderResponse, nil
	}
}
