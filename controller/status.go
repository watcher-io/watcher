package controller

import (
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/repo"
	"github.com/aka-achu/watcher/response"
	"net/http"
)

// StatusController is an empty struct. All application status related handle functions will be implemented
// on this struct. This is used as a logical partition for all the handle functions
// in controller package.
type StatusController struct {}

// TakeBackup handle function takes a db snapshot and sends as client response
func (*StatusController) TakeBackup(w http.ResponseWriter, r *http.Request) {

	// Getting the request tracing id from the request context
	requestTraceID := r.Context().Value("trace_id").(string)

	// Setting the response headers for backup file download
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="watcher.db"`)

	// Taking the db backup
	if err:= repo.BackupRepo(w); err != nil {
		logging.Error.Printf(" [DB] Failed to take backup of the embedded database. Error-%v TraceID-%s", err, requestTraceID )
		response.InternalServerError(w,"3000", err.Error())
	}
}
