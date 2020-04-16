package core

import (
	"net/http"
	"strconv"
)

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
