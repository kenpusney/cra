package core

import (
	"encoding/json"
	"github.com/cbroglie/mustache"
	"github.com/oliveagle/jsonpath"
	"reflect"
)

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

type CascadeContext = map[string]interface{}

func Cascaded(craRequest *Request, context *Context, completer ResponseCompleter) {
	var craResponses []*ResponseItem
	var resItem *ResponseItem
	cascadeContext := make(CascadeContext)

	for _, request := range craRequest.Requests {
		if len(cascadeContext) > 0 {
			request = rebuildRequestItem(cascadeContext, request)
		}
		resItem = context.Proceed(request)
		craResponses = append(craResponses,
			resItem)
		if request.Cascading != nil && len(request.Cascading) != 0 {
			for key, value := range request.Cascading {
				lookup, _ := jsonpath.JsonPathLookup(resItem.Body, value)
				cascadeContext[key] = lookup
			}
		}
	}

	response := &Response{
		Id:       "id",
		Response: craResponses,
	}

	completer(response)
}

func rebuildRequestItem(cascadeContext CascadeContext, request *RequestItem) *RequestItem {

	endPoint, _ := mustache.Render(request.Endpoint, cascadeContext)
	marshal, _ := json.Marshal(request.Body)
	body, _ := mustache.Render(string(marshal), cascadeContext)

	var jsonBody interface{}

	_ = json.Unmarshal([]byte(body), jsonBody)
	newRequest := &RequestItem{
		Id:       request.Id,
		Endpoint: endPoint,
		Method:   request.Method,
		Type:     request.Type,
		Body:     jsonBody,
	}

	return newRequest
}

func Batch(craRequest *Request, context *Context, completer ResponseCompleter) {
	var craResponses []*ResponseItem
	var resItem *ResponseItem

	var seed = make(CascadeContext)

	if craRequest.Seed == nil {
		seed = craRequest.Data
	} else {
		request := craRequest.Seed
		resItem := context.Proceed(request)
		if request.Cascading != nil && len(request.Cascading) != 0 {
			for key, value := range request.Cascading {
				lookup, _ := jsonpath.JsonPathLookup(resItem.Body, value)
				seed[key] = lookup
			}
		}
	}

	for _, originalRequest := range craRequest.Requests {
		if originalRequest.Batch == nil {
			break
		}

		data := asArray(seed[*originalRequest.Batch])

		if len(data) > 0 {
			for _, v := range data {
				request := rebuildRequestItem(map[string]interface{}{
					*originalRequest.Batch: v,
				}, originalRequest)
				resItem = context.Proceed(request)
				craResponses = append(craResponses,
					resItem)
			}
		}
	}

	response := &Response{
		Id:       "id",
		Response: craResponses,
	}
	completer(response)
}

func asArray(object interface{}) []interface{} {
	var result []interface{}

	switch reflect.TypeOf(object).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(object)
		for i := 0; i < s.Len(); i++ {
			result = append(result, s.Index(i).Interface())
		}
	}

	return result
}
