package main

import (
	"github.com/aka-achu/watcher/cmd"
	"github.com/aka-achu/watcher/etcd"
	"github.com/aka-achu/watcher/logging"
	"github.com/subosito/gotenv"
)

func init() {
	if err := gotenv.Load(".env"); err != nil {
		logging.Error.Fatalf("Failed to load the env file. %v", err)
	}
	etcd.Initialize()
}
func main() {
	//etcd.C(&model.ClusterProfile{
	//	ID:         "",
	//	Name:       "",
	//	Endpoints:  []string{"http://65.0.213.27:2379","http://65.0.213.27:3379","http://65.0.213.27:1379"},
	//	Username:   "",
	//	Password:   "",
	//	ServerName: "",
	//	CreatedAt:  0,
	//})

	cmd.Execute()
}
