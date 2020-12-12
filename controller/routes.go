package controller

import (
	"github.com/aka-achu/watcher/middleware"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/repository"
	"github.com/aka-achu/watcher/service"
	"github.com/gorilla/mux"
)

func Initialize() *mux.Router {

	router := mux.NewRouter()
	db := repository.NewDatabase()

	registerUserRoute(
		router,
		NewUserController(),
		repository.NewUserRepo(db.Conn),
		service.NewUserService(),
	)
	registerAuthRoute(
		router,
		NewAuthController(),
		repository.NewUserRepo(db.Conn),
		service.NewAuthService(),
	)
	registerClusterProfileRoute(
		router,
		NewClusterProfileController(),
		repository.NewClusterProfileRepo(db.Conn),
		service.NewClusterProfileService(),
	)
	registerDashboardRoute(
		router,
		NewDashboardController(),
		repository.NewClusterProfileRepo(db.Conn),
		service.NewDashboardService(),
	)
	registerKVRoute(
		router,
		NewKVController(),
		repository.NewClusterProfileRepo(db.Conn),
		service.NewKVService(),
	)

	return router
}

func registerUserRoute(
	r *mux.Router,
	controller model.UserController,
	repo model.UserRepo,
	svc model.UserService,
) {
	var userRouter = r.PathPrefix("/api/v1/user").Subrouter()
	userRouter.Use(middleware.NoAuthLogging)
	userRouter.HandleFunc("/create", controller.Create(repo, svc)).Methods("POST")
	userRouter.HandleFunc("/exists", controller.Exists(repo, svc)).Methods("GET")
}

func registerAuthRoute(
	r *mux.Router,
	controller model.AuthController,
	repo model.UserRepo,
	svc model.AuthService,
) {
	var authRouter = r.PathPrefix("/api/v1/auth").Subrouter()
	authRouter.Use(middleware.NoAuthLogging)
	authRouter.HandleFunc("/login", controller.Login(repo, svc)).Methods("POST")
}

func registerClusterProfileRoute(
	r *mux.Router,
	controller model.ClusterProfileController,
	repo model.ClusterProfileRepo,
	svc model.ClusterProfileService,
) {
	var clusterProfileRouter = r.PathPrefix("/api/v1/clusterProfile").Subrouter()
	clusterProfileRouter.Use(middleware.NoAuthLogging)
	clusterProfileRouter.HandleFunc("/create", controller.Create(repo, svc)).Methods("POST")
	clusterProfileRouter.HandleFunc("/fetch", controller.Fetch(repo, svc)).Methods("GET")
}

func registerDashboardRoute(
	r *mux.Router,
	controller model.DashboardController,
	repo model.ClusterProfileRepo,
	svc model.DashboardService,
) {
	var dashboardRouter = r.PathPrefix("/api/v1/dashboard").Subrouter()
	dashboardRouter.Use(middleware.NoAuthLogging)
	dashboardRouter.HandleFunc("/view/{cluster_profile_id}", controller.View(repo, svc)).Methods("GET")
}

func registerKVRoute(
	r *mux.Router,
	controller model.KVController,
	repo model.ClusterProfileRepo,
	svc model.KVService,
) {
	var kvRouter = r.PathPrefix("/api/v1/kv").Subrouter()
	kvRouter.Use(middleware.NoAuthLogging)
	kvRouter.HandleFunc("/put/{cluster_profile_id}", controller.Put(repo, svc)).Methods("POST")
}
