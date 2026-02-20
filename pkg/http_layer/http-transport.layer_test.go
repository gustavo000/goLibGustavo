package http_layer

import (
	c "context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stretchr/testify/assert"
	constants_headers "github.com/gustavo000/goLibGustavo/constants/constants_headers"
	headers_request "github.com/gustavo000/goLibGustavo/pkg/headers_request"
	"github.com/gustavo000/goLibGustavo/resources/properties"
)

func TestGetHttpClient(t *testing.T) {
	propertiesService := properties.Service{
		Name:       "test",
		Timeout:    1000,
		Client:     nil,
		IsExternal: true,
	}

	httpClient := GetHttpClient(&propertiesService, "http://test.com")
	assert.NotNil(t, httpClient)
	assert.Equal(t, &HttpClient{
		Url:               "http://test.com",
		ExternalService:   "test",
		Timeout:           int64(1000),
		Client:            nil,
		IsExternalService: true,
	}, httpClient)
}

func TestHttpClient_WithContext(t *testing.T) {
	httpClient := &HttpClient{}
	app := iris.New()
	ctx := context.NewContext(app)
	var ok bool
	httpClient, ok = httpClient.WithContext(ctx).(*HttpClient)
	assert.True(t, ok)
	assert.NotNil(t, httpClient)
	assert.Equal(t, ctx, httpClient.Context)
	assert.Equal(t, ctx.Request(), httpClient.HttpRequest)
}

func TestHttpClient_WithHeaders(t *testing.T) {
	httpClient := &HttpClient{}
	headers := headers_request.RequestHeaders{}
	headers.AddHeaders(headers_request.RequestHeader{Key: "test"})
	var ok bool
	httpClient, ok = httpClient.WithHeaders(headers).(*HttpClient)
	assert.True(t, ok)
	assert.NotNil(t, httpClient)
	assert.Equal(t, headers, httpClient.Headers)
}

func TestHttpClient_WithBody(t *testing.T) {
	httpClient := &HttpClient{}
	body := []byte("test")
	var ok bool
	httpClient, ok = httpClient.WithBody(body).(*HttpClient)
	assert.True(t, ok)
	assert.NotNil(t, httpClient)
	assert.Equal(t, body, httpClient.RequestBody)
}

func TestHttpClient_SetRequestHeaders(t *testing.T) {
	properties.NewProperties(
		properties.WithEnvironment("BETA"),
	)

	testTable := []struct {
		name             string
		newHeader        http.Header
		request          *http.Request
		headers          headers_request.RequestHeaders
		httpHeaderExpect http.Header
	}{
		{
			name:      "Test case 1",
			newHeader: http.Header{},
			httpHeaderExpect: map[string][]string{
				constants_headers.EXT_SERVICE:         {""},
				constants_headers.EXT_REQUEST_TIMEOUT: {"0s"},
				"Ext-Test":                            {"ForceValueTest"},
			},
			headers: headers_request.RequestHeaders{
				{
					Key:        "Ext-Test",
					ForceValue: "ForceValueTest",
				},
			},
		},
		{
			name:      "Test case 2",
			newHeader: http.Header{},
			httpHeaderExpect: map[string][]string{
				constants_headers.EXT_SERVICE:         {""},
				constants_headers.EXT_REQUEST_TIMEOUT: {"0s"},
			},
			headers: headers_request.RequestHeaders{
				{
					Key: "Ext-Test",
				},
			},
		},
		{
			name:      "Test case 3",
			newHeader: http.Header{},
			httpHeaderExpect: map[string][]string{
				constants_headers.EXT_SERVICE:         {""},
				constants_headers.EXT_REQUEST_TIMEOUT: {"0s"},
			},
			headers: headers_request.RequestHeaders{
				{
					Key:       "Ext-Test",
					Translate: "TranslateTest",
				},
			},
		},
		{
			name:      "Test case 4",
			newHeader: http.Header{},
			httpHeaderExpect: map[string][]string{
				constants_headers.EXT_SERVICE:         {""},
				constants_headers.EXT_REQUEST_TIMEOUT: {"0s"},
				"EnvTest":                             {"PROD"},
			},
			headers: headers_request.RequestHeaders{
				{
					Key:       "X-Environment",
					Translate: "EnvTest",
				},
			},
		},
		{
			name:      "Test case 5",
			newHeader: http.Header{},
			httpHeaderExpect: map[string][]string{
				constants_headers.EXT_SERVICE:         {""},
				constants_headers.EXT_REQUEST_TIMEOUT: {"0s"},
			},
		},
	}

	httpClient := &HttpClient{}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			httpHeader := httpClient.SetRequestHeaders(tt.newHeader, tt.request, tt.headers)
			u := httpHeader.Get(constants_headers.EXT_REQUEST_ID)
			_, err := uuid.Parse(u)
			assert.Nil(t, err)

			httpHeader.Del(constants_headers.EXT_REQUEST_ID)
			assert.Equal(t, tt.httpHeaderExpect, httpHeader)
		})
	}
}

