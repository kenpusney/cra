package core

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type basicContext struct {
	Opts     *Opts
	Endpoint string
	client   *http.Client
	server   *http.Server
	mux      *http.ServeMux

	strategies map[string]Strategy
}

func NewContext(opts *Opts) Context {
	context := &basicContext{}

	context.Opts = opts

	context.mux = http.NewServeMux()

	context.client = &http.Client{}
	context.server = &http.Server{Addr: context.addr(), Handler: context.mux}
	context.Endpoint = opts.Endpoint

	context.strategies = make(map[string]Strategy)

	context.strategies["seq"] = Sequential
	context.strategies["con"] = Concurrent
	context.strategies["cascaded"] = Cascaded
	context.strategies["batch"] = Batch
	return context
}

func (context *basicContext) Register(ty string, strategy Strategy) {
	context.strategies[ty] = strategy
}

func (context *basicContext) Serve() error {
	context.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var craRequest Request

		UnmarshallJsonObject(ConvertJsonBodyToBytes(r.Body),
			&craRequest)

		context.processRequest(&craRequest, func(response *Response) {
			marshal, err := json.Marshal(response)
			if err != nil {
				return
			}

			_, _ = w.Write(marshal)
		})

	})

	return context.server.ListenAndServe()
}

func (context *basicContext) processRequest(craRequest *Request, completer ResponseCompleter) {

	if craRequest.Mode == "" {
		craRequest.Mode = "seq"
	}

	strategy := context.strategies[craRequest.Mode]
	if strategy != nil {
		strategy(craRequest, context, completer)
	}
}

func (context *basicContext) Proceed(reqItem *RequestItem) *ResponseItem {
	resItem := &ResponseItem{Id: reqItem.Id}

	fillRequests(reqItem)

	var requestBody io.Reader

	if reqItem.Type == "json" {
		b, _ := json.Marshal(reqItem.Body)
		requestBody = bytes.NewReader(b)
	} else {
		requestBody = base64.NewDecoder(base64.RawStdEncoding, strings.NewReader(reqItem.Body.(string)))
	}

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

func (context *basicContext) addr() string {
	if context.Opts.Port > 1000 {
		return ":" + strconv.Itoa(context.Opts.Port)
	}
	return ":9511"
}
