package core

import (
	"context"
	"encoding/json"
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

	config *Config
}

func NewContext(opts *Opts) Context {
	ctx := &basicContext{}

	ctx.Opts = opts
	ctx.config = LoadConfig()

	ctx.mux = http.NewServeMux()

	ctx.client = &http.Client{}
	ctx.server = &http.Server{Addr: ctx.addr(), Handler: ctx.mux}
	ctx.Endpoint = opts.Endpoint

	ctx.strategies = make(map[string]Strategy)

	ctx.strategies["seq"] = Sequential
	ctx.strategies["con"] = Concurrent
	ctx.strategies["cascaded"] = Cascaded
	ctx.strategies["batch"] = Batch
	return ctx
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

		for _, request := range craRequest.Requests {
			request.r = r
		}

		bc.processRequest(&craRequest, func(response *Response) {
			marshal, err := json.Marshal(response)
			if err != nil {
				return
			}

			if len(response.Response) > 0 {
				lastResponse := response.Response[len(response.Response)-1]

				if lastResponse.r != nil {
					for header := range lastResponse.r.Header {
						if bc.shouldByPassHeader(header) {
							w.Header().Set(header, lastResponse.r.Header.Get(header))
						}
					}
				}
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

	if reqItem.r != nil {
		for header := range reqItem.r.Header {
			if bc.shouldByPassHeader(header) {
				request.Header.Set(header, reqItem.r.Header.Get(header))
			}
		}
	}

	response, err := bc.client.Do(request)
	if err != nil {
		return &ResponseItem{Id: reqItem.Id, Error: SystemErrorResponse(err), r: response}
	}

	return formatResponse(response, reqItem.Id)
}

func (bc *basicContext) addr() string {
	if bc.Opts.Port > 1000 {
		return ":" + strconv.Itoa(bc.Opts.Port)
	}
	return ":9511"
}

func (bc *basicContext) shouldByPassHeader(header string) bool {
	if strings.HasPrefix(header, "X-") {
		return true
	}

	if bc.config != nil && bc.config.BypassedHeaders != nil {
		for _, bypassConfig := range *bc.config.BypassedHeaders {
			if strings.EqualFold(header, bypassConfig) {
				return true
			}
		}
	}

	return false
}
