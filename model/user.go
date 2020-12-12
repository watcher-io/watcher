package model

import (
	"context"
	"net/http"
	"os"
)

type User struct {
	ID        string `json:"id"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"     validate:"required"`
	FirstName string `json:"first_name"   validate:"required"`
	LastName  string `json:"last_name"    validate:"required"`
}

type UserRepo interface {
	Create(*User, context.Context) error
	Fetch(string, context.Context) (*User, error)
}

type UserService interface {
	Create(*User, UserRepo, context.Context) (*User, error)
	Fetch(string, UserRepo, context.Context) (*User, error)
	Exists(string, UserRepo, context.Context) (bool, error)
}

type UserController interface {
	Create(UserRepo, UserService) http.HandlerFunc
	Exists(UserRepo, UserService) http.HandlerFunc
}

func (*User) Prefix() string {
	return os.Getenv("USER_PREFIX")
}
