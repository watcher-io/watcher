package controller

import (
	"encoding/json"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/response"
	"net/http"
)

type clusterProfileController struct{}

func NewClusterProfileController() *clusterProfileController {
	return &clusterProfileController{}
}

func (*clusterProfileController) Fetch(
	clusterProfileRepo model.ClusterProfileRepo,
	clusterProfileService model.ClusterProfileService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if profiles, err := clusterProfileService.FetchAll(clusterProfileRepo, r.Context()); err != nil {
			response.InternalServerError(w,err.Error())
		} else {
			response.Success(w,"cluster profiles retrieved", profiles)
		}
	}
}

func (*clusterProfileController) Create(
	clusterProfileRepo model.ClusterProfileRepo,
	clusterProfileService model.ClusterProfileService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestTraceID := r.Context().Value("trace_id").(string)
		var clusterCreateRequest model.ClusterProfile
		if err := json.NewDecoder(r.Body).Decode(&clusterCreateRequest); err != nil {
			logging.Error.Printf(" [APP] Failed to decode the request body. Error-%v TraceID-%s",
				err, requestTraceID)
			response.BadRequest(w, err.Error())
			return
		}

		if profile, err := clusterProfileService.Create(&clusterCreateRequest, clusterProfileRepo, r.Context()); err != nil {
			response.InternalServerError(w,err.Error())
		} else {
			response.Success(w,"cluster profile created", profile)
		}
	}
}
