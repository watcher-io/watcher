package model

type User struct {

	// Password is a key for the admin account
	Password             string `json:"password"`

	// InitializationStatus represents whether the admin account is configured or not.
	InitializationStatus bool `json:"initialization_status"`
}
