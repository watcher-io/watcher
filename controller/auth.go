package controller

import (
	"encoding/json"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/repo"
	"github.com/aka-achu/watcher/response"
	"github.com/aka-achu/watcher/utility"
	"github.com/aka-achu/watcher/validator"
	"net/http"
)

// AuthController is an empty struct. All authentication related handle functions will be implemented
// on this struct. This is used as a logical partition for all the handle functions
// in controller package.
type AuthController struct{}

// CheckAdminInitStatus returns the status of admin profile initialization.
func (*AuthController) CheckAdminInitStatus(w http.ResponseWriter, r *http.Request) {

	// Getting the request tracing id from the request context
	requestTraceID := r.Context().Value("trace_id").(string)

	// Fetching the admin details from the database
	user, err := repo.DB.GetUserDetails()
	if err != nil {
		logging.Error.Printf(" [DB] Failed to fetch the admin details. Error-%v TraceID-%s", err, requestTraceID)
		response.InternalServerError(w, "1002", err.Error())
		return
	}

	// If user.InitializationStatus is false then the admin profile has not been initialized
	if !user.InitializationStatus {
		logging.Info.Printf(" [APP] The admin profile does not exist in the application. TraceID-%s", requestTraceID)
		response.Success(w, "1003", "Admin profile isn't initialized", user)
	} else {
		logging.Info.Printf(" [APP] The admin profile exists in the application. TraceID-%s", requestTraceID)
		// Emptying the password field before sending the data to view layer
		user.Password = ""
		response.Success(w, "1004", "Admin profile is initialized", user)
	}
}

// SaveAdminProfile will be used to save admin profile details (only password in the current release)
// It will check the initialization status of admin profile, the admin password will be save only if
// the admin profile is not initialized before.
func (*AuthController) SaveAdminProfile(w http.ResponseWriter, r *http.Request) {

	// Getting the request tracing id from the request context
	requestTraceID := r.Context().Value("trace_id").(string)

	// Decoding the request body to the model.SaveAdminProfileRequest object
	var saveAdminProfileRequest model.SaveAdminProfileRequest
	err := json.NewDecoder(r.Body).Decode(&saveAdminProfileRequest)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to decode the request body. Error-%v TraceID-%s", err, requestTraceID)
		response.BadRequest(w, "1000", err.Error())
		return
	}

	// Validating the fields (password) present in the request body.
	err = validator.Validate.Struct(saveAdminProfileRequest)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to validate the request body. Error-%v TraceID-%s", err, requestTraceID)
		response.BadRequest(w, "1001", err.Error())
		return
	}

	// Fetching admin details from the database
	user, err := repo.DB.GetUserDetails()
	if err != nil {
		logging.Error.Printf(" [DB] Failed to fetch the admin details. Error-%v TraceID-%s", err, requestTraceID)
		response.InternalServerError(w, "1002", err.Error())
		return
	}

	// If the admin profile is not initialized before, user.Password is empty then save the password for admin
	if user.Password == "" {
		// Hashing the password
		user.Password = utility.Hash(saveAdminProfileRequest.Password)
		// Changing the initialization status to true
		user.InitializationStatus = true
		// Saving the admin profile details
		if err := repo.DB.SaveUserDetails(user); err != nil {
			logging.Error.Printf(" [DB] Failed to save admin details. Error-%v TraceID-%s", err, requestTraceID)
			response.InternalServerError(w, "1005", err.Error())
		} else {
			logging.Info.Printf(" [DB] Successfully saved admin details. TraceID-%s", requestTraceID)
			response.Success(w, "1006", "Successfully saved admin details", nil)
		}
	} else {
		logging.Warn.Printf(" [APP] The admin profile already exists in the application. Attemp to re-create. TraceID-%s", requestTraceID)
		response.Conflict(w, "1007", "Attempt to re-initialize the admin profile")
	}

}

// Login handle function will be used to log into the application.
// It validates the request body for existence of the password, if the hashed password matched
// the password store in the database, a JWT token will be created and sent to the user which will
// be used as an access_token for the application
func (*AuthController) Login(w http.ResponseWriter, r *http.Request) {

	// Getting the request tracing id from the request context
	requestTraceID := r.Context().Value("trace_id").(string)

	// Decoding the request body data to model.LoginRequest object
	var loginRequest model.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to decode the request body. Error-%v TraceID-%s", err, requestTraceID)
		response.BadRequest(w, "1000", err.Error())
		return
	}

	// Validating the request body for existence fot password
	err = validator.Validate.Struct(loginRequest)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to validate the request body. Error-%v TraceID-%s", err, requestTraceID)
		response.BadRequest(w, "1001", err.Error())
		return
	}

	// Fetching the user details from the database
	user, err := repo.DB.GetUserDetails()
	if err != nil {
		logging.Error.Printf(" [DB] Failed to fetch the admin details. Error-%v TraceID-%s", err, requestTraceID)
		response.InternalServerError(w, "1002", err.Error())
		return
	}

	// If the password field of the user data present in the database is empty then
	// the admin profile has not been initialized
	if user.Password == "" {
		logging.Error.Printf(" [APP] Login attempt when the admin profile is not initialized. TraceID-%s", requestTraceID)
		response.InternalServerError(w, "1008", "Please initialize admin profile first.")
		return
	}

	// Validating the hash of user password with the store password hash
	if user.Password == utility.Hash(loginRequest.Password) {
		// Generating JWT token which will be used as access_token
		if token, err := utility.CreateToken("admin"); err != nil {
			logging.Error.Printf(" [APP] Valid user credential but failed to generate access token. Error-%v TraceID-%s", err, requestTraceID)
			response.InternalServerError(w, "1009", err.Error())
		} else {
			logging.Info.Printf(" [APP] Valid user credential. Successfully generated access token. TraceID-%s", requestTraceID)
			response.Success(w, "1010", "Successful login", model.LoginResponse{Token: token})
		}
	} else {
		logging.Warn.Printf(" [APP] Invalid user credential. TraceID-%s", requestTraceID)
		response.UnAuthorized(w, "1011", "Invalid user credential")
	}
}
