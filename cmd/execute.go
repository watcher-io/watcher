package cmd

import (
	"crypto/tls"
	"github.com/aka-achu/watcher/controller"
	"github.com/aka-achu/watcher/logging"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"os"
	"path/filepath"
	"time"
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
	server := getServer(os.Getenv("BUILD") == "Prod", router)

	if os.Getenv("BUILD") == "Prod" {
		logging.Info.Printf(" [APP] Starting server @%s", os.Getenv("SERVER_ADDRESS"))
		logging.Error.Fatal(
			server.ListenAndServeTLS(
				filepath.Join("cert", os.Getenv("TLS_CERTIFICATE_FILE")),
				filepath.Join("cert", os.Getenv("TLS_KEY_FILE")),
			),
		)
	} else {
		logging.Info.Printf(" [APP] Starting server @%s", os.Getenv("SERVER_ADDRESS"))
		logging.Error.Fatal(
			server.ListenAndServe(),
		)
	}
}

func getServer(secure bool, router *mux.Router) *http.Server {
	if secure {
		return &http.Server{
			Addr:         os.Getenv("SERVER_ADDRESS"),
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
			TLSConfig: &tls.Config{
				PreferServerCipherSuites: true,
				CurvePreferences: []tls.CurveID{
					tls.CurveP256,
					tls.X25519, // Go 1.8 only
				},
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				},
			},
		}
	} else {
		return &http.Server{
			Addr:         os.Getenv("SERVER_ADDRESS"),
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		}
	}
}
