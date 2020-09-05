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
	KeyPrefix bool   `json:"key_prefix"`
	Key       string `json:"key"    validate:"required"`
	Value     string `json:"value"  validate:"required"`
}
