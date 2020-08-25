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

type AuthController struct {}

func (*AuthController) CheckAdminInitStatus( w http.ResponseWriter, r *http.Request) {
	user, err := repo.GetUserDetails()
	if err != nil {
		logging.Error.Printf(" [DB] Failed to fetch the admin details. %v", err)
		response.InternalServerError(w,"1002", err.Error())
		return
	}
	if !user.InitializationStatus {
		logging.Info.Printf(" [APP] The admin profile does not exist in the application.")
		response.Success(w,"1003","Admin profile isn't initialized", user)
		return
	} else {
		logging.Info.Printf(" [APP] The admin profile exists in the application.")
		user.Password = ""
		response.Success(w,"1004","Admin profile is initialized", user)
		return
	}

}

func (*AuthController) SaveAdminProfile( w http.ResponseWriter, r *http.Request) {
	var saveAdminProfileRequest model.SaveAdminProfileRequest
	err := json.NewDecoder(r.Body).Decode(&saveAdminProfileRequest)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to decode the request body. %v", err)
		response.BadRequest(w, "1000", err.Error())
		return
	}
	err = validator.Validate.Struct(saveAdminProfileRequest)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to validate the request body. %v", err)
		response.BadRequest(w, "1001", err.Error())
		return
	}
	user, err := repo.GetUserDetails()
	if err != nil {
		logging.Error.Printf(" [DB] Failed to fetch the admin details. %v", err)
		response.InternalServerError(w,"1002", err.Error())
		return
	}
	if user.Password == "" {
		user.Password = utility.Hash(saveAdminProfileRequest.Password)
		user.InitializationStatus = true
		if err := repo.SaveUserDetails(user); err != nil {
			logging.Error.Printf(" [DB] Failed to save admin details. %v", err)
			response.InternalServerError(w,"1005", err.Error())
		} else {
			logging.Info.Printf(" [DB] Successfully saved admin details.")
			response.Success(w,"1006","Successfully saved admin details", nil)
		}
		return
	} else {
		logging.Warn.Printf(" [APP] The admin profile already exists in the application. Attemp to re-create")
		response.Conflict(w,"1007","Attempt to re-initialize the admin profile")
		return
	}

}

func (*AuthController) Login( w http.ResponseWriter, r *http.Request) {
	var loginRequest model.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to decode the request body. %v", err)
		response.BadRequest(w, "1000", err.Error())
		return
	}
	err = validator.Validate.Struct(loginRequest)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to validate the request body. %v", err)
		response.BadRequest(w, "1001", err.Error())
		return
	}
	user, err := repo.GetUserDetails()
	if err != nil {
		logging.Error.Printf(" [DB] Failed to fetch the admin details. %v", err)
		response.InternalServerError(w,"1002", err.Error())
		return
	}
	if user.Password == "" {
		logging.Error.Printf(" [APP] Login attempt when the admin profile is not initialized.")
		response.InternalServerError(w,"1008", "Please initialize admin profile first.")
		return
	}
	if user.Password == utility.Hash(loginRequest.Password){
		if token, err := utility.CreateToken("admin"); err != nil {
			logging.Error.Printf(" [APP] Valid user credential but failed to generate access token. %v", err)
			response.InternalServerError(w,"1009", err.Error())
		} else {
			logging.Info.Printf(" [APP] Valid user credential. Successfully generated access token.")
			response.Success(w,"1010","Successful login", model.LoginResponse{Token: token})
		}
		return
	} else {
		logging.Warn.Printf(" [APP] Invalid user credential")
		response.UnAuthorized(w,"1011","Invalid user credential")
		return
	}
}