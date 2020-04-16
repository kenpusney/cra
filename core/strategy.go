package core

import (
	"encoding/json"
	"github.com/cbroglie/mustache"
	"github.com/google/uuid"
	"github.com/oliveagle/jsonpath"
	"reflect"
)

func Sequential(craRequest *Request, context Context, completer ResponseCompleter) {
	var craResponses []*ResponseItem
	for _, request := range craRequest.Requests {
		request.Id = GenerateId(request.Id, uuid.New().String(), -1)
		craResponses = append(craResponses,
			context.Proceed(request))
	}

	completer(&Response{
		Id:       GenerateId(craRequest.Id, uuid.New().String(), -1),
		Response: craResponses,
	})
}

func Concurrent(craRequest *Request, context Context, completer ResponseCompleter) {
	ch := make(chan *ResponseItem, len(craRequest.Requests))
	var craResponses []*ResponseItem

	for _, request := range craRequest.Requests {
		request.Id = GenerateId(request.Id, uuid.New().String(), -1)
		go func(request *RequestItem) {
			ch <- context.Proceed(request)
		}(request)
	}

	for range craRequest.Requests {
		it := <-ch
		craResponses = append(craResponses, it)
	}

	completer(&Response{
		Id:       GenerateId(craRequest.Id, uuid.New().String(), -1),
		Response: craResponses,
	})
	close(ch)
}

type CascadeContext = map[string]interface{}

func Cascaded(craRequest *Request, context Context, completer ResponseCompleter) {
	var craResponses []*ResponseItem
	var resItem *ResponseItem
	cascadeContext := make(CascadeContext)

	for _, request := range craRequest.Requests {
		request.Id = GenerateId(request.Id, uuid.New().String(), -1)
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

	completer(&Response{
		Id:       GenerateId(craRequest.Id, uuid.New().String(), -1),
		Response: craResponses,
	})
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

func Batch(craRequest *Request, context Context, completer ResponseCompleter) {
	var craResponses []*ResponseItem
	var resItem *ResponseItem

	var seed = make(CascadeContext)

	if craRequest.Seed == nil {
		seed = craRequest.Data
	} else {
		request := craRequest.Seed
		request.Id = GenerateId(request.Id, uuid.New().String(), -1)
		resItem := context.Proceed(request)
		if request.Cascading != nil && len(request.Cascading) != 0 {
			for key, value := range request.Cascading {
				lookup, _ := jsonpath.JsonPathLookup(resItem.Body, value)
				seed[key] = lookup
			}
		}
	}

	for _, originalRequest := range craRequest.Requests {
		originalRequest.Id = GenerateId(originalRequest.Id, uuid.New().String(), -1)

		if originalRequest.Batch == nil {
			break
		}

		name := *originalRequest.Batch
		seedData := asArray(seed[name])

		if len(seedData) > 0 {
			for index, value := range seedData {
				request := rebuildRequestItem(newContext(name, value), originalRequest)
				request.Id = GenerateId(originalRequest.Id, uuid.New().String(), index)
				resItem = context.Proceed(request)
				craResponses = append(craResponses,
					resItem)
			}
		}
	}

	completer(&Response{
		Id:       GenerateId(craRequest.Id, uuid.New().String(), -1),
		Response: craResponses,
	})
}

func newContext(name string, value interface{}) map[string]interface{} {
	return map[string]interface{}{name: value}
}

func asArray(object interface{}) []interface{} {
	var result []interface{}

	switch reflect.TypeOf(object).Kind() {
	case reflect.Slice, reflect.Array:
		value := reflect.ValueOf(object)
		for index := 0; index < value.Len(); index++ {
			result = append(result, value.Index(index).Interface())
		}
	}

	return result
}
