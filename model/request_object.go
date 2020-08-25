package model

type SaveAdminProfileRequest struct {
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Password string `json:"password" validate:"required"`
}
