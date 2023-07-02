package common

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}

type Request struct {
	RequestId   string    `json:"requestId"`
	RequestTime time.Time `json:"requestTime"`
	Signature   string    `json:"signature"`
}
type Response struct {
	ResponseID   string    `json:"responseId"`
	ResponseTime time.Time `json:"responseTime"`
}
type UserRequest struct {
	Request Request `json:"request"`
	Data    User    `json:"data" validate:"required"`
}
type UserResponse struct {
	ResponseCode    string   `json:"responseCode"`
	Response        Response `json:"response"`
	ResponseMessage string   `json:"responseMessage"`
}
