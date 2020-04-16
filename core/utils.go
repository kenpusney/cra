package core

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func ConvertJsonBodyToObject(body io.Reader) interface{} {
	bytes := ReadBytes(body)
	var object interface{}
	UnmarshallJsonObject(bytes, &object)
	if object == nil {
		return nil
	}
	return object
}

func ReadBytes(body io.Reader) []byte {
	all, err := ioutil.ReadAll(body)
	if err != nil {
		return nil
	}
	return all
}

func UnmarshallJsonObject(bytes []byte, object interface{}) interface{} {
	if bytes != nil {
		_ = json.Unmarshal(bytes, object)
		return object
	}
	return nil
}

func GenerateId(requestId string, uuid string, index int) string {
	if len(requestId) == 0 {
		requestId = uuid
	}
	if index >= 0 {
		return fmt.Sprintf("%s-%d", requestId, index)
	}
	return requestId
}
