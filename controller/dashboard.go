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
	clusterProfileRepo model.ClusterProfileRepo,
	dashboardService model.DashboardService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		clusterProfileID := mux.Vars(r)["cluster_profile_id"]
		if clusterState, err := dashboardService.ViewCluster(clusterProfileID, clusterProfileRepo, r.Context()); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "cluster state", clusterState)
		}

	}

}
