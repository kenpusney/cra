package core

import "net/http"

type Opts struct {
	Endpoint string `arg:"positional"`
	//FollowLocation bool   `arg:"-f" default:"false"`
	Port int `arg:"-p" default:"9511"`
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
	Body     interface{} `json:"body"`
	// only for cascaded mode
	Cascading map[string]string `json:"cascading"`
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
