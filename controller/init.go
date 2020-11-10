package controller

import (
	"github.com/aka-achu/watcher/middleware"
	"github.com/gorilla/mux"
)

// Initialize, initialized a router, created sub routers, integrates middlewares,
// registers the handle functions
func Initialize() *mux.Router {

	// Creating a new router
	router := mux.NewRouter()

	// Initializing controller object
	// Creating a sub router
	// Adding no-auth or auth middleware based on the need
	// Registering the handle functions

	authController := NewAuthController()
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.Use(middleware.NoAuthLogging)
	authRouter.HandleFunc("/checkAdminStatus", authController.checkAdminInitStatus).Methods("GET")
	authRouter.HandleFunc("/saveAdminProfile", authController.saveAdminProfile).Methods("POST")
	authRouter.HandleFunc("/login", authController.login)

	clusterController := NewClusterController()
	clusterProfileRouter := router.PathPrefix("/api/clusterProfile").Subrouter()
	clusterProfileRouter.Use(middleware.AuthLogging)
	clusterProfileRouter.HandleFunc("/fetchProfiles", clusterController.fetchProfiles).Methods("GET")
	clusterProfileRouter.HandleFunc("/createProfile", clusterController.createProfile).Methods("POST")

	dashboardController := NewDashboardController()
	dashboardRouter := router.PathPrefix("/api/dashboard").Subrouter()
	dashboardRouter.Use(middleware.AuthLogging)
	dashboardRouter.HandleFunc("/view/{cluster_profile_id}", dashboardController.view).Methods("GET")

	kvController := NewKVController()
	kvRouter := router.PathPrefix("/api/kv").Subrouter()
	kvRouter.Use(middleware.AuthLogging)
	kvRouter.HandleFunc("/put/{cluster_profile_id}", kvController.put).Methods("POST")

	return router
}
