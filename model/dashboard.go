package model

import (
	"context"
	"net/http"
)

type Cluster struct {
	Members []ClusterMember `json:"members"`
	Leader  uint64          `json:"leader"`
	ID      uint64          `json:"id"`
}

type ClusterMember struct {
	ID         uint64       `json:"id"`
	Name       string       `json:"name"`
	PeerURLS   []string     `json:"peer_urls"`
	ClientURLS []string     `json:"client_urls"`
	IsLearner  bool         `json:"is_learner"`
	Status     MemberStatus `json:"status"`
}

type MemberStatus struct {
	Version          string `json:"version"`
	DbSize           int64  `json:"db_size"`
	DbSizeInUse      int64  `json:"db_size_in_use"`
	RaftIndex        uint64 `json:"raft_index"`
	RaftTerm         uint64 `json:"raft_term"`
	RaftAppliedIndex uint64 `json:"raft_applied_index"`
}

type DashboardService interface {
	ViewCluster(string, ClusterProfileRepo, context.Context) (*Cluster, error)
}

type DashboardController interface {
	View(ClusterProfileRepo, DashboardService) http.HandlerFunc
}