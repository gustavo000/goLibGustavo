package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Response struct {
	Http  *http.Response
	Mutex sync.Mutex
	Body  any
}

type SpanInfo struct {
	Span         trace.Span      `json:"-"`
	SpanID       trace.SpanID    `json:"spanID"`
	FunctionName string          `json:"functionName"`
	CreatedAt    time.Time       `json:"timestamp"`
	Context      context.Context `json:"-"`
	TraceID      trace.TraceID   `json:"traceID"`
	ParentSpanID trace.SpanID    `json:"parentSpanID"`
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

func (r *Response) SetBody(object any) error {
	marshal, err := json.Marshal(object)
	if err != nil {
		return err
	}
	r.GetHttp().Body = io.NopCloser(bytes.NewReader(marshal))
	return nil
}

func (r *Response) GetBodyBytes() ([]byte, error) {
	if r.GetHttp().Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, r.GetHttp().Body)
	if err != nil {
		return nil, err
	}

	responseBytes := buffer.Bytes()
	r.GetHttp().Body = io.NopCloser(bytes.NewReader(responseBytes))
	return responseBytes, nil
}

func (r *Response) GetBodyString() string {
	responseBytes, err := r.GetBodyBytes()
	if err != nil {
		return ""
	}
	return string(responseBytes)
}

func (r *Response) GetObject(v any) error {
	responseBytes, err := r.GetBodyBytes()
	if err != nil {
		return err
	}
	return unmarshalToObject(responseBytes, v)
}

func (r *Response) IsException() bool {
	if r.GetHttp() == nil {
		return true
	}
	return r.GetHttp().StatusCode >= 400 && r.GetHttp().StatusCode <= 599
}

func (r *Response) GetHeader() http.Header {
	return r.GetHttp().Header
}

func (r *Response) SetHeader(headers http.Header) {
	r.GetHttp().Header = headers
}

func (r *Response) GetHeaderClone() (header http.Header, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	if httpHeader := r.GetHeader(); httpHeader != nil {
		header = httpHeader.Clone()
	}
	return header, err
}

func (r *Response) CopyBody() io.ReadCloser {
	responseBytes, err := r.GetBodyBytes()
	if err != nil {
		return io.NopCloser(bytes.NewReader([]byte(err.Error())))
	}
	return io.NopCloser(bytes.NewReader(responseBytes))
}

func unmarshalToObject(toDecode []byte, v any) error {
	defer handlerPanic()
	if json.Valid(toDecode) {
		bufferOfBytes := &bytes.Buffer{}
		json.HTMLEscape(bufferOfBytes, toDecode)
		if errUnmarshal := json.Unmarshal(bufferOfBytes.Bytes(), v); errUnmarshal != nil {
			return errUnmarshal
		}
		return nil
	} else {
		return fmt.Errorf("json data is not valid")
	}
}

func handlerPanic() {
	if r := recover(); r != nil {

	}
}
