package controller

import (
	"github.com/aka-achu/watcher/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Initialize() *mux.Router {
	router := mux.NewRouter()
	router.Use(cors.AllowAll().Handler)

	var authController AuthController
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.Use(middleware.NoAuthLogging)
	authRouter.HandleFunc("/checkAdminStatus", authController.CheckAdminInitStatus)
	authRouter.HandleFunc("/saveAdminProfile", authController.SaveAdminProfile)
	authRouter.HandleFunc("/login", authController.Login)

	var clusterController ClusterController
	clusterRouter := router.PathPrefix("/api/cluster").Subrouter()
	clusterRouter.Use(middleware.AuthLogging)
	clusterRouter.HandleFunc("/fetch", clusterController.FetchClusterProfiles)
	clusterRouter.HandleFunc("/create", clusterController.CreateClusterProfile)

	return router
}
