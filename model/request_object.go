package model

// SaveAdminProfileRequest is a object containing fields required to create/initialize the admin profile
type SaveAdminProfileRequest struct {
	Password string `json:"password" validate:"required"`
}

// LoginRequest is a object containing fields required to log into the application
type LoginRequest struct {
	Password string `json:"password" validate:"required"`
}
