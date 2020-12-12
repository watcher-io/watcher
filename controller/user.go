package controller

import (
	"encoding/json"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/response"
	"github.com/aka-achu/watcher/validator"
	"net/http"
)

type userController struct{}

func NewUserController() *userController {
	return &userController{}
}

func (*userController) Create(
	repo model.UserRepo,
	svc model.UserService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestTraceID := r.Context().Value("trace_id").(string)
		var createRequest model.User
		if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to decode the request body. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if err := validator.Validate.Struct(createRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to validate request body for required fields. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if user, err := svc.Create(&createRequest, repo, r.Context()); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "user created successfully", user)
		}
	}
}

func (*userController) Exists(
	repo model.UserRepo,
	svc model.UserService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if status, err := svc.Exists("admin", repo, r.Context()); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "admin existence status", status)
		}
	}
}
