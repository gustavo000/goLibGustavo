package models

import (
	"io"
	"net/http"
	"strings"
	"sync"
)

type Response struct {
	Http  *http.Response
	Mutex sync.Mutex
	Body  any
}

func (r *Response) GetHttp() *http.Response {
	var response *http.Response
	response = r.Http
	if response == nil {
		return &http.Response{StatusCode: http.StatusInternalServerError,
			Status: http.StatusText(http.StatusInternalServerError),
			Body:   io.NopCloser(strings.NewReader("can't get response from external service")),
		}
	}
	return response
}
