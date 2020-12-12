package service

import (
	"context"
	"errors"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/pkg_error"
	"github.com/aka-achu/watcher/utility"
	"github.com/aka-achu/watcher/validator"
	"github.com/dgraph-io/badger/v2"
	"github.com/google/uuid"
)

type userService struct{}

func NewUserService() *userService {
	return &userService{}
}
func (s *userService) Create(
	user *model.User,
	r model.UserRepo,
	ctx context.Context,
) (
	*model.User,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	if err := validator.Validate.Struct(user); err != nil {
		logging.Error.Printf(" [SVC] Failed to validate the request object. Error-%v TraceID-%s",
			err, requestTraceID)
		return nil, pkg_error.ErrMissingRequiredFields
	}
	user.UserName = "admin"
	if status, _ := s.Exists(user.UserName, r, ctx); !status {
		user.ID = uuid.New().String()
		user.Password = utility.Hash(user.Password)
		if err := r.Create(user, ctx); err != nil {
			logging.Error.Printf(" [DB] Failed to create the user Error-%v TraceID-%s",
				err, requestTraceID)
			return nil, pkg_error.ErrFailedToCreateUser
		} else {
			logging.Info.Printf(" [DB] User profile created. Error-%v TraceID-%s",
				err, requestTraceID)
			user.Password = ""
			return user, nil
		}
	} else {
		logging.Info.Printf(" [DB] User already exists with same user_name. TraceID-%s",
			requestTraceID)
		return nil, pkg_error.ErrUserAlreadyExists
	}
}

func (*userService) Fetch(
	userName string,
	r model.UserRepo,
	ctx context.Context,
) (
	*model.User,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	if user, err := r.Fetch(userName, ctx); err != nil {
		logging.Error.Printf(" [DB] Failed to fetch user profile details. Error-%v TraceID-%s",
			err, requestTraceID)
		return nil, pkg_error.ErrFailedToFetchUser
	} else {
		logging.Info.Printf(" [DB] User profile fetched. TraceID-%s",
			requestTraceID)
		user.Password = ""
		return user, nil
	}
}

func (*userService) Exists(
	userName string,
	r model.UserRepo,
	ctx context.Context,
) (
	bool,
	error,
) {
	_, err := r.Fetch(userName, ctx)
	if errors.Is(err, badger.ErrKeyNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
