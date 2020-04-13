package core

import "net/http"

type Opts struct {
	FollowLocation bool
	Port           int
}

type Strategy = func(craRequest *Request, context *Context, completer ResponseCompleter)

type ResponseCompleter = func(response *Response)

type Context struct {
	Opts     *Opts
	Endpoint string
	client   *http.Client
	server   *http.Server
	mux      *http.ServeMux

	strategies map[string]Strategy
}

type RequestItem struct {
	Id       string      `json:"id"`
	Endpoint string      `json:"endpoint"`
	Method   string      `json:"method"`
	Type     string      `json:"type"`
	Json     interface{} `json:"json"`
	Bytes    string      `json:"bytes"`
}

type Request struct {
	Mode     string         `json:"mode"`
	Id       string         `json:"id"`
	Requests []*RequestItem `json:"requests"`
}

type ResponseItem struct {
	Id    string         `json:"id"`
	Type  string         `json:"type"`
	Error *ResponseError `json:"error"`
	Body  interface{}    `json:"body"`
}

type ResponseError struct {
	Type   string      `json:"type"`
	Reason string      `json:"reason"`
	Body   interface{} `json:"body"`
}

type Response struct {
	Id       string          `json:"id"`
	Response []*ResponseItem `json:"response"`
}
