package repo

import (
	"errors"
	"github.com/aka-achu/watcher/logging"
	"go.etcd.io/bbolt"
	"path/filepath"
	"time"
)

// ValidateSnapshot will validate the user upload db snapshot
// The validation includes checking existence of user profiles and
// cluster profiles.
func ValidateSnapshot() error {
	// Creating a connection with the snapshot db
	snapshotConnection, err := bbolt.Open(filepath.Join("data", "watcher.db.snap"), 0666, &bbolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		logging.Error.Printf(" [DB] Failed to connect to the snapshot database. Error-%v", err)
		return err
	}
	defer snapshotConnection.Close()
	var db = Repo{
		Conn: snapshotConnection,
	}

	// Fetching user details for validation
	user ,err := db.GetUserDetails()
	if err != nil {
		logging.Error.Printf(" [DB] Failed to fetch user profile details from snapshot database. Error-%v", err)
		return err
	}
	if user.InitializationStatus == false || user.Password == "" {
		logging.Error.Printf(" [DB] Invalid database snapshot. Error-%v", err)
		return errors.New("invalid db snapshot. No user profile found")
	}

	// Fetching cluster details for validation
	clusterProfiles, err := db.GetClusterProfiles()
	if err != nil {
		logging.Error.Printf(" [DB] Failed to fetch cluser profiles from snapshot database. Error-%v", err)
		return err
	}
	if len(clusterProfiles) == 0 {
		logging.Error.Printf(" [DB] Invalid database snapshot. Error-%v", err)
		return errors.New("invalid db snapshot. No cluster profile found")
	}
	return nil
}
