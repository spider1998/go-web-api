package conf

import (
	"encoding/json"
	"net/http"
)

func NewAPIError(code Code) APIError {
	var err APIError
	return err.WithCode(code)
}

type APIError struct {
	status  int         `json:"-"`
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e APIError) StatusCode() int {
	return e.status
}

func (e APIError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (e APIError) WithStatus(status int) APIError {
	e.status = status
	return e
}

func (e APIError) WithCode(code Code) APIError {
	e.Code = code
	e.Message = string(code)
	if status, ok := codeStatusMap[e.Code]; ok {
		e.status = status
	} else {
		e.status = http.StatusInternalServerError
	}
	return e
}

func (e APIError) WithMessage(message string) APIError {
	e.Message = message
	return e
}

func (e APIError) WithDetails(details interface{}) APIError {
	e.Details = details
	return e
}
