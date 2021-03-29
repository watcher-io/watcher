package service

import (
	"context"
	"errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/model"
	"github.com/watcher-io/watcher/utility"
)

type userService struct {
	repo model.UserRepo
}

func NewUserService(
	repo model.UserRepo,
) *userService {
	return &userService{repo}
}
func (s *userService) Create(
	ctx context.Context,
	user *model.User,
) (
	*model.User,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	user.UserName = "admin"
	if status, _ := s.Exists(ctx, user.UserName); !status {
		user.ID = uuid.New().String()
		user.Password = utility.Hash(user.Password)
		if err := s.repo.Create(ctx, user); err != nil {
			logging.Error.Printf(" [DB] TraceID-%s Failed to create the user Error-%v",
				requestTraceID, err)
			return nil, err
		} else {
			logging.Info.Printf(" [DB] TraceID-%s User profile created.",
				requestTraceID)
			user.Password = ""
			return user, nil
		}
	} else {
		return nil, errors.New("user already exists")
	}
}

func (s *userService) Fetch(
	ctx context.Context,
	userName string,
) (
	*model.User,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	if user, err := s.repo.Fetch(ctx, userName); err != nil {
		logging.Error.Printf(" [DB] TraceID-%s Failed to fetch user profile details. Error-%v",
			requestTraceID, err)
		return nil, err
	} else {
		user.Password = ""
		return user, nil
	}
}

func (s *userService) Exists(
	ctx context.Context,
	userName string,
) (
	bool,
	error,
) {
	_, err := s.repo.Fetch(ctx, userName)
	if errors.Is(err, badger.ErrKeyNotFound) {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}
