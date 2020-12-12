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
	cmd.Execute()
}
