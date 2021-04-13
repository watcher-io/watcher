package response

import (
	"encoding/json"
	"github.com/watcher-io/watcher/model"
	"net/http"
)

func BadRequest(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(model.Format{
		ResponseMessage: msg,
	})
}

func InternalServerError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(model.Format{
		ResponseMessage: msg,
	})
}

func Success(w http.ResponseWriter, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.Format{
		ResponseMessage: msg,
		Data:            data,
	})
}

func Conflict(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	_ = json.NewEncoder(w).Encode(model.Format{
		ResponseMessage: msg,
	})
}

func UnAuthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(model.Format{
		ResponseMessage: msg,
	})
}
