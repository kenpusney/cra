package core

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func NewContext(opts *Opts) *Context {
	context := &Context{}

	context.Opts = opts

	context.mux = http.NewServeMux()

	context.client = &http.Client{}
	context.server = &http.Server{Addr: context.addr(), Handler: context.mux}

	return context
}

func (context *Context) Serve() error {
	context.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var craRequest Request

		UnmarshallJsonObject(ConvertJsonBodyToBytes(r.Body),
			&craRequest)

		context.processRequest(&craRequest, func(response *Response) {
			marshal, err := json.Marshal(response)
			if err != nil {
				return
			}

			w.Write(marshal)
		})

	})

	return context.server.ListenAndServe()
}

func (context *Context) processRequest(craRequest *Request, completer ResponseCompleter) {
	Concurrent(craRequest, context, completer)

}

func (context *Context) Proceed(reqItem *RequestItem) *ResponseItem {
	resItem := &ResponseItem{Id: "internal_id_"}

	fillRequests(reqItem)

	var requestBody io.Reader

	request, err := http.NewRequest(reqItem.Method, context.Endpoint+reqItem.Endpoint, requestBody)
	if err != nil {
		resItem.Error = SystemErrorResponse(err)
		return resItem
	}

	response, err := context.client.Do(request)
	if err != nil {
		resItem.Error = SystemErrorResponse(err)
		return resItem
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		resItem.Error = UpstreamErrorResponse(response)
	}

	if strings.Contains(response.Header.Get("Content-Type"), "json") {
		resItem.Type = "json"
		resItem.Body = ConvertJsonBodyToObject(response.Body)
	} else {
		resItem.Type = "binary"
		resItem.Body = ConvertJsonBodyToBytes(response.Body)
	}
	if resItem.Error != nil {
		resItem.Error.Body = resItem.Body
	}
	return resItem
}
func fillRequests(req *RequestItem) {
	if req.Method == "" {
		req.Method = "GET"
	}
}

func (context *Context) addr() string {
	if context.Opts.Port > 1000 {
		return ":" + strconv.Itoa(context.Opts.Port)
	}
	return ":9511"
}
