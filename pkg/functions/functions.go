package functions

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gustavo000/goLibGustavo/models"
)

func GenerateHttpResponse(status int, body any) *models.Response {
	switch s := body.(type) {
	case string:
		return &models.Response{Http: &http.Response{
			Status:     http.StatusText(status),
			Header:     make(http.Header),
			StatusCode: status, Body: io.NopCloser(bytes.NewReader([]byte(s))),
		}, Body: body}
	case []byte:
		return &models.Response{Http: &http.Response{
			Status:     http.StatusText(status),
			Header:     make(http.Header),
			StatusCode: status, Body: io.NopCloser(bytes.NewReader(s)),
		}, Body: body}
	default:
		marshal, err := json.Marshal(s)
		if err != nil {
			return &models.Response{Http: &http.Response{
				Status:     http.StatusText(500),
				Header:     make(http.Header),
				StatusCode: 500, Body: io.NopCloser(bytes.NewReader([]byte("response can't be parsed"))),
			}, Body: body}
		}
		return &models.Response{Http: &http.Response{
			Status:     http.StatusText(status),
			Header:     make(http.Header),
			StatusCode: status, Body: io.NopCloser(bytes.NewReader(marshal)),
		}, Body: body}
	}
}
