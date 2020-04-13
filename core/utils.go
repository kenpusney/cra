package core

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func ConvertJsonBodyToObject(body io.ReadCloser) interface{} {
	bytes := ConvertJsonBodyToBytes(body)
	if bytes != nil {
		var result interface{}
		json.Unmarshal(bytes, &result)
		return result
	}
	return nil
}

func ConvertJsonBodyToBytes(body io.ReadCloser) []byte {
	all, err := ioutil.ReadAll(body)
	if err != nil {
		return nil
	}
	return all
}

func UnmarshallJsonObject(bytes []byte, object interface{}) interface{} {
	if bytes != nil {
		json.Unmarshal(bytes, object)
	}
	return object
}
