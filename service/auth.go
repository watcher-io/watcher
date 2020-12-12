package service

import (
	"context"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/pkg_error"
	"github.com/aka-achu/watcher/utility"
)


type authService struct {}

func NewAuthService() *authService {
	return &authService{}
}

func (*authService) Login(
	object *model.LoginRequest,
	r model.UserRepo,
	ctx context.Context,
) (*model.LoginResponse, error) {
	requestTraceID := ctx.Value("trace_id").(string)
	user, err := r.Fetch(object.UserName, ctx)
	if err != nil {
		logging.Error.Printf(" [DB] Failed to fetch user profile details. Error-%v TraceID-%s",
			err, requestTraceID)
		return nil, pkg_error.ErrFailedToFetchUser
	}
	if user.Password == utility.Hash(object.Password) {
		if token, err := utility.CreateToken(user.UserName); err != nil {
			logging.Error.Printf(" [APP] Valid user credential but failed to generate access token. Error-%v TraceID-%s",
				err, requestTraceID)
			return nil, err
		} else {
			logging.Info.Printf(" [APP] Valid user credential. Successfully generated access token. TraceID-%s",
				requestTraceID)
			return &model.LoginResponse{
				AccessToken: token,
				UserName:    user.UserName,
				FirstName:   user.FirstName,
				LastName:    user.LastName,
			}, nil
		}
	} else {
		logging.Warn.Printf(" [APP] Invalid user credential. TraceID-%s", requestTraceID)
		return nil, pkg_error.ErrInvalidCredential
	}
}
