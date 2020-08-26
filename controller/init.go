package controller

import (
	"github.com/aka-achu/watcher/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Initialize, initialized a router, created sub routers, integrates middlewares,
// registers the handle functions
func Initialize() *mux.Router {

	// Creating a new router
	router := mux.NewRouter()
	// Adding a middleware for handling cors
	router.Use(cors.AllowAll().Handler)

	// Declaring Controllers in order to access the handle functions
	var authController AuthController
	var clusterController ClusterController
	var statusController StatusController

	// Creating sub routers with different path prefix
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	clusterRouter := router.PathPrefix("/api/cluster").Subrouter()
	statusWithAuthRouter := router.PathPrefix("/api/status").Subrouter()
	statusWithOutAuthRouter := router.PathPrefix("/api/status").Subrouter()


	// Integrating the Auth and NoAuth middlewares
	authRouter.Use(middleware.NoAuthLogging)
	clusterRouter.Use(middleware.AuthLogging)
	statusWithAuthRouter.Use(middleware.AuthLogging)
	statusWithOutAuthRouter.Use(middleware.AuthLogging)

	// Registering the handle function for different request paths
	authRouter.HandleFunc("/checkAdminStatus", authController.CheckAdminInitStatus)
	authRouter.HandleFunc("/saveAdminProfile", authController.SaveAdminProfile)
	authRouter.HandleFunc("/login", authController.Login)

	clusterRouter.HandleFunc("/fetch", clusterController.FetchClusterProfiles)
	clusterRouter.HandleFunc("/create", clusterController.CreateClusterProfile)

	statusWithAuthRouter.HandleFunc("/backup", statusController.TakeBackup)
	statusWithOutAuthRouter.HandleFunc("/useSnapshot", statusController.ReInitDBWithSnapshot)

	return router
}
