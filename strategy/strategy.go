package strategy

import "github.com/kenpusney/cra/core"

func Sequential(craRequest *core.Request, context *core.Context, completer core.ResponseCompleter) {
	var craResponses []*core.ResponseItem
	for _, request := range craRequest.Requests {
		craResponses = append(craResponses,
			context.Proceed(request))
	}

	response := &core.Response{
		Id:       "id",
		Response: craResponses,
	}

	completer(response)
}

func Concurrent(craRequest *core.Request, context *core.Context, completer core.ResponseCompleter) {
	ch := make(chan *core.ResponseItem, len(craRequest.Requests))
	var craResponses []*core.ResponseItem

	for _, request := range craRequest.Requests {
		go func(request *core.RequestItem) {
			ch <- context.Proceed(request)
		}(request)
	}

	for range craRequest.Requests {
		it := <-ch
		craResponses = append(craResponses, it)
	}

	completer(&core.Response{
		Id:       "id",
		Response: craResponses,
	})
	close(ch)
}
