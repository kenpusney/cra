package core

import (
	"context"
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

func (bc *basicContext) Shutdown() {
	_ = bc.server.Shutdown(context.Background())
}

func (bc *basicContext) Register(ty string, strategy Strategy) {
	bc.strategies[ty] = strategy
}

func (bc *basicContext) Serve() error {
	bc.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var craRequest Request

		UnmarshallJsonObject(ReadBytes(r.Body),
			&craRequest)

		bc.processRequest(&craRequest, func(response *Response) {
			marshal, err := json.Marshal(response)
			if err != nil {
				return
			}
			_, _ = w.Write(marshal)
		})

	})

	return bc.server.ListenAndServe()
}

func (bc *basicContext) processRequest(craRequest *Request, completer ResponseCompleter) {

	if craRequest.Mode == "" {
		craRequest.Mode = "seq"
	}

	strategy := bc.strategies[craRequest.Mode]
	if strategy != nil {
		strategy(craRequest, bc, completer)
	}
}

func (bc *basicContext) Proceed(reqItem *RequestItem) *ResponseItem {

	fillRequests(reqItem)

	requestBody := decodeRequestBody(reqItem)

	request, err := http.NewRequest(reqItem.Method, bc.Endpoint+reqItem.Endpoint, requestBody)

	if err != nil {
		return &ResponseItem{Id: reqItem.Id, Error: SystemErrorResponse(err)}
	}

	response, err := bc.client.Do(request)
	if err != nil {
		return &ResponseItem{Id: reqItem.Id, Error: SystemErrorResponse(err)}
	}

	return formatResponse(response, reqItem.Id)
}

func (bc *basicContext) addr() string {
	if bc.Opts.Port > 1000 {
		return ":" + strconv.Itoa(bc.Opts.Port)
	}
	return ":9511"
}
