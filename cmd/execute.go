package cmd

import (
	"github.com/aka-achu/watcher/controller"
	"github.com/aka-achu/watcher/logging"
	"github.com/rs/cors"
	"net/http"
	"os"
	"path/filepath"
)

// Execute starts the web server on the specified server address
// in the ".env" file.
func Execute() {
	// Initializing the controller, registering the endpoints
	router := controller.Initialize()

	// Getting CORS_ORIGIN from env
	CorsOrigin := os.Getenv("CORS_ORIGIN")

	// Adding a middleware for handling cors
	if CorsOrigin != "" {
		router.Use(cors.New(cors.Options{
			AllowedOrigins:   []string{CorsOrigin},
			AllowCredentials: true,
			Debug:            os.Getenv("BUILD") != "Prod",
		},
		).Handler)
	} else {
		router.Use(cors.AllowAll().Handler)
	}

	if os.Getenv("BUILD") == "Prod" {
		logging.Info.Printf(" [APP] Starting server @%s", os.Getenv("SERVER_ADDRESS"))
		logging.Error.Fatal(http.ListenAndServeTLS(
			os.Getenv("SERVER_ADDRESS"),
			filepath.Join("cert", os.Getenv("TLS_CERTIFICATE_FILE")),
			filepath.Join("cert", os.Getenv("TLS_KEY_FILE")),
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
