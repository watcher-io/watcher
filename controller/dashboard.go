package controller

import (
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/response"
	"github.com/gorilla/mux"
	"net/http"
)

type dashboardController struct{}

func NewDashboardController() *dashboardController {
	return &dashboardController{}
}

func (*dashboardController) View(
	repo model.ClusterProfileRepo,
	svc model.DashboardService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clusterProfileID := mux.Vars(r)["cluster_profile_id"]
		if clusterState, err := svc.ViewCluster(clusterProfileID, repo, r.Context()); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "cluster state", clusterState)
		}
	}
}
