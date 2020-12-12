package service

import (
	"context"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/utility"
)

type authService struct{}

func NewAuthService() *authService {
	return &authService{}
}

func (*authService) Login(
	object *model.LoginRequest,
	repo model.UserRepo,
	ctx context.Context,
) (
	*model.LoginResponse,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	user, err := repo.Fetch(object.UserName, ctx)
	if err != nil {
		logging.Error.Printf(" [DB] TraceID-%s Failed to fetch user profile details. Error-%v",
			requestTraceID, err)
		return nil, err
	}
	if user.Password == utility.Hash(object.Password) {
		if token, err := utility.CreateToken(user.UserName); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to generate access token. Error-%v",
				requestTraceID, err)
			return nil, err
		} else {
			return &model.LoginResponse{
				AccessToken: token,
				UserName:    user.UserName,
				FirstName:   user.FirstName,
				LastName:    user.LastName,
			}, nil
		}
	} else {
		logging.Warn.Printf(" [APP] TraceID-%s Invalid user credential.", requestTraceID)
		return nil, err
	}
}
