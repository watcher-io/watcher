package cmd

import (
	"context"
	"crypto/tls"
	"github.com/gorilla/mux"
	"github.com/watcher-io/watcher/controller"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/middleware"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Execute() {
	router := controller.Initialize()
	//router.Use(cors.AllowAll().Handler)
	router.Use(middleware.CORS)
	server := getServer(os.Getenv("BUILD") == "Prod", router)

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logging.Info.Printf(" [APP] Server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logging.Error.Fatalf(" [ERROR] Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	//if os.Getenv("BUILD") == "Prod" {
	//	logging.Info.Printf(" [APP] Starting server @%s", os.Getenv("SERVER_ADDRESS"))
	//	if err := server.ListenAndServeTLS(
	//		filepath.Join("cert", os.Getenv("TLS_CERTIFICATE_FILE")),
	//		filepath.Join("cert", os.Getenv("TLS_KEY_FILE")),
	//	); err != nil && err != http.ErrServerClosed {
	//		logging.Error.Fatal(err)
	//	}
	//} else {
	//}
	logging.Info.Printf(" [APP] Starting server @%s", os.Getenv("SERVER_ADDRESS"))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logging.Error.Fatal(err)
	}

	<-done
	logging.Info.Printf(" [APP] Server has stopped gracefully")
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
