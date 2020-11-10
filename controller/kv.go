package controller

import (
	"context"
	"encoding/json"
	"github.com/aka-achu/watcher/etcd"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/repo"
	"github.com/aka-achu/watcher/response"
	"github.com/aka-achu/watcher/validator"
	"github.com/gorilla/mux"
	"net/http"
)

// KVController is an empty struct
// for all the key-value handle function to be implemented on
type kvController struct{}

func NewKVController() *kvController {
	return &kvController{}
}

// put handle function stores a key value pair in the cluster
func (*kvController) put(db *repo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Getting the request tracing id from the request context
		requestTraceID := r.Context().Value("trace_id").(string)

		// fetching the cluster_profile_id from the request URI
		clusterProfileID := mux.Vars(r)["cluster_profile_id"]

		// Decoding the request body to the model.PutKVRequest object
		var putKVRequest model.PutKVRequest
		err := json.NewDecoder(r.Body).Decode(&putKVRequest)
		if err != nil {
			logging.Error.Printf(" [APP] Failed to decode the request body. Error-%v TraceID-%s", err, requestTraceID)
			response.BadRequest(w, "5000", err.Error())
			return
		}

		// Validating the fields (password) present in the request body.
		err = validator.Validate.Struct(putKVRequest)
		if err != nil {
			logging.Error.Printf(" [APP] Failed to validate the request body. Error-%v TraceID-%s", err, requestTraceID)
			response.BadRequest(w, "5001", err.Error())
			return
		}

		// Getting the cluster connection for the requested cluster profile
		conn, err := etcd.ClusterConnection.Get(db, clusterProfileID)
		if err != nil {
			logging.Error.Printf(" [APP] Failed to establish connection with the cluster. Error-%v TraceID-%s ClusterProfileID-%s", err, requestTraceID, clusterProfileID)
			response.InternalServerError(w, "5002", err.Error())
			return
		}

		// Storing the requested key-value pair in the cluster
		putResponse, err := etcd.PutKV(context.Background(), conn, &putKVRequest)
		if err != nil {
			logging.Error.Printf(" [APP] Failed to put the key-value pair in the cluster. Error-%v TraceID-%s ClusterProfileID-%s", err, requestTraceID, clusterProfileID)
			response.InternalServerError(w, "5003", err.Error())
			return
		}
		logging.Info.Printf(" [APP] Successfully stored the key-vlaue pair in the cluster. TraceID-%s ClusterProfileID-%s", requestTraceID, clusterProfileID)
		response.Success(w, "5002", "Successfully stored the kv pair", putResponse)
	}
}
