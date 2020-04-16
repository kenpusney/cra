package core

import (
	"encoding/json"
	"net/http"
	"strconv"
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

		UnmarshallJsonObject(ReadBytes(r.Body),
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

	fillRequests(reqItem)

	requestBody := decodeRequestBody(reqItem)

	request, err := http.NewRequest(reqItem.Method, context.Endpoint+reqItem.Endpoint, requestBody)

	if err != nil {
		return &ResponseItem{Id: reqItem.Id, Error: SystemErrorResponse(err)}
	}

	response, err := context.client.Do(request)
	if err != nil {
		return &ResponseItem{Id: reqItem.Id, Error: SystemErrorResponse(err)}
	}

	return formatResponse(response, reqItem.Id)
}

func (context *basicContext) addr() string {
	if context.Opts.Port > 1000 {
		return ":" + strconv.Itoa(context.Opts.Port)
	}
	return ":9511"
}
