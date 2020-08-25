package cmd

import (
	"github.com/aka-achu/watcher/controller"
	"github.com/aka-achu/watcher/logging"
	"net/http"
	"os"
	"path/filepath"
)

func Execute() {
	router := controller.Initialize()
	if os.Getenv("BUILD") == "Prod" {
		logging.Error.Fatal(http.ListenAndServeTLS(
			os.Getenv("SERVER_ADDRESS"),
			filepath.Join("config", "tls.cert"),
			filepath.Join("config", "tls.key"),
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
