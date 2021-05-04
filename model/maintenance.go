package model

import (
	"context"
	"net/http"
)

type ListAlarmResponse struct {
	Header Header           `json:"header"`
	Alarms map[uint64]int32 `json:"alarms"`
}

type DisarmAlarmRequest struct {
	MemberID uint64 `json:"member_id" validate:"required"`
	Alarm    int32  `json:"alarm" validate:"required"`
}

type DisarmAlarmResponse ListAlarmResponse

type DefragmentRequest struct {
	Endpoint string `json:"endpoint" validate:"required"`
}

type CompactRequest struct {
	Revision int64 `json:"revision" validate:"required"`
}

type CompactResponse struct {
	Header Header `json:"header"`
}

type MoveLeaderRequest struct {
	TransfereeID    uint64   `json:"transferee_id"`
	LeaderEndPoints []string `json:"leader_end_points"`
}

type MoveLeaderResponse struct {
	Header Header `json:"header"`
}

type MaintenanceService interface {
	DisarmAlarm(context.Context, string, *DisarmAlarmRequest) (*DisarmAlarmResponse, error)
	MoveLeader(context.Context, string, *MoveLeaderRequest) (*MoveLeaderResponse, error)
	Defragment(context.Context, string, *DefragmentRequest) error
	Snapshot(context.Context, string) (string, error)
}

type MaintenanceController interface {
	DisarmAlarm() http.HandlerFunc
	Defragment() http.HandlerFunc
	Snapshot() http.HandlerFunc
	MoveLeader() http.HandlerFunc
}
