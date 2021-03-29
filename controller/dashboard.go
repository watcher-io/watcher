package controller

import (
	"github.com/gorilla/mux"
	"github.com/watcher-io/watcher/model"
	"github.com/watcher-io/watcher/response"
	"net/http"
)

type dashboardController struct {
	svc model.DashboardService
}

func NewDashboardController(
	svc model.DashboardService,
) *dashboardController {
	return &dashboardController{svc}
}

func (c *dashboardController) View() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clusterProfileID := mux.Vars(r)["cluster_profile_id"]
		if clusterState, err := c.svc.ViewCluster(r.Context(), clusterProfileID); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "cluster state", clusterState)
		}
	}
}
