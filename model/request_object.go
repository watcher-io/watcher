package model

// SaveAdminProfileRequest is a object containing fields required to create/initialize the admin profile
type SaveAdminProfileRequest struct {
	Password string `json:"password" validate:"required"`
}

// LoginRequest is a object containing fields required to log into the application
type LoginRequest struct {
	Password string `json:"password" validate:"required"`
}

// PutKVRequest is a object containing fields required to store key-value in etcd
type PutKVRequest struct {
	Key   string `json:"key"    validate:"required"`
	Value string `json:"value"  validate:"required"`
}

// GetKVRequest is a object containing fields required to fetch key-value in etcd
type GetKVWithPrefixRequest struct {
	// Key is the user requested key or the prefix of the key
	Key          string `json:"key"                validate:"required"`
	// WithPrefix enables key search with prefix
	WithPrefix   bool   `json:"with_prefix"`
	// Limit specifies number of keys to be returned when WithPrefix is enabled
	Limit        int64  `json:"limit"`
	// SortTarget is the sorting criteria - by key or by value
	SortTarget   string `json:"sort_target"        validate:"required"`
	SortOrder    string `json:"sort_order"         validate:"required"`
	Revision     int64  `json:"revision"`
	KeyIteration bool   `json:"key_iteration"`
}
