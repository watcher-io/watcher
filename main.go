package main

import (
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/validation"
	"github.com/subosito/gotenv"
)

func init() {
	if err := gotenv.Load(".env"); err != nil {
		logging.Error.Fatalf("Failed to load the env file. %v", err)
	}
	validation.Validate()
}
func main(){

}
