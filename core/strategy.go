package core

func Sequential(craRequest *Request, context *Context, completer ResponseCompleter) {
	var craResponses []*ResponseItem
	for _, request := range craRequest.Requests {
		craResponses = append(craResponses,
			context.Proceed(request))
	}

	response := &Response{
		Id:       "id",
		Response: craResponses,
	}

	completer(response)
}

func Concurrent(craRequest *Request, context *Context, completer ResponseCompleter) {
	ch := make(chan *ResponseItem, len(craRequest.Requests))
	var craResponses []*ResponseItem

	for _, request := range craRequest.Requests {
		go func(request *RequestItem) {
			ch <- context.Proceed(request)
		}(request)
	}

	for range craRequest.Requests {
		it := <-ch
		craResponses = append(craResponses, it)
	}

	completer(&Response{
		Id:       "id",
		Response: craResponses,
	})
	close(ch)
}
