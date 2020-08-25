package controller

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Initialize() *mux.Router {
	router := mux.NewRouter()
	router.Use(cors.AllowAll().Handler)

	var authController AuthController

	router.HandleFunc("/api/auth/checkAdminStatus", authController.CheckAdminInitStatus)
	router.HandleFunc("/api/auth/saveAdminProfile", authController.SaveAdminProfile)
	router.HandleFunc("/api/auth/login", authController.Login)

	return router
}
