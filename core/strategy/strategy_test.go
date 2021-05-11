package strategy

import (
	"github.com/kenpusney/cra/core/common"
	"github.com/kenpusney/cra/core/contract"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestContext struct{}

func (TestContext) Register(ty string, strategy common.Strategy) {
	panic("implement me")
}

func (TestContext) Serve() error {
	panic("implement me")
}

func (TestContext) Shutdown() {
	panic("implement me")
}

func (TestContext) Proceed(reqItem *contract.RequestItem) *contract.ResponseItem {
	return &contract.ResponseItem{
		Id: reqItem.Id,
		Body: map[string]interface{}{
			"value": "test",
			"data":  []int{1, 2, 3},
		},
	}
}

func TestSequential(t *testing.T) {
	Sequential(&contract.Request{
		Id:   "seq",
		Mode: "seq",
		Requests: []*contract.RequestItem{
			{
				Id: "1",
			},
			{
				Id: "2",
			},
		},
	}, &TestContext{}, func(response *contract.Response) {
		assert.NotEmpty(t, response.Response)
		assert.Equal(t, response.Id, "seq")
		assert.Len(t, response.Response, 2)
	})
}

func TestConcurrent(t *testing.T) {
	Concurrent(&contract.Request{
		Id:   "con",
		Mode: "con",
		Requests: []*contract.RequestItem{
			{
				Id: "1",
			},
			{
				Id: "2",
			},
		},
	}, &TestContext{}, func(response *contract.Response) {
		assert.NotEmpty(t, response.Response)
		assert.Equal(t, response.Id, "con")
		assert.Len(t, response.Response, 2)
	})
}

func TestCascaded(t *testing.T) {
	Cascaded(&contract.Request{
		Id:   "cascaded",
		Mode: "cascaded",
		Requests: []*contract.RequestItem{
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
	}, &TestContext{}, func(response *contract.Response) {
		assert.NotEmpty(t, response.Response)
		assert.Equal(t, response.Id, "cascaded")
		assert.Len(t, response.Response, 2)
	})
}

func TestBatchWithBatchRequest(t *testing.T) {

	var batchName = "values"

	Batch(&contract.Request{
		Id: "batch",
		Seed: &contract.RequestItem{
			Id: "seed",
			Cascading: map[string]string{
				"values": "$.data",
			},
		},
		Mode: "batch",
		Requests: []*contract.RequestItem{
			{
				Id:    "batch",
				Batch: &batchName,
			},
		},
	}, &TestContext{}, func(response *contract.Response) {
		assert.NotEmpty(t, response.Response)
		assert.Equal(t, response.Id, "batch")
		assert.Len(t, response.Response, 3)
	})
}

func TestBatchWithData(t *testing.T) {

	var batchName = "values"

	Batch(&contract.Request{
		Id: "batch",
		Data: CascadeContext{
			"values": []string{"a", "b", "c"},
		},
		Mode: "batch",
		Requests: []*contract.RequestItem{
			{
				Id:    "batch",
				Batch: &batchName,
			},
		},
	}, &TestContext{}, func(response *contract.Response) {
		assert.NotEmpty(t, response.Response)
		assert.Equal(t, response.Id, "batch")
		assert.Len(t, response.Response, 3)
	})
}
