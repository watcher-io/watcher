package controller

import (
	"encoding/json"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/response"
	"github.com/aka-achu/watcher/validator"
	"net/http"
)

type authController struct{}

func NewAuthController() *authController {
	return &authController{}
}

func (*authController) Login(
	repo model.UserRepo,
	svc model.AuthService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestTraceID := r.Context().Value("trace_id").(string)
		var loginRequest model.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to decode the request body. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if err := validator.Validate.Struct(loginRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to validate request body for required fields. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if resp, err := svc.Login(&loginRequest, repo, r.Context()); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "login successful", resp)
		}
	}
}
