package conf

import (
	"net/http"
)

type Code string

const (
	InternalServerError Code = "INTERNAL_SERVER_ERROR"
	NotFound                 = "NOT_FOUND"
	MethodNotAllowed         = "METHOD_NOT_ALLOWED"
	InvalidData              = "INVALID_DATA"
)

var codeStatusMap = map[Code]int{
	InternalServerError: http.StatusInternalServerError,
	NotFound:            http.StatusNotFound,
	MethodNotAllowed:    http.StatusMethodNotAllowed,
	InvalidData:         http.StatusBadRequest,
}
