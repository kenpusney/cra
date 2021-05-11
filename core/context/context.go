package context

import (
	"context"
	"encoding/json"
	"github.com/kenpusney/cra/core/common"
	"github.com/kenpusney/cra/core/contract"
	"github.com/kenpusney/cra/core/util"
	"net/http"
	"strconv"
	"strings"
)

type basicContext struct {
	Opts     *common.Opts
	Endpoint string
	client   *http.Client
	server   *http.Server
	mux      *http.ServeMux

	strategies map[string]common.Strategy

	config *common.Config
}

func MakeContext(opts *common.Opts, config *common.Config) common.Context {
	ctx := &basicContext{}

	ctx.Opts = opts
	ctx.config = config

	ctx.mux = http.NewServeMux()

	ctx.client = &http.Client{}
	ctx.server = &http.Server{Addr: ctx.addr(), Handler: ctx.mux}
	ctx.Endpoint = opts.Endpoint

	ctx.strategies = make(map[string]common.Strategy)
	return ctx
}

func (bc *basicContext) Shutdown() {
	_ = bc.server.Shutdown(context.Background())
}

func (bc *basicContext) Register(ty string, strategy common.Strategy) {
	bc.strategies[ty] = strategy
}

func (bc *basicContext) Serve() error {
	bc.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var craRequest contract.Request

		util.UnmarshallJsonObject(util.ReadBytes(r.Body),
			&craRequest)

		craRequest.AttachOriginalRequest(r)

		bc.processRequest(&craRequest, func(response *contract.Response) {
			marshal, err := json.Marshal(response)
			if err != nil {
				return
			}

			headers := response.LatestHeaders()
			for headerName := range headers {
				if bc.shouldByPassHeader(headerName) {
					w.Header().Set(headerName, headers.Get(headerName))
				}
			}
			_, _ = w.Write(marshal)
		})

	})

	return bc.server.ListenAndServe()
}

func (bc *basicContext) processRequest(craRequest *contract.Request, completer contract.ResponseCompleter) {

	if craRequest.Mode == "" {
		craRequest.Mode = "seq"
	}

	strategy := bc.strategies[craRequest.Mode]
	if strategy != nil {
		strategy(craRequest, bc, completer)
	}
}

func (bc *basicContext) Proceed(reqItem *contract.RequestItem) *contract.ResponseItem {

	contract.FillRequest(reqItem)

	requestBody := contract.DecodeRequestBody(reqItem)

	req, err := http.NewRequest(reqItem.Method, bc.Endpoint+reqItem.Endpoint, requestBody)

	if err != nil {
		return &contract.ResponseItem{Id: reqItem.Id, Error: contract.SystemErrorResponse(err)}
	}

	headers := reqItem.RequestHeaders()

	for headerName := range headers {
		if bc.shouldByPassHeader(headerName) {
			req.Header.Set(headerName, headers.Get(headerName))
		}
	}

	response, err := bc.client.Do(req)
	if err != nil {
		return contract.ErrorResponse(reqItem, err, response)
	}

	return contract.FormatResponse(response, reqItem.Id)
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
