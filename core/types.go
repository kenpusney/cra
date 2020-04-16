package core

type Opts struct {
	Endpoint string `arg:"positional"`
	//FollowLocation bool   `arg:"-f" default:"false"`
	Port int `arg:"-p" default:"9511"`
}

type Strategy func(craRequest *Request, context Context, completer ResponseCompleter)

type ResponseCompleter = func(response *Response)

type Context interface {
	Register(ty string, strategy Strategy)
	Serve() error
	Proceed(reqItem *RequestItem) *ResponseItem
	Shutdown()
}

type RequestItem struct {
	Id       string      `json:"id"`
	Endpoint string      `json:"endpoint"`
	Method   string      `json:"method"`
	Type     string      `json:"type"`
	Body     interface{} `json:"body"`
	// only for cascaded mode
	Cascading map[string]string `json:"cascading"`
	Batch     *string           `json:"batch"`
}

type Request struct {
	Mode     string                 `json:"mode"`
	Id       string                 `json:"id"`
	Requests []*RequestItem         `json:"requests"`
	Seed     *RequestItem           `json:"seed"`
	Data     map[string]interface{} `json:"data"`
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
