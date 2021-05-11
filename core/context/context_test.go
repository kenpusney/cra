package context

import (
	"github.com/kenpusney/cra/core/common"
	"github.com/kenpusney/cra/core/contract"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestServerStart(t *testing.T) {
	withCraContext(t, func(context *basicContext) {
		assert.NotNil(t, context)
	})
}

func TestContextProcessingRequest(t *testing.T) {
	withCraContext(t, func(context *basicContext) {
		context.processRequest(&contract.Request{
			Id: "github-contract",
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

func withCraContext(t *testing.T, fn func(context *basicContext)) {
	stop, wait, context := createCraContext()

	fn(context)

	stop()
	assert.True(t, wait())
}

func createCraContext() (func(), func() bool, *basicContext) {
	context := MakeContext(&common.Opts{
		Endpoint: "https://api.github.com/",
	}, nil).(*basicContext)

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
