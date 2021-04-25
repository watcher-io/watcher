package model

import (
	"context"
	"net/http"
)

type Header struct {
	ClusterID uint64 `json:"cluster_id"`
	MemberID  uint64 `json:"member_id"`
	Revision  int64  `json:"revision"`
	RaftTerm  uint64 `json:"raft_term"`
}

type PutKVRequest struct {
	Key   string `json:"key"    validate:"required"`
	Value string `json:"value"  validate:"required"`
}

type KV struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Version int64  `json:"version"`
	// create_revision is the revision of last creation on this key.
	CreateRevision int64 `json:"create_revision"`
	// mod_revision is the revision of last modification on this key.
	ModRevision int64 `json:"mod_revision"`
	Lease       int64 `json:"lease"`
}

type PutKVResponse struct {
	Header     Header `json:"header"`
	NewKV      bool   `json:"new_kv"`
	PreviousKV KV     `json:"previous_kv"`
}

type GetKVRequest struct {
	// Key is the user requested key or the prefix of the key
	Key string `json:"key"                validate:"required"`
	// WithPrefix enables key search with prefix
	Prefix bool `json:"prefix"`
	// Limit specifies number of keys to be returned when WithPrefix is enabled
	Limit    int64 `json:"limit"`
	Revision int64 `json:"revision"`
	// Get keys that are greater than or equal to the given key using byte compare
	FromKey bool `json:"from_key"`
	// Get only the keys
	KeysOnly bool `json:"keys_only"`
	// Get will return the keys in the range [key, end)
	Range string `json:"range"`
	// Returns on the count of the keys
	CountOnly bool `json:"count_only"`
}

type GetKVResponse struct {
	Header     Header `json:"header"`
	// more indicates if there are more keys to return in the requested range.
	More bool `json:"more"`
	// count is set to the number of keys within the range when requested.
	Count     int64 `json:"count"`
	KeyValues []KV  `json:"key_values"`
}

type DeleteKVRequest struct {
	Key string `json:"key"                validate:"required"`
	// WithPrefix enables key delete with prefix
	Prefix bool `json:"prefix"`
	// Get keys that are greater than or equal to the given key using byte compare
	FromKey bool `json:"from_key"`
	// Delete will remove the keys in the range [key, end)
	Range string `json:"range"`
}

type DeleteKVResponse struct {
	Header     Header `json:"header"`
	// count is set to the number of keys within the range when requested.
	Count int64 `json:"count"`
	//KeyValues []KV  `json:"key_values"`
}

type KVService interface {
	Put(context.Context, string, *PutKVRequest) (*PutKVResponse, error)
	Get(context.Context, string, *GetKVRequest) (*GetKVResponse, error)
	Delete(context.Context, string, *DeleteKVRequest) (*DeleteKVResponse, error)
}

type KVController interface {
	Put() http.HandlerFunc
	Get() http.HandlerFunc
	Delete() http.HandlerFunc
}
