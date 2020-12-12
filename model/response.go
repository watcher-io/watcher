package model

type Format struct {
	ResponseMessage string      `json:"response_message"`
	Data            interface{} `json:"data"`
}
