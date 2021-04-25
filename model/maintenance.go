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

type DefragmentResponse struct {
	Header Header `json:"header"`
}

type MaintenanceService interface {
	ListAlarm(context.Context, string) (*ListAlarmResponse, error)
	DisarmAlarm(context.Context, string, *DisarmAlarmRequest) (*DisarmAlarmResponse, error)
	Defragment(context.Context, string, *DefragmentRequest) (*DefragmentResponse, error)
	Snapshot(context.Context, string) (string, error)
}

type MaintenanceController interface {
	ListAlarm() http.HandlerFunc
	DisarmAlarm() http.HandlerFunc
	Defragment() http.HandlerFunc
	Snapshot() http.HandlerFunc
}
