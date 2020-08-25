package response

import (
	"encoding/json"
	"net/http"
)

// Format, is a generic response object for all the handle functions
type Format struct {
	// ResponseCode, is a code which is mapped with a application event status
	// This is not same as http.Status code
	ResponseCode    string      `json:"response_code"`

	// ResponseMessage is application event message
	ResponseMessage string      `json:"response_message"`

	// Data is a generic interface which will be used to
	// send any type of data to the view layer
	Data            interface{} `json:"data"`
}

// getResponseBody, return a Format containing requested response code, message and data
func getResponseBody(code string, message string, data ...interface{}) Format {
	if len(data) == 0 {
		return Format{
			ResponseCode:    code,
			ResponseMessage: message,
		}
	} else {
		return Format{
			ResponseCode:    code,
			ResponseMessage: message,
			Data:            data[0],
		}
	}
}

func BadRequest(w http.ResponseWriter, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(getResponseBody(code, message))
}

func InternalServerError(w http.ResponseWriter, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(getResponseBody(code, message))
}

func Success(w http.ResponseWriter, code string, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(getResponseBody(code, message, data))
}

func Conflict(w http.ResponseWriter, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	_ = json.NewEncoder(w).Encode(getResponseBody(code, message))
}

func UnAuthorized(w http.ResponseWriter, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(getResponseBody(code, message))
}
