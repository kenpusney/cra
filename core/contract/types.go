package contract

import "net/http"

type ResponseCompleter = func(response *Response)

type RequestItem struct {
	Id       string      `json:"id"`
	Endpoint string      `json:"endpoint"`
	Method   string      `json:"method"`
	Type     string      `json:"type"`
	Body     interface{} `json:"body"`
	// only for cascaded mode
	Cascading map[string]string `json:"cascading"`
	Batch     *string           `json:"batch"`
	r         *http.Request
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
	r     *http.Response
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

func (req *Request) AttachOriginalRequest(r *http.Request) {
	for _, request := range req.Requests {
		request.r = r
	}
}

func (req *RequestItem) RequestHeaders() http.Header {
	if req.r != nil {
		return req.r.Header
	}

	return http.Header{}
}

func (res *Response) LatestHeaders() http.Header {
	if len(res.Response) > 0 {
		lastResponse := res.Response[len(res.Response)-1]

		if lastResponse.r != nil {
			return lastResponse.r.Header
		}
	}

	return http.Header{}
}
