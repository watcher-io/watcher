package controller

import (
	"encoding/json"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/response"
	"github.com/gorilla/mux"
	"net/http"
)

type kvController struct{}

func NewKVController() *kvController {
	return &kvController{}
}

func (*kvController) Put(clusterProfileRepo model.ClusterProfileRepo, kvService model.KVService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		requestTraceID := r.Context().Value("trace_id").(string)
		clusterProfileID := mux.Vars(r)["cluster_profile_id"]

		var putKVRequest model.PutKVRequest
		if err := json.NewDecoder(r.Body).Decode(&putKVRequest); err != nil {
			logging.Error.Printf(" [APP] Failed to decode the request body. Error-%v TraceID-%s", err, requestTraceID)
			response.BadRequest(w, err.Error())
			return
		}
		if putResponse, err := kvService.Put(clusterProfileID, &putKVRequest, clusterProfileRepo, r.Context()); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "key-value stored successfully", putResponse)
		}
	}
}
