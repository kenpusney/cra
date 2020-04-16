package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestContext struct{}

func (TestContext) Register(ty string, strategy Strategy) {
	panic("implement me")
}

func (TestContext) Serve() error {
	panic("implement me")
}

func (TestContext) Proceed(reqItem *RequestItem) *ResponseItem {
	return &ResponseItem{
		Id: reqItem.Id,
		Body: map[string]interface{}{
			"value": "test",
			"data":  []int{1, 2, 3},
		},
	}
}

func TestSequential(t *testing.T) {
	Sequential(&Request{
		Id:   "seq",
		Mode: "seq",
		Requests: []*RequestItem{
			{
				Id: "1",
			},
			{
				Id: "2",
			},
		},
	}, &TestContext{}, func(response *Response) {
		assert.NotEmpty(t, response.Response)
		assert.Equal(t, response.Id, "seq")
		assert.Len(t, response.Response, 2)
	})
}

func TestConcurrent(t *testing.T) {
	Concurrent(&Request{
		Id:   "con",
		Mode: "con",
		Requests: []*RequestItem{
			{
				Id: "1",
			},
			{
				Id: "2",
			},
		},
	}, &TestContext{}, func(response *Response) {
		assert.NotEmpty(t, response.Response)
		assert.Equal(t, response.Id, "con")
		assert.Len(t, response.Response, 2)
	})
}

func TestCascaded(t *testing.T) {
	Cascaded(&Request{
		Id:   "cascaded",
		Mode: "cascaded",
		Requests: []*RequestItem{
			{
				Id: "1",
				Cascading: map[string]string{
					"value": "$.value",
				},
			},
			{
				Id: "2",
				Body: map[string]string{
					"value": "{{value}}",
				},
			},
		},
	}, &TestContext{}, func(response *Response) {
		assert.NotEmpty(t, response.Response)
		assert.Equal(t, response.Id, "cascaded")
		assert.Len(t, response.Response, 2)
	})
}

func TestBatchWithBatchRequest(t *testing.T) {

	var batchName = "values"

	Batch(&Request{
		Id: "batch",
		Seed: &RequestItem{
			Id: "seed",
			Cascading: map[string]string{
				"values": "$.data",
			},
		},
		Mode: "batch",
		Requests: []*RequestItem{
			{
				Id:    "batch",
				Batch: &batchName,
			},
		},
	}, &TestContext{}, func(response *Response) {
		assert.NotEmpty(t, response.Response)
		assert.Equal(t, response.Id, "batch")
		assert.Len(t, response.Response, 3)
	})
}

func TestBatchWithData(t *testing.T) {

	var batchName = "values"

	Batch(&Request{
		Id: "batch",
		Data: CascadeContext{
			"values": []string{"a", "b", "c"},
		},
		Mode: "batch",
		Requests: []*RequestItem{
			{
				Id:    "batch",
				Batch: &batchName,
			},
		},
	}, &TestContext{}, func(response *Response) {
		assert.NotEmpty(t, response.Response)
		assert.Equal(t, response.Id, "batch")
		assert.Len(t, response.Response, 3)
	})
}
