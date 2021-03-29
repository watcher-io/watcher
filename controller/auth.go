package controller

import (
	"encoding/json"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/model"
	"github.com/watcher-io/watcher/response"
	"github.com/watcher-io/watcher/validator"
	"net/http"
)

type authController struct {
	svc model.AuthService
}

func NewAuthController(
	svc model.AuthService,
) *authController {
	return &authController{svc}
}

func (c *authController) Login() http.HandlerFunc {
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
		if resp, err := c.svc.Login(r.Context(), &loginRequest); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "login successful", resp)
		}
	}
}
