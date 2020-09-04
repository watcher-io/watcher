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
	var clusterController ClusterProfileController
	var statusController StatusController
	var dashboardController DashboardController

	// Creating sub routers with different path prefix
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	clusterProfileRouter := router.PathPrefix("/api/clusterProfile").Subrouter()
	statusWithAuthRouter := router.PathPrefix("/api/status").Subrouter()
	statusWithOutAuthRouter := router.PathPrefix("/api/status").Subrouter()
	dashboardRouter := router.PathPrefix("/api/dashboard").Subrouter()


	// Integrating the Auth and NoAuth middlewares
	authRouter.Use(middleware.NoAuthLogging)
	clusterProfileRouter.Use(middleware.AuthLogging)
	statusWithAuthRouter.Use(middleware.AuthLogging)
	statusWithOutAuthRouter.Use(middleware.AuthLogging)
	dashboardRouter.Use(middleware.AuthLogging)

	// Registering the handle function for different request paths
	authRouter.HandleFunc("/checkAdminStatus", authController.CheckAdminInitStatus).Methods("GET")
	authRouter.HandleFunc("/saveAdminProfile", authController.SaveAdminProfile).Methods("POST")
	authRouter.HandleFunc("/login", authController.Login)

	clusterProfileRouter.HandleFunc("/fetch", clusterController.FetchClusterProfiles).Methods("GET")
	clusterProfileRouter.HandleFunc("/create", clusterController.CreateClusterProfile).Methods("POST")

	statusWithAuthRouter.HandleFunc("/backup", statusController.TakeBackup).Methods("GET")
	statusWithOutAuthRouter.HandleFunc("/useSnapshot", statusController.ReInitDBWithSnapshot).Methods("POST")

	dashboardRouter.HandleFunc("/fetch/{cluster_profile_id}", dashboardController.fetch).Methods("GET")

	return router
}
