package controller

import (
	"encoding/json"
	"errors"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/repo"
	"github.com/aka-achu/watcher/response"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// clusterController is an empty struct
// for all the cluster profile handle function to be implemented on
type clusterController struct{}

// NewClusterController return a initialized clusterProfileController object
func NewClusterController() *clusterController {
	return &clusterController{}
}

// fetchProfiles handle function returns all the cluster profiles present in
// the application.
func (*clusterController) fetchProfiles(db *repo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Getting the request tracing id from the request context
		requestTraceID := r.Context().Value("trace_id").(string)

		// Fetching the cluster profiles
		clusterProfiles, err := db.GetClusterProfiles()
		if err != nil {
			logging.Error.Printf(" [DB] Failed to fetch cluster info from the database. Error-%v TraceID-%s", err, requestTraceID)
			response.InternalServerError(w, "2002", err.Error())
		} else {
			logging.Info.Printf(" [DB] Successfully fetched cluster info from the database. TraceID-%s", requestTraceID)
			response.Success(w, "2003", "Successfully fetched cluster info", clusterProfiles)
		}
	}
}

// createProfile handle function creates a cluster profile give cluster details.
// After validating the required fields in the request body, ID and CreatedTime fields
// are populated and the cluster profile is created
func (*clusterController) createProfile(db *repo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Getting the request tracing id from the request context
		requestTraceID := r.Context().Value("trace_id").(string)

		// Decoding the request body to the model.ClusterProfile object
		var clusterCreateRequest model.ClusterProfile
		err := json.NewDecoder(r.Body).Decode(&clusterCreateRequest)
		if err != nil {
			logging.Error.Printf(" [APP] Failed to decode the request body. Error-%v TraceID-%s", err, requestTraceID)
			response.BadRequest(w, "2000", err.Error())
			return
		}

		// Validating the data present in the request body for cluster profile creation
		// Checking the existence of cluster name and the cluster endpoints
		if clusterCreateRequest.Name == "" || len(clusterCreateRequest.Endpoints) == 0 {
			logging.Error.Printf(" [APP] Failed to validate the request body. Error-%v TraceID-%s", errors.New("all the required fields are not present in the request body"), requestTraceID)
			response.BadRequest(w, "2001", "All the required fields are not present in the request body")
			return
		}

		// Populating the cluster id and creation time
		clusterCreateRequest.ID = uuid.New().String()
		clusterCreateRequest.CreationTime = time.Now().Unix()

		// Creating the cluster profile with the requested data
		err = db.CreateClusterProfile(&clusterCreateRequest)
		if err != nil {
			logging.Error.Printf(" [DB] Failed to create the cluster profile. Error-%v TraceID-%s", err, requestTraceID)
			response.InternalServerError(w, "2004", err.Error())
		} else {
			logging.Info.Printf(" [DB] Successfully created the cluster profile. TraceID-%s", requestTraceID)
			response.Success(w, "2005", "Successfully created the cluster profile", nil)
		}
	}
}
