package service

import (
	"context"
	"errors"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/utility"
	"github.com/dgraph-io/badger/v2"
	"github.com/google/uuid"
)

type userService struct{}

func NewUserService() *userService {
	return &userService{}
}
func (svc *userService) Create(
	user *model.User,
	repo model.UserRepo,
	ctx context.Context,
) (
	*model.User,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	user.UserName = "admin"
	if status, _ := svc.Exists(user.UserName, repo, ctx); !status {
		user.ID = uuid.New().String()
		user.Password = utility.Hash(user.Password)
		if err := repo.Create(user, ctx); err != nil {
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

func (*userService) Fetch(
	userName string,
	repo model.UserRepo,
	ctx context.Context,
) (
	*model.User,
	error,
) {
	requestTraceID := ctx.Value("trace_id").(string)
	if user, err := repo.Fetch(userName, ctx); err != nil {
		logging.Error.Printf(" [DB] TraceID-%s Failed to fetch user profile details. Error-%v",
			requestTraceID, err)
		return nil, err
	} else {
		user.Password = ""
		return user, nil
	}
}

func (*userService) Exists(
	userName string,
	repo model.UserRepo,
	ctx context.Context,
) (
	bool,
	error,
) {
	_, err := repo.Fetch(userName, ctx)
	if errors.Is(err, badger.ErrKeyNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
