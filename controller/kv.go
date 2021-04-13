package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/model"
	"github.com/watcher-io/watcher/response"
	"github.com/watcher-io/watcher/validator"
	"net/http"
)

type kvController struct {
	svc model.KVService
}

func NewKVController(
	svc model.KVService,
) *kvController {
	return &kvController{svc}
}

func (c *kvController) Put() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestTraceID := r.Context().Value("trace_id").(string)
		clusterProfileID := mux.Vars(r)["cluster_profile_id"]
		var putKVRequest model.PutKVRequest
		if err := json.NewDecoder(r.Body).Decode(&putKVRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to decode the request body. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if err := validator.Validate.Struct(putKVRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to validate request body for required fields. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if putResponse, err := c.svc.Put(r.Context(), clusterProfileID, &putKVRequest); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "key-value stored successfully", putResponse)
		}
	}
}

func (c *kvController) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestTraceID := r.Context().Value("trace_id").(string)
		clusterProfileID := mux.Vars(r)["cluster_profile_id"]
		var getKVRequest model.GetKVRequest
		if err := json.NewDecoder(r.Body).Decode(&getKVRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to decode the request body. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if err := validator.Validate.Struct(getKVRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to validate request body for required fields. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if getResponse, err := c.svc.Get(r.Context(), clusterProfileID, &getKVRequest); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "key-value stored successfully", getResponse)
		}
	}
}