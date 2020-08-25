package model

// User signifies the admin profile object
type User struct {

	// Password is a key for the admin account
	// The user given key will be hashed  and stored in the database
	Password string `json:"password"`

	// InitializationStatus represents whether the admin account is initialized or not
	InitializationStatus bool `json:"initialization_status"`
}
