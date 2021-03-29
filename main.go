package main

import (
	"github.com/subosito/gotenv"
	"github.com/watcher-io/watcher/cmd"
	"github.com/watcher-io/watcher/etcd"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/repository"
	"log"
	"os"
)

func init() {
	if len(os.Args) != 2 {
		logging.Error.Fatalf(" [APP] Invalid of nil runtime arguments")
	}

	switch os.Args[1] {
	case "dev":
		if err := gotenv.Load(".env.dev"); err != nil {
			log.Fatalf(" [APP] Failed to load the dev env file. Err-%s", err.Error())
		}
	case "prod":
		if err := gotenv.Load(".env.dev"); err != nil {
			log.Fatalf(" [APP] Failed to load the production env file. Err-%s", err.Error())
		}
	default:
		log.Fatal(" [APP] Invalid argument")
	}

	repository.Initialize()
	etcd.Initialize()
}
func main() {
	cmd.Execute()
}
