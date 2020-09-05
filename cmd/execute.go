package cmd

import (
	"github.com/aka-achu/watcher/controller"
	"github.com/aka-achu/watcher/logging"
	"net/http"
	"os"
	"path/filepath"
)

// Execute starts the web server on the specified server address
// in the ".env" file.
func Execute() {
	// Initializing the controller, registering the endpoints
	router := controller.Initialize()

	if os.Getenv("BUILD") == "Prod" {
		logging.Info.Printf(" [APP] Starting server @%s", os.Getenv("SERVER_ADDRESS"))
		logging.Error.Fatal(http.ListenAndServeTLS(
			os.Getenv("SERVER_ADDRESS"),
			filepath.Join("cert",os.Getenv("TLS_CERTIFICATE_FILE")),
			filepath.Join("cert",os.Getenv("TLS_KEY_FILE")),
			router,
		))
	} else {
		logging.Info.Printf(" [APP] Starting server @%s", os.Getenv("SERVER_ADDRESS"))
		logging.Error.Fatal(http.ListenAndServe(
			os.Getenv("SERVER_ADDRESS"),
			router,
		))
	}
}
