package service

import (
	"context"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/model"
	"github.com/watcher-io/watcher/utility"
)

type authService struct {
	repo model.UserRepo
}

func NewAuthService(repo model.UserRepo) *authService {
	return &authService{repo}
}

func (s *authService) Login(
	ctx context.Context,
	object *model.LoginRequest,
) (
	*model.LoginResponse,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	user, err := s.repo.Fetch(ctx, object.UserName)
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
