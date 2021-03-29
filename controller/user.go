package controller

import (
	"encoding/json"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/model"
	"github.com/watcher-io/watcher/response"
	"github.com/watcher-io/watcher/validator"
	"net/http"
)

type userController struct {
	svc model.UserService
}

func NewUserController(
	svc model.UserService,
) *userController {
	return &userController{svc}
}

func (c *userController) Create() http.HandlerFunc {
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
		if user, err := c.svc.Create(r.Context(), &createRequest); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "user created successfully", user)
		}
	}
}

func (c *userController) Exists() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if status, err := c.svc.Exists(r.Context(), "admin"); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "admin existence status", status)
		}
	}
}
