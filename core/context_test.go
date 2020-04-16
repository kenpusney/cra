package core

import (
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
		context.processRequest(&Request{
			Id: "github-request",
			Requests: []*RequestItem{
				{
					Id:       "home",
					Endpoint: "",
				},
			},
		}, func(response *Response) {
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
	context := NewContext(&Opts{
		Endpoint: "https://api.github.com/",
	}).(*basicContext)

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
