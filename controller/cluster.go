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

type ClusterController struct{}

func (*ClusterController) FetchClusterProfiles(w http.ResponseWriter, r *http.Request) {
	requestTraceID := r.Context().Value("trace_id").(string)
	clusterProfiles, err := repo.GetClusterProfiles()
	if err != nil {
		logging.Error.Printf(" [DB] Failed to fetch cluster info from the database. Error-%v TraceID-%s", err, requestTraceID)
		response.InternalServerError(w, "2002", err.Error())
	} else {
		logging.Info.Printf(" [DB] Successfully fetched cluster info from the database. TraceID-%s", requestTraceID)
		response.Success(w, "2003", "Successfully fetched cluster info", clusterProfiles)
	}
}

func (*ClusterController) CreateClusterProfile(w http.ResponseWriter, r *http.Request) {
	requestTraceID := r.Context().Value("trace_id").(string)
	var clusterCreateRequest model.Cluster
	err := json.NewDecoder(r.Body).Decode(&clusterCreateRequest)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to decode the request body. Error-%v TraceID-%s", err, requestTraceID)
		response.BadRequest(w, "2000", err.Error())
		return
	}
	if clusterCreateRequest.Name == "" || len(clusterCreateRequest.Endpoints) == 0 {
		logging.Error.Printf(" [APP] Failed to validate the request body. Error-%v TraceID-%s", errors.New("all the required fields are not present in the request body"), requestTraceID)
		response.BadRequest(w, "2001", err.Error())
		return
	}
	clusterCreateRequest.ID = uuid.New().String()
	clusterCreateRequest.CreationTime = time.Now().Unix()

	err = repo.CreateClusterProfile(&clusterCreateRequest)
	if err != nil {
		logging.Error.Printf(" [DB] Failed to create the cluster profile. Error-%v TraceID-%s", err, requestTraceID)
		response.InternalServerError(w,"2004", err.Error())
	} else {
		logging.Info.Printf(" [DB] Successfully created the cluster profile. TraceID-%s", requestTraceID)
		response.Success(w,"2005", "Successfully created the cluster profile", nil)
	}
}
