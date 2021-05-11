package context

import (
	"github.com/kenpusney/cra/core/common"
	"github.com/kenpusney/cra/core/contract"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestServerStart(t *testing.T) {
	withCraContext(t, func(context *CoreContext) {
		assert.NotNil(t, context)
	})
}

func TestContextProcessingRequest(t *testing.T) {
	withCraContext(t, func(context *CoreContext) {
		context.processRequest(&contract.Request{
			Id:   "baidu-contract",
			Mode: "test",
			Requests: []*contract.RequestItem{
				{
					Id:       "home",
					Endpoint: "",
				},
			},
		}, func(response *contract.Response) {
			assert.Len(t, response.Response, 1)
		})
	})
}

func TestServerHttpEndpoint(t *testing.T) {
	withCraContext(t, func(context *CoreContext) {
		request, _ := http.NewRequest("POST", "http://127.0.0.1:9511/", contract.EncodeRequestBody(&contract.RequestItem{
			Type: "json",
			Body: contract.Request{
				Id:   "baidu-contract",
				Mode: "test",
				Requests: []*contract.RequestItem{
					{
						Id:       "home",
						Endpoint: "",
					},
				},
			},
		}))

		response, _ := context.client.Do(request)

		assert.NotNil(t, response)
	})
}

func withCraContext(t *testing.T, fn func(context *CoreContext)) {
	stop, wait, context := createCraContext()

	fn(context)

	stop()
	assert.True(t, wait())
}

func createCraContext() (func(), func() bool, *CoreContext) {
	context := MakeContext(&common.Opts{
		Endpoint: "https://www.baidu.com/",
	}, nil)

	context.Register("test", func(craRequest *contract.Request, context common.Context, completer contract.ResponseCompleter) {

		httpReq := &http.Request{}
		httpReq.Header = http.Header{
			"X-POWERED-BY": []string{"CRA"},
		}
		craRequest.AttachOriginalRequest(httpReq)
		response := context.Proceed(craRequest.Requests[0])
		completer(&contract.Response{Response: []*contract.ResponseItem{response}})
	})

	chStop := make(chan string)

	chWait := make(chan bool)

	go func() {
		context.Serve()
		chWait <- true
	}()

	go func() {
		<-chStop
		context.Shutdown()
	}()

	stop := func() {
		time.Sleep(50000000)
		chStop <- "stop"
	}

	wait := func() bool {
		return <-chWait
	}

	return stop, wait, context
}
