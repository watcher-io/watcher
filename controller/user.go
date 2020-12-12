package controller

import (
	"encoding/json"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/pkg_error"
	"github.com/aka-achu/watcher/response"
	"net/http"
)

type userController struct{}

func NewUserController() *userController {
	return &userController{}
}

func (*userController) Create(ur model.UserRepo, us model.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestTraceID := r.Context().Value("trace_id").(string)
		var createRequest model.User
		if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
			logging.Error.Printf(" [APP] Failed to decode the request body. Error-%v TraceID-%s",
				err, requestTraceID)
			response.BadRequest(w, pkg_error.ErrFailedToDecodeRequestBody.Error())
			return
		}

		if user, err := us.Create(&createRequest, ur, r.Context()); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "user created successfully", user)
		}
	}
}

func (*userController) Exists(ur model.UserRepo, us model.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if status, err := us.Exists("admin",ur, r.Context()); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "admin existence status", status)
		}
	}
}