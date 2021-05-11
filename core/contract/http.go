package contract

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/kenpusney/cra/core/util"
	"io"
	"net/http"
	"strings"
)

func ErrorResponse(req *RequestItem, err error, response *http.Response) *ResponseItem {
	return &ResponseItem{Id: req.Id, Error: SystemErrorResponse(err), r: response}
}

func FormatResponse(response *http.Response, requestId string) *ResponseItem {
	resItem := &ResponseItem{Id: requestId, r: response}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		resItem.Error = UpstreamErrorResponse(response)
	}

	if strings.Contains(response.Header.Get("Content-Type"), "json") {
		resItem.Type = "json"
		resItem.Body = util.ConvertJsonBodyToObject(response.Body)
	} else {
		resItem.Type = "binary"
		resItem.Body = util.ReadBytes(response.Body)
	}
	if resItem.Error != nil {
		resItem.Error.Body = resItem.Body
	}
	return resItem
}

func DecodeRequestBody(reqItem *RequestItem) io.Reader {
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

func FillRequest(req *RequestItem) {
	if req.Method == "" {
		req.Method = "GET"
	}
}
