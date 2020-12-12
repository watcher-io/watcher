package model

// Format, is a generic response object for all the handle functions
type Format struct {
	// ResponseMessage is application event message
	ResponseMessage string      `json:"response_message"`

	// Data is a generic interface which will be used to
	// send any type of data to the view layer
	Data            interface{} `json:"data"`
}
