package model

import (
	"context"
	"net/http"
)

type Cluster struct {
	Members []Member `json:"members"`
	Leader  uint64   `json:"leader"`
	ID      uint64   `json:"id"`
}

type Member struct {
	ID         uint64   `json:"id"`
	Name       string   `json:"name"`
	DbSize     int64    `json:"db_size"`
	Version    string   `json:"version"`
	Alarms     []int32  `json:"alarms"`
	ClientURLS []string `json:"client_urls"`
	PeerURLs   []string `json:"peer_urls"`
	RaftIndex  uint64   `json:"raft_index"`
	RaftTerm   uint64   `json:"raft_term"`
	IsLearner  bool     `json:"is_learner"`
}

type DashboardService interface {
	ViewCluster(context.Context, string) (*Cluster, error)
}

type DashboardController interface {
	View() http.HandlerFunc
}
