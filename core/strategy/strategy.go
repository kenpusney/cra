package strategy

import (
	"encoding/json"
	"github.com/cbroglie/mustache"
	"github.com/google/uuid"
	"github.com/kenpusney/cra/core/common"
	"github.com/kenpusney/cra/core/contract"
	"github.com/kenpusney/cra/core/util"
	"github.com/oliveagle/jsonpath"
	"reflect"
)

func Sequential(craRequest *contract.Request, context common.Context, completer contract.ResponseCompleter) {
	var craResponses []*contract.ResponseItem
	for _, r := range craRequest.Requests {
		r.Id = util.GenerateId(r.Id, uuid.New().String(), -1)
		craResponses = append(craResponses,
			context.Proceed(r))
	}

	completer(&contract.Response{
		Id:       util.GenerateId(craRequest.Id, uuid.New().String(), -1),
		Response: craResponses,
	})
}

func Concurrent(craRequest *contract.Request, context common.Context, completer contract.ResponseCompleter) {
	ch := make(chan *contract.ResponseItem, len(craRequest.Requests))
	var craResponses []*contract.ResponseItem

	for _, r := range craRequest.Requests {
		r.Id = util.GenerateId(r.Id, uuid.New().String(), -1)
		go func(request *contract.RequestItem) {
			ch <- context.Proceed(request)
		}(r)
	}

	for range craRequest.Requests {
		it := <-ch
		craResponses = append(craResponses, it)
	}

	completer(&contract.Response{
		Id:       util.GenerateId(craRequest.Id, uuid.New().String(), -1),
		Response: craResponses,
	})
	close(ch)
}

type CascadeContext = map[string]interface{}

func Cascaded(craRequest *contract.Request, context common.Context, completer contract.ResponseCompleter) {
	var craResponses []*contract.ResponseItem
	var resItem *contract.ResponseItem
	cascadeContext := make(CascadeContext)

	for _, request := range craRequest.Requests {
		request.Id = util.GenerateId(request.Id, uuid.New().String(), -1)
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

	completer(&contract.Response{
		Id:       util.GenerateId(craRequest.Id, uuid.New().String(), -1),
		Response: craResponses,
	})
}

func rebuildRequestItem(cascadeContext CascadeContext, req *contract.RequestItem) *contract.RequestItem {

	endPoint, _ := mustache.Render(req.Endpoint, cascadeContext)
	marshal, _ := json.Marshal(req.Body)
	body, _ := mustache.Render(string(marshal), cascadeContext)

	var jsonBody interface{}

	_ = json.Unmarshal([]byte(body), jsonBody)
	newRequest := &contract.RequestItem{
		Id:       req.Id,
		Endpoint: endPoint,
		Method:   req.Method,
		Type:     req.Type,
		Body:     jsonBody,
	}

	return newRequest
}

func Batch(craRequest *contract.Request, context common.Context, completer contract.ResponseCompleter) {
	var craResponses []*contract.ResponseItem
	var resItem *contract.ResponseItem

	var seed = make(CascadeContext)

	if craRequest.Seed == nil {
		seed = craRequest.Data
	} else {
		request := craRequest.Seed
		request.Id = util.GenerateId(request.Id, uuid.New().String(), -1)
		resItem := context.Proceed(request)
		if request.Cascading != nil && len(request.Cascading) != 0 {
			for key, value := range request.Cascading {
				lookup, _ := jsonpath.JsonPathLookup(resItem.Body, value)
				seed[key] = lookup
			}
		}
	}

	for _, originalRequest := range craRequest.Requests {
		originalRequest.Id = util.GenerateId(originalRequest.Id, uuid.New().String(), -1)

		if originalRequest.Batch == nil {
			break
		}

		name := *originalRequest.Batch
		seedData := asArray(seed[name])

		if len(seedData) > 0 {
			for index, value := range seedData {
				request := rebuildRequestItem(newContext(name, value), originalRequest)
				request.Id = util.GenerateId(originalRequest.Id, uuid.New().String(), index)
				resItem = context.Proceed(request)
				craResponses = append(craResponses,
					resItem)
			}
		}
	}

	completer(&contract.Response{
		Id:       util.GenerateId(craRequest.Id, uuid.New().String(), -1),
		Response: craResponses,
	})
}

func newContext(name string, value interface{}) map[string]interface{} {
	return map[string]interface{}{name: value}
}

func asArray(object interface{}) []interface{} {
	var result []interface{}

	if object == nil {
		return result
	}

	switch reflect.TypeOf(object).Kind() {
	case reflect.Slice, reflect.Array:
		value := reflect.ValueOf(object)
		for index := 0; index < value.Len(); index++ {
			result = append(result, value.Index(index).Interface())
		}
	}

	return result
}
