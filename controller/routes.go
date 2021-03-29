package controller

import (
	"github.com/gorilla/mux"
	"github.com/watcher-io/watcher/middleware"
	"github.com/watcher-io/watcher/model"
	"github.com/watcher-io/watcher/repository"
	"github.com/watcher-io/watcher/service"
	"github.com/watcher-io/watcher/store"
)

func Initialize() *mux.Router {

	router := mux.NewRouter()
	db := repository.NewDatabase()
	objectStore := store.NewObjectStore()

	registerUserRoute(
		router,
		NewUserController(
			service.NewUserService(repository.NewUserRepo(db.Conn)),
		),
	)
	registerAuthRoute(
		router,
		NewAuthController(
			service.NewAuthService(repository.NewUserRepo(db.Conn)),
		),
	)
	registerClusterProfileRoute(
		router,
		NewClusterProfileController(
			service.NewClusterProfileService(repository.NewClusterProfileRepo(db.Conn), objectStore),
		),
	)
	registerDashboardRoute(
		router,
		NewDashboardController(
			service.NewDashboardService(repository.NewClusterProfileRepo(db.Conn), objectStore),
		),
	)
	registerKVRoute(
		router,
		NewKVController(
			service.NewKVService(repository.NewClusterProfileRepo(db.Conn), objectStore),
		),
	)
	return router
}

func registerUserRoute(
	r *mux.Router,
	controller model.UserController,
) {
	var userRouter = r.PathPrefix("/api/v1/user").Subrouter()
	userRouter.Use(middleware.NoAuthLogging)
	userRouter.HandleFunc("/create", controller.Create()).Methods("POST")
	userRouter.HandleFunc("/exists", controller.Exists()).Methods("GET")
}

func registerAuthRoute(
	r *mux.Router,
	controller model.AuthController,
) {
	var authRouter = r.PathPrefix("/api/v1/auth").Subrouter()
	authRouter.Use(middleware.NoAuthLogging)
	authRouter.HandleFunc("/login", controller.Login()).Methods("POST")
}

func registerClusterProfileRoute(
	r *mux.Router,
	controller model.ClusterProfileController,
) {
	var clusterProfileRouter = r.PathPrefix("/api/v1/clusterProfile").Subrouter()
	clusterProfileRouter.Use(middleware.NoAuthLogging)
	clusterProfileRouter.HandleFunc("/create", controller.Create()).Methods("POST")
	clusterProfileRouter.HandleFunc("/fetch", controller.Fetch()).Methods("GET")
	clusterProfileRouter.HandleFunc("/uploadCertificate", controller.UploadCertificate()).Methods("POST")
}

func registerDashboardRoute(
	r *mux.Router,
	controller model.DashboardController,
) {
	var dashboardRouter = r.PathPrefix("/api/v1/dashboard").Subrouter()
	dashboardRouter.Use(middleware.NoAuthLogging)
	dashboardRouter.HandleFunc("/view/{cluster_profile_id}", controller.View()).Methods("GET")
}

func registerKVRoute(
	r *mux.Router,
	controller model.KVController,
) {
	var kvRouter = r.PathPrefix("/api/v1/kv").Subrouter()
	kvRouter.Use(middleware.NoAuthLogging)
	kvRouter.HandleFunc("/put/{cluster_profile_id}", controller.Put()).Methods("POST")
}
