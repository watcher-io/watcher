package controller

import (
	"github.com/aka-achu/watcher/etcd"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/response"
	"github.com/gorilla/mux"
	"net/http"
)

// DashboardController is an empty struct. All dashboard related handle functions will be implemented
// on this struct. This is used as a logical partition for all the handle functions
// in controller package.
type DashboardController struct{}

func (DashboardController) fetch(w http.ResponseWriter, r *http.Request) {
	// Getting the request tracing id from the request context
	requestTraceID := r.Context().Value("trace_id").(string)

	// fetching the cluster_profile_id from the request URI
	clusterProfileID := mux.Vars(r)["cluster_profile_id"]

	// Getting the cluster connection for the requested cluster profile
	conn, err := etcd.ClusterConnection.Get(clusterProfileID)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to establish connection with the cluster. Error-%v TraceID-%s ClusterProfileID-%s", err, requestTraceID, clusterProfileID)
		response.InternalServerError(w, "4000", err.Error())
		return
	}

	// Fetching the state of the cluster and members
	cluster, err := etcd.FetchMember(conn)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to fetch cluster info. Error-%v TraceID-%s ClusterProfileID-%s", err, requestTraceID, clusterProfileID)
		response.InternalServerError(w, "4001", err.Error())
		return
	}
	logging.Info.Printf(" [APP] Successfully fetched cluster info. TraceID-%s ClusterProfileID-%s", requestTraceID, clusterProfileID)
	response.Success(w,"4002","Cluster state", cluster)
}