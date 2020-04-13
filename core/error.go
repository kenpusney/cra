package core

import (
	"errors"
	"net/http"
	"strconv"
)

var ErrServerStarted = errors.New("server already started")

func UpstreamErrorResponse(response *http.Response) *ResponseError {
	return &ResponseError{
		Type:   strconv.Itoa(response.StatusCode),
		Reason: response.Status,
	}
}

func SystemErrorResponse(err error) *ResponseError {
	return &ResponseError{
		Type:   "xxx",
		Reason: err.Error(),
	}
}
