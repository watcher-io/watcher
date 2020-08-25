package main

import (
	"github.com/aka-achu/watcher/cmd"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/repo"
	"github.com/aka-achu/watcher/state"
	"github.com/subosito/gotenv"
)

func init() {
	if err := gotenv.Load(".env"); err != nil {
		logging.Error.Fatalf("Failed to load the env file. %v", err)
	}
	state.Validate()
	repo.Initialize()
}
func main() {
	cmd.Execute()
}
