package controller

import (
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/repo"
	"github.com/aka-achu/watcher/response"
	"github.com/aka-achu/watcher/state"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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
	if err:= repo.DB.BackupRepo(w); err != nil {
		logging.Error.Printf(" [DB] Failed to take backup of the embedded database. Error-%v TraceID-%s", err, requestTraceID )
		response.InternalServerError(w,"3000", err.Error())
	}
}

func (*StatusController) ReInitDBWithSnapshot(w http.ResponseWriter, r *http.Request) {

	// Getting the request tracing id from the request context
	requestTraceID := r.Context().Value("trace_id").(string)

	// Parsing the multipart form
	err := r.ParseMultipartForm(200000)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to parse multi part form. Error-%v TraceID-%s", err, requestTraceID)
		response.InternalServerError(w,"3001", err.Error())
		return
	}

	// Getting the file from the multipart form
	file := r.MultipartForm.File["db_snapshot"][0]
	// Validating the file extension of the uploaded snapshot
	if filepath.Ext(file.Filename) != ".db" {
		logging.Error.Printf(" [APP] Unexpected db snapshot file extension. TraceID-%s", requestTraceID)
		response.InternalServerError(w,"3002", "Unexpected db snapshot file extension")
		return
	}

	// Opening the uploaded file
	f, err := file.Open()
	if err != nil {
		logging.Error.Printf(" [APP] Failed to open the snapshot file present in the multi part from. Error-%v TraceID-%s", err, requestTraceID)
		response.InternalServerError(w,"3002", err.Error())
		return
	}

	// Reading the content of the uploaded file
	fileByte, err := ioutil.ReadAll(f)
	if err != nil {
		logging.Error.Printf(" [APP] Failed to read the content of the uploaded file. Error-%v TraceID-%s", err, requestTraceID)
		response.InternalServerError(w,"3003", err.Error())
		return
	}
	_ = f.Close()

	// Writing the content of the upload file to "data/watcher.db.snap"
	if err := ioutil.WriteFile(filepath.Join("data", "watcher.db.snap"), fileByte, 0666); err != nil {
		logging.Error.Printf(" [APP] Failed to write the content of the upload snapshot. Error-%v TraceID-%s", err, requestTraceID)
		response.InternalServerError(w,"3004", err.Error())
		return
	}

	// Validating the user uploaded db snapshot
	if err := repo.ValidateSnapshot(); err != nil {
		logging.Error.Printf(" [DB] Failed to validate the snapshot db. Error-%v TraceID-%s", err, requestTraceID)
		response.InternalServerError(w,"3005", err.Error())
		return
	}

	// Closing the existing database connection
	_ = repo.DB.Conn.Close()

	// Removing the existing database
	if err := os.Remove(filepath.Join("data","watcher.db")); err != nil {
		logging.Error.Printf(" [APP] Failed to remove the current database file. Error-%v TraceID-%s", err, requestTraceID)
		// re-initializing the current database as we failed to remove the db
		repo.Initialize()
		response.InternalServerError(w,"3006", err.Error())
		return
	}

	// Reaming the snapshot database to current database
	if err := os.Rename(filepath.Join("data","watcher.db.snap"), filepath.Join("data", "watcher.db")); err != nil {
		logging.Error.Printf(" [APP] Failed to rename the snapshot database. Error-%v TraceID-%s", err, requestTraceID)
		// As we failed to rename the snapshot database and we have delete the current database
		// we need to reinitialize the current database and all the buckets
		state.Validate()
		response.InternalServerError(w,"3007", err.Error())
		return
	}

	// If we have successfully deleted the previous database and renamed the snapshot database,
	// then we simply need to re-initialize the database connection.
	repo.Initialize()
	response.Success(w,"3008","Successfully replaced the default database with the snapshot", nil)

}
