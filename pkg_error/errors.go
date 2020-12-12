package pkg_error

import "errors"

var (
	ErrAdminDoesNotExist         = errors.New("admin profile does not exist")
	ErrFailedToDecodeRequestBody = errors.New("failed to decode the data in request body")
	ErrMissingRequiredFields     = errors.New("required fields are not present in the request")
	ErrFailedToCreateUser        = errors.New("failed to create the user profile")
	ErrUserAlreadyExists         = errors.New("user already exist with same user_name")
	ErrFailedToFetchUser         = errors.New("failed to fetch user profile details")
	ErrInvalidCredential         = errors.New("invalid user credential")
)
