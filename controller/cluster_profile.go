package controller

import (
	"encoding/json"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/response"
	"github.com/aka-achu/watcher/validator"
	"net/http"
)

type clusterProfileController struct{}

func NewClusterProfileController() *clusterProfileController {
	return &clusterProfileController{}
}

func (*clusterProfileController) Fetch(
	repo model.ClusterProfileRepo,
	svc model.ClusterProfileService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if profiles, err := svc.FetchAll(repo, r.Context()); err != nil {
			response.InternalServerError(w,err.Error())
		} else {
			response.Success(w,"cluster profiles retrieved", profiles)
		}
	}
}

func (*clusterProfileController) Create(
	repo model.ClusterProfileRepo,
	svc model.ClusterProfileService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestTraceID := r.Context().Value("trace_id").(string)
		var clusterCreateRequest model.ClusterProfile
		if err := json.NewDecoder(r.Body).Decode(&clusterCreateRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to decode the request body. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if err := validator.Validate.Struct(clusterCreateRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to validate request body for required fields. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if profile, err := svc.Create(&clusterCreateRequest, repo, r.Context()); err != nil {
			response.InternalServerError(w,err.Error())
		} else {
			response.Success(w,"cluster profile created", profile)
		}
	}
}
