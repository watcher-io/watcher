package controller

import (
	"github.com/aka-achu/watcher/middleware"
	"github.com/aka-achu/watcher/model"
	"github.com/aka-achu/watcher/repo"
	"github.com/aka-achu/watcher/service"
	"github.com/gorilla/mux"
)

func Initialize() *mux.Router {

	router := mux.NewRouter()
	db := repo.NewDatabase()

	registerUserRoute(router, NewUserController(), repo.NewUserRepo(db.Conn), service.NewUserService())
	registerAuthRoute(router, NewAuthController(), repo.NewUserRepo(db.Conn), service.NewAuthService())
	registerClusterProfileRoute(router, NewClusterProfileController(), repo.NewClusterProfileRepo(db.Conn), service.NewClusterProfileService())
	registerDashboardRoute(router, NewDashboardController(), repo.NewClusterProfileRepo(db.Conn), service.NewDashboardService())
	registerKVRoute(router, NewKVController(), repo.NewClusterProfileRepo(db.Conn), service.NewKVService())

	return router
}

func registerUserRoute(
	r *mux.Router,
	userController model.UserController,
	userRepo model.UserRepo,
	userService model.UserService,
) {
	var userRouter = r.PathPrefix("/api/v1/user").Subrouter()
	userRouter.Use(middleware.NoAuthLogging)
	userRouter.HandleFunc("/create", userController.Create(userRepo, userService)).Methods("POST")
	userRouter.HandleFunc("/exists", userController.Exists(userRepo, userService)).Methods("GET")
}

func registerAuthRoute(
	r *mux.Router,
	authController model.AuthController,
	userRepo model.UserRepo,
	authService model.AuthService,
) {
	var authRouter = r.PathPrefix("/api/v1/auth").Subrouter()
	authRouter.Use(middleware.NoAuthLogging)
	authRouter.HandleFunc("/login", authController.Login(userRepo, authService)).Methods("POST")
}

func registerClusterProfileRoute(
	r *mux.Router,
	clusterProfileController model.ClusterProfileController,
	clusterProfileRepo model.ClusterProfileRepo,
	clusterProfileService model.ClusterProfileService,
) {
	var clusterProfileRouter = r.PathPrefix("/api/v1/clusterProfile").Subrouter()
	clusterProfileRouter.Use(middleware.NoAuthLogging)
	clusterProfileRouter.HandleFunc("/create", clusterProfileController.Create(clusterProfileRepo, clusterProfileService)).Methods("POST")
	clusterProfileRouter.HandleFunc("/fetch", clusterProfileController.Fetch(clusterProfileRepo, clusterProfileService)).Methods("GET")
}

func registerDashboardRoute(
	r *mux.Router,
	dashboardController model.DashboardController,
	clusterProfileRepo model.ClusterProfileRepo,
	dashboardService model.DashboardService,
) {
	var dashboardRouter = r.PathPrefix("/api/v1/dashboard").Subrouter()
	dashboardRouter.Use(middleware.NoAuthLogging)
	dashboardRouter.HandleFunc("/view/{cluster_profile_id}", dashboardController.View(clusterProfileRepo, dashboardService)).Methods("GET")
}

func registerKVRoute(
	r *mux.Router,
	kvController model.KVController,
	clusterProfileRepo model.ClusterProfileRepo,
	kvService model.KVService,
) {
	var kvRouter = r.PathPrefix("/api/v1/kv").Subrouter()
	kvRouter.Use(middleware.NoAuthLogging)
	kvRouter.HandleFunc("/put/{cluster_profile_id}", kvController.Put(clusterProfileRepo, kvService)).Methods("POST")

}
