package model

type SaveAdminProfileRequest struct {
	Password string `json:"password" validate:"required"`
}
