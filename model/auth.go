package model

import (
	"context"
	"net/http"
)

type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"  validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	UserName    string `json:"user_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

type AuthService interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
}

type AuthController interface {
	Login() http.HandlerFunc
}
