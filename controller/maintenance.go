package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/model"
	"github.com/watcher-io/watcher/response"
	"github.com/watcher-io/watcher/validator"
	"net/http"
)

type maintenanceController struct {
	svc model.MaintenanceService
}

func NewMaintenanceController(
	svc model.MaintenanceService,
) *maintenanceController {
	return &maintenanceController{svc}
}

func (c *maintenanceController) ListAlarm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clusterProfileID := mux.Vars(r)["cluster_profile_id"]
		if listAlarmResponse, err := c.svc.ListAlarm(r.Context(), clusterProfileID); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "cluster alarms fetched successfully", listAlarmResponse)
		}
	}
}

func (c *maintenanceController) DisarmAlarm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestTraceID := r.Context().Value("trace_id").(string)
		clusterProfileID := mux.Vars(r)["cluster_profile_id"]
		var disarmAlarmRequest model.DisarmAlarmRequest
		if err := json.NewDecoder(r.Body).Decode(&disarmAlarmRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to decode the request body. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if err := validator.Validate.Struct(disarmAlarmRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to validate request body for required fields. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if disarmAlarmResponse, err := c.svc.DisarmAlarm(r.Context(), clusterProfileID, &disarmAlarmRequest); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "disarmed the requested alarm successfully", disarmAlarmResponse)
		}
	}
}

func (c *maintenanceController) Defragment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestTraceID := r.Context().Value("trace_id").(string)
		clusterProfileID := mux.Vars(r)["cluster_profile_id"]
		var defragmentRequest model.DefragmentRequest
		if err := json.NewDecoder(r.Body).Decode(&defragmentRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to decode the request body. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if err := validator.Validate.Struct(defragmentRequest); err != nil {
			logging.Error.Printf(" [APP] TraceID-%s Failed to validate request body for required fields. Error-%v",
				requestTraceID, err)
			response.BadRequest(w, err.Error())
			return
		}
		if defragmentResponse, err := c.svc.Defragment(r.Context(), clusterProfileID, &defragmentRequest); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			response.Success(w, "de-fragmented the node successfully", defragmentResponse)
		}
	}
}

func (c *maintenanceController) Snapshot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clusterProfileID := mux.Vars(r)["cluster_profile_id"]
		if fileName, err := c.svc.Snapshot(r.Context(), clusterProfileID); err != nil {
			response.InternalServerError(w, err.Error())
		} else {
			http.ServeFile(w, r, fileName)
		}
	}
}
