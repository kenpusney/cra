package core

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func formatResponse(response *http.Response, requestId string) *ResponseItem {
	resItem := &ResponseItem{Id: requestId}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		resItem.Error = UpstreamErrorResponse(response)
	}

	if strings.Contains(response.Header.Get("Content-Type"), "json") {
		resItem.Type = "json"
		resItem.Body = ConvertJsonBodyToObject(response.Body)
	} else {
		resItem.Type = "binary"
		resItem.Body = ReadBytes(response.Body)
	}
	if resItem.Error != nil {
		resItem.Error.Body = resItem.Body
	}
	return resItem
}

func decodeRequestBody(reqItem *RequestItem) io.Reader {
	var requestBody io.Reader = strings.NewReader("")

	if reqItem.Type == "json" {
		b, _ := json.Marshal(reqItem.Body)
		requestBody = bytes.NewReader(b)
	} else {
		if reqItem.Body != nil {
			requestBody = base64.NewDecoder(base64.RawStdEncoding, strings.NewReader(reqItem.Body.(string)))
		}
	}
	return requestBody
}

func fillRequests(req *RequestItem) {
	if req.Method == "" {
		req.Method = "GET"
	}
}