func getMockAppExternalService() *iris.Application {
	appMock := iris.New()
	appMock.Get("/", func(ctx iris.Context) {
		ctx.JSON(map[string]interface{}{
			"Id":   1,
			"Name": "test",
			"Values": []string{
				"test",
			},
		})
	})

	appMock.Post("/", func(ctx iris.Context) {
		ctx.JSON(map[string]interface{}{
			"Id":   1,
			"Name": "test",
			"Values": []string{
				"test",
			},
		})
	})

	return appMock
}

func TestHttpClient_PerformHttpRequest(t *testing.T) {
	properties.NewProperties(
		properties.WithVersion("1.0.0"),
		properties.WithServiceName("test"),
	)

	appMock := getMockAppExternalService()
	go func() {
		err := appMock.Listen(":50020")
		assert.Equal(t, "http: Server closed", err.Error())
	}()

	defer func() {
		err := appMock.Shutdown(c.Background())
		assert.NoError(t, err)
	}()

	for {
		_, err := http.Get("http://localhost:50020")
		if err != nil && strings.Contains(err.Error(), "connection refused") {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}

	testTable := []struct {
		name             string
		client           *HttpClient
		method           string
		expectedResponse struct {
			Status     string
			StatusCode int
			Body       []byte
		}
	}{
		{
			name: "GET response",
			client: &HttpClient{
				Url:    "http://localhost:50020",
				Client: &http.Client{},
			},
			method: "GET",
			expectedResponse: struct {
				Status     string
				StatusCode int
				Body       []byte
			}{
				Status:     "200 OK",
				StatusCode: 200,
				Body:       []byte(`{"Id":1,"Name":"test","Values":["test"]}` + "\n"),
			},
		},
		{
			name: "GET response with is external service",
			client: &HttpClient{
				Url:               "http://localhost:50020",
				Client:            &http.Client{},
				IsExternalService: true,
			},
			method: "GET",
			expectedResponse: struct {
				Status     string
				StatusCode int
				Body       []byte
			}{
				Status:     "200 OK",
				StatusCode: 200,
				Body:       []byte(`{"Id":1,"Name":"test","Values":["test"]}` + "\n"),
			},
		},
		{
			name: "POST response",
			client: &HttpClient{
				Url:    "http://localhost:50020",
				Client: &http.Client{},
			},
			method: "POST",
			expectedResponse: struct {
				Status     string
				StatusCode int
				Body       []byte
			}{
				Status:     "200 OK",
				StatusCode: 200,
				Body:       []byte(`{"Id":1,"Name":"test","Values":["test"]}` + "\n"),
			},
		},
		{
			name: "UNKOWN response",
			client: &HttpClient{
				Url:    "http://localhost:50020",
				Client: &http.Client{},
			},
			method: "UNKOWN",
			expectedResponse: struct {
				Status     string
				StatusCode int
				Body       []byte
			}{
				Status:     "Method Not Allowed",
				StatusCode: 405,
				Body:       []byte("Method UNKOWN not allowed"),
			},
		},
		{
			name:   "return error",
			client: &HttpClient{},
			expectedResponse: struct {
				Status     string
				StatusCode int
				Body       []byte
			}{
				Status:     "Internal Server Error",
				StatusCode: 500,
				Body:       []byte("client http not exists"),
			},
		},
		{
			name: "GET response error",
			client: &HttpClient{
				Url:    "http://localhost:8081",
				Client: &http.Client{},
			},
			method: "GET",
			expectedResponse: struct {
				Status     string
				StatusCode int
				Body       []byte
			}{
				Status:     "Internal Server Error",
				StatusCode: 500,
				Body:       []byte("Get \"http://localhost:8081\": dial tcp 127.0.0.1:8081: connect: connection refused"),
			},
		},
		{
			name: "GET response error with is external service",
			client: &HttpClient{
				Url:               "http://localhost:8081",
				Client:            &http.Client{},
				IsExternalService: true,
			},
			method: "GET",
			expectedResponse: struct {
				Status     string
				StatusCode int
				Body       []byte
			}{
				Status:     "Internal Server Error",
				StatusCode: 500,
				Body:       []byte("Get \"http://localhost:8081\": dial tcp 127.0.0.1:8081: connect: connection refused"),
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			modelResponse := tt.client.PerformHttpRequest(tt.method)
			assert.Equal(t, tt.expectedResponse.Status, modelResponse.Http.Status)
			assert.Equal(t, tt.expectedResponse.StatusCode, modelResponse.Http.StatusCode)
			body, err := io.ReadAll(modelResponse.Http.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(tt.expectedResponse.Body), string(body))
		})
	}
}
