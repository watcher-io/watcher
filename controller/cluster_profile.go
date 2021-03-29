package controller

import (
	"encoding/json"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/model"
	"github.com/watcher-io/watcher/response"
	"github.com/watcher-io/watcher/validator"
	"net/http"
)

type clusterProfileController struct {
	svc model.ClusterProfileService
}

func NewClusterProfileController(
	svc model.ClusterProfileService,
) *clusterProfileController {
	return &clusterProfileController{svc}
}

func (c *clusterProfileController) Fetch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if profiles, err := c.svc.FetchAll(r.Context()); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "cluster profiles retrieved", profiles)
		}
	}
}

func (c *clusterProfileController) Create() http.HandlerFunc {
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
		if profile, err := c.svc.Create(r.Context(), &clusterCreateRequest); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "cluster profile created", profile)
		}
	}
}

func (c *clusterProfileController) UploadCertificate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if certIDs, err := c.svc.UploadCertificate(r.Context(), r); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "certificate ids", certIDs)
		}
	}
}
