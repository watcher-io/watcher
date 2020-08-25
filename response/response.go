package response

import (
	"encoding/json"
	"net/http"
)

type Format struct {
	ResponseCode    string      `json:"response_code"`
	ResponseMessage string      `json:"response_message"`
	Data            interface{} `json:"data"`
}

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
