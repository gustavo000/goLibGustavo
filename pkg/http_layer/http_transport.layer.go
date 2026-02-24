package http_layer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gustavo000/goLibGustavo/models/properties"
	"github.com/gustavo000/goLibGustavo/models/rest"
	"github.com/gustavo000/goLibGustavo/pkg/logger"

	"github.com/google/uuid"
	"github.com/gustavo000/goLibGustavo/constants/constants_headers"
	"github.com/gustavo000/goLibGustavo/pkg/functions"
	"github.com/gustavo000/goLibGustavo/pkg/headers_request"
)

type IHttpTransport interface {
	WithContext(ctx *http.Request) IHttpTransport
	WithBody(body []byte) IHttpTransport
	WithHeaders(headers headers_request.RequestHeaders) IHttpTransport
	WithParentSpan(span *rest.SpanInfo) IHttpTransport
	SetRequestHeaders(newHeader http.Header, request *http.Request, headers headers_request.RequestHeaders) http.Header
	MakeGet() *rest.Response
	MakePost() *rest.Response
	MakePut() *rest.Response
	MakePatch() *rest.Response
	MakeDelete() *rest.Response
	PerformHttpRequest(method string) *rest.Response
	GetBody() (io.ReadCloser, error)
}

type HttpClient struct {
	Url               string
	ExternalService   string
	RequestBody       []byte
	Timeout           int64
	Context           context.Context
	HttpRequest       *http.Request
	Headers           headers_request.RequestHeaders
	Client            *http.Client
	IsExternalService bool
	ParentSpan        *rest.SpanInfo
}

func GetHttpClient(service *rest.Service, url string) *HttpClient {
	return &HttpClient{
		Url:               url,
		ExternalService:   service.Name,
		Timeout:           service.Timeout,
		Client:            service.Client,
		IsExternalService: service.IsExternal,
	}
}

func (h *HttpClient) WithContext(ctx *http.Request) IHttpTransport {
	if ctx != nil {
		h.HttpRequest = ctx
		h.Context = ctx.Context()
	}
	return h
}

func (h *HttpClient) WithParentSpan(span *rest.SpanInfo) IHttpTransport {
	if span != nil {
		h.ParentSpan = span
	}
	return h
}

func (h *HttpClient) WithHeaders(headers headers_request.RequestHeaders) IHttpTransport {
	h.Headers = headers
	return h
}
func (h *HttpClient) WithBody(body []byte) IHttpTransport {
	h.RequestBody = body
	return h
}

// SetRequestHeaders It takes a request header, a list of headers-validator to set, and returns a new header with the headers-validator set
func (h *HttpClient) SetRequestHeaders(newHeader http.Header, request *http.Request, headers headers_request.RequestHeaders) http.Header {
	var requestHeader = make(http.Header)
	if request != nil {
		requestHeader = request.Header
	}
	for _, v := range headers {
		var key string
		var value string
		value = requestHeader.Get(v.Key)
		if v.ForceValue != "" {
			value = v.ForceValue
		} else if requestHeader.Get(v.Key) == "" {
			value = v.DefaultValue
		}
		if v.Translate != "" {
			key = v.Translate
		} else {
			key = v.Key
		}

		if strings.Contains(v.Key, "X-Environment") &&
			v.ForceValue == "" &&
			properties.GetProperty().GetEnv() == "BETA" &&
			!strings.Contains(h.Url, "cust-rtmn-orch") &&
			!strings.Contains(h.Url, "policyorch") &&
			!strings.Contains(h.Url, "api-qa-ftc-sc.falabella.com") {
			value = "PROD"
		}
		if value != "" {
			newHeader[key] = []string{value}
		}
	}
	requestId := uuid.NewString()
	newHeader.Set(constants_headers.EXT_REQUEST_ID, requestId)
	newHeader.Set(constants_headers.EXT_SERVICE, h.ExternalService)
	timeout := time.Second * time.Duration(h.Timeout)
	newHeader.Set(constants_headers.EXT_REQUEST_TIMEOUT, timeout.String())
	return newHeader
}

// MakeGet function to perform any Http HttpRequest under the method GET
func (h *HttpClient) MakeGet() *rest.Response {
	return h.PerformHttpRequest("GET")
}

// MakePost function to perform any Http HttpRequest under the method POST
func (h *HttpClient) MakePost() *rest.Response {
	return h.PerformHttpRequest("POST")
}

// MakePut function to perform any Http HttpRequest under the method POST
func (h *HttpClient) MakePut() *rest.Response {
	return h.PerformHttpRequest("PUT")
}

// MakePatch function to perform any Http HttpRequest under the method PATCH
func (h *HttpClient) MakePatch() *rest.Response {
	return h.PerformHttpRequest("PATCH")
}

// MakeDelete function to perform any Http HttpRequest under the method DELETE
func (h *HttpClient) MakeDelete() *rest.Response {
	return h.PerformHttpRequest("DELETE")
}

func (h *HttpClient) GetBody() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewBuffer(h.RequestBody)), nil
}

func (h *HttpClient) PerformHttpRequest(method string) *rest.Response {

	if h.Client == nil {
		return functions.GenerateHttpResponse(500, "client http not exists")
	}
	var req *http.Request
	switch method {
	case "GET", "DELETE":
		req, _ = http.NewRequest(method, h.Url, nil)
	case "PUT", "POST", "PATCH":
		req, _ = http.NewRequest(method, h.Url, io.NopCloser(bytes.NewBuffer(h.RequestBody)))
	default:
		return functions.GenerateHttpResponse(http.StatusMethodNotAllowed, fmt.Sprintf("Method %s not allowed", method))
	}
	req.GetBody = h.GetBody
	req.Header = h.SetRequestHeaders(req.Header, h.HttpRequest, h.Headers)

	if req.Context() != nil && h.Context != nil {
		logger.Info("HOLA")
	}

	if h.IsExternalService {
		logger.Info("HOLA")
	}
	resp, errClient := h.Client.Do(req)
	if errClient != nil {
		if h.IsExternalService {
			logger.Info("HOLA")
		}

		if strings.Contains(errClient.Error(), "Client.Timeout") {
			return functions.GenerateHttpResponse(http.StatusRequestTimeout, errClient.Error())
		}
		return functions.GenerateHttpResponse(http.StatusInternalServerError, errClient.Error())
	}
	if resp == nil {
		if h.IsExternalService {
			logger.Info("HOLA")
		}
		return functions.GenerateHttpResponse(http.StatusInternalServerError, "can't get response from URL "+h.Url)
	}
	resp.Header.Set(constants_headers.EXT_REQUEST_ID, req.Header.Get("Ext-Request-Id"))
	if h.IsExternalService {
		logger.Info("HOLA")
	}
	return &rest.Response{Http: resp}
}
