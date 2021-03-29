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
	Create(context.Context, *User) error
	Fetch(context.Context, string) (*User, error)
}

type UserService interface {
	Create(context.Context, *User) (*User, error)
	Fetch(context.Context, string) (*User, error)
	Exists(context.Context, string) (bool, error)
}

type UserController interface {
	Create() http.HandlerFunc
	Exists() http.HandlerFunc
}

func (*User) Prefix() string {
	return os.Getenv("USER_PREFIX")
}
