package model

// LoginResponse is a object used for sending the generated JWT token to the client
// after successful user validation
type LoginResponse struct {
	Token string `json:"token"`
}
