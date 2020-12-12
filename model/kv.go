package model

import (
	"context"
	"net/http"
)

type PutKVRequest struct {
	Key   string `json:"key"    validate:"required"`
	Value string `json:"value"  validate:"required"`
}

type PutKVResponse struct {
	Revision             int64  `json:"revision"`
	MemberID             uint64 `json:"member_id"`
	RaftTerm             uint64 `json:"raft_term"`
	PreviousRevision     int64  `json:"previous_revision"`
	PreviousRevisionVale string `json:"previous_revision_vale"`
}

type GetKVRequest struct {
	// Key is the user requested key or the prefix of the key
	Key string `json:"key"                validate:"required"`
	// WithPrefix enables key search with prefix
	WithPrefix bool `json:"with_prefix"`
	// Limit specifies number of keys to be returned when WithPrefix is enabled
	Limit int64 `json:"limit"`
	// SortTarget is the sorting criteria - by key or by value
	SortTarget   string `json:"sort_target"        validate:"required"`
	SortOrder    string `json:"sort_order"         validate:"required"`
	Revision     int64  `json:"revision"`
	KeyIteration bool   `json:"key_iteration"`
}

type KVService interface {
	Put(string, *PutKVRequest, ClusterProfileRepo, context.Context) (*PutKVResponse, error)
}

type KVController interface {
	Put(ClusterProfileRepo, KVService) http.HandlerFunc
}
