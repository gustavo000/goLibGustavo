package url_helper

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	http_layer "github.com/gustavo000/goLibGustavo/pkg/http_layer"
	"github.com/gustavo000/goLibGustavo/models"
	"github.com/gustavo000/goLibGustavo/resources/properties"
)

func Test_GetClientBy(t *testing.T) {
	os.Setenv("ENVIRONMENT", "BETA")

	testTable := []struct {
		name       string
		layer      string
		url        string
		path       string
		forceAtEnd string
		IsExternal bool
	}{
		{
			layer:      "EXTERNAL",
			name:       "without force at end",
			url:        "http://localhost:8080/mock",
			path:       "/dl-core",
			IsExternal: true,
		},
		{
			layer:      "DL",
			name:       "with force at end",
			url:        "http://localhost:8080/mock",
			path:       "/dl-core",
			forceAtEnd: "/abc",
		},
	}

	properties.NewProperties(
		properties.WithEnvironment("BETA"),
		properties.WithServices([]*properties.Service{
			{
				Layer:   "EXTERNAL",
				Name:    "Core",
				Ingress: "http://localhost:8080/mock",
				Path:    "/dl-core",
				Endpoints: []properties.Endpoint{
					{
						Name: "GetConstantByName",
					},
				},
			},
			{
				Layer:   "DL",
				Name:    "Core",
				Ingress: "http://localhost:8080/mock",
				Path:    "/dl-core",
				Endpoints: []properties.Endpoint{
					{
						Name: "GetConstantByName",
					},
				},
			},
		}...),
	)

	service := &properties.Service{
		Name:    "Core",
		Ingress: "http://localhost:8080/mock",
		Path:    "/dl-core",
		Endpoints: []properties.Endpoint{
			{
				Name: "GetConstantByName",
			},
		},
		Timeout: 10,
		Client:  &http.Client{Timeout: time.Second * time.Duration(10)},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			helperUrl := models.Url{
				LayerName:    tt.layer,
				ServiceName:  "Core",
				EndpointName: "GetConstantByName",
				ForceAtEnd:   tt.forceAtEnd,
			}

			client := GetClientBy(helperUrl)
			service.Layer = tt.layer
			service.IsExternal = tt.IsExternal
			expectedClient := http_layer.GetHttpClient(service, tt.url+tt.path+tt.forceAtEnd)
			assert.NotNil(t, client)
			assert.Equal(t, expectedClient, client)
		})
	}
}

func Test_SetIngress(t *testing.T) {
	testTable := []struct {
		name         string
		url          string
		ingressForce string
		layer        string
		release      string
		namespace    string
		isLocal      bool
		ingress      string
		expectedUrl  string
	}{
		{
			name:         "with ingress force",
			url:          "http://localhost:8080/mock",
			ingressForce: "/force",
			expectedUrl:  "http://localhost:8080/mock/force",
		},
		{
			name:        "with service ingress",
			url:         "http://localhost:8080/mock",
			ingress:     "/ingress",
			expectedUrl: "http://localhost:8080/mock/ingress",
		},
		{
			name:        "with service ingress and release and namespace",
			url:         "http://localhost:8080/{release}/{namespace}/mock",
			ingress:     "/ingress",
			expectedUrl: "http://localhost:8080/v1/namespace/mock/ingress",
			release:     "v1",
			namespace:   "namespace",
		},
		{
			name:        "with get url by env",
			url:         "",
			isLocal:     true,
			expectedUrl: "https://local-policyorch.ftc-cx.tech/",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.isLocal {
				os.Setenv("ENVIRONMENT", "LOCAL")
			} else {
				os.Unsetenv("ENVIRONMENT")
			}

			properties.NewProperties(
				properties.WithEnvironment("LOCAL"),
				properties.WithRelease(tt.release),
				properties.WithNameSpace(tt.namespace),
			)

			helperUrl := models.Url{
				IngressForce: tt.ingressForce,
			}

			service := &properties.Service{
				Ingress: tt.ingress,
			}
			url := SetIngress(tt.url, helperUrl, service)
			assert.Equal(t, tt.expectedUrl, url)
		})
	}
}

func Test_SetEndpoint(t *testing.T) {
	testTable := []struct {
		name          string
		url           string
		endpointForce string
		endpoints     []properties.Endpoint
		expectedUrl   string
	}{
		{
			name:          "with endpoint force",
			url:           "http://localhost:8080/mock",
			endpointForce: "/force",
			expectedUrl:   "http://localhost:8080/mock/force",
		},
		{
			name: "with service endpoint",
			url:  "http://localhost:8080/mock",
			endpoints: []properties.Endpoint{
				{
					Name: "GetConstantByName",
					Path: "/get-constant",
				},
			},
			expectedUrl: "http://localhost:8080/mock/get-constant",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			helperUrl := models.Url{
				EndpointForce: tt.endpointForce,
				EndpointName:  "GetConstantByName",
			}
			service := &properties.Service{
				Endpoints: tt.endpoints,
			}
			url := SetEndpoint(tt.url, helperUrl, service)
			assert.Equal(t, tt.expectedUrl, url)
		})
	}
}

func Test_SetPathParams(t *testing.T) {
	testTable := []struct {
		name         string
		url          string
		parserParams []models.ParserParam
		expected     string
	}{
		{
			name: "with parser params",
			url:  "http://test/{id}",
			parserParams: []models.ParserParam{
				{
					Key:   "{id}",
					Value: "1",
				},
			},
			expected: "http://test/1",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			helperUrl := models.Url{
				ParserParams: tt.parserParams,
			}

			url := SetPathParams(tt.url, helperUrl)
			assert.Equal(t, tt.expected, url)
		})
	}
}

func Test_SetQueryParams(t *testing.T) {
	testTable := []struct {
		name        string
		queryForce  string
		queryParams []models.QueryParam
		expected    string
	}{

		{

			name:        "with query force",
			queryForce:  "force",
			queryParams: []models.QueryParam{},
			expected:    "http://test?force",
		},
		{
			name: "with query params",
			queryParams: []models.QueryParam{
				{
					Key:   "key",
					Value: "value",
				},
			},
			expected: "http://test?key=value",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			helperUrl := models.Url{
				QueryForce:  tt.queryForce,
				QueryParams: tt.queryParams,
			}

			url := SetQueryParams("http://test", helperUrl)
			assert.Equal(t, tt.expected, url)
		})
	}
}

func Test_GetUrlByEnv(t *testing.T) {
	testTable := []struct {
		name     string
		env      string
		expected string
	}{
		{
			name:     "with test env",
			env:      "TEST",
			expected: "https://uat-policyorch.ftc-cx.tech/{release}",
		},
		{
			name:     "with prod env",
			env:      "PROD",
			expected: "https://prd-policyorch.ftc-cx.tech/{release}",
		},
		{
			name:     "with dev env",
			env:      "DEV",
			expected: "https://dev-policyorch.ftc-cx.tech/{release}",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			properties.NewProperties(
				properties.WithEnvironment(tt.env),
			)
			url := GetUrlByEnv()
			assert.Equal(t, tt.expected, url)
		})
	}
}

func Test_GetServiceByLayerAndName2(t *testing.T) {
	testTable := []struct {
		name            string
		layer           string
		serviceName     string
		endpoints       []properties.Endpoint
		expectedService *properties.Service
	}{
		{
			name:        "case 1",
			layer:       "DL",
			serviceName: "Core",
			endpoints: []properties.Endpoint{
				{
					Name: "GetConstantByName",
				},
			},
			expectedService: &properties.Service{
				Layer:   "DL",
				Name:    "Core",
				Ingress: "http://localhost:8080/mock",
				Path:    "/dl-core",
				Endpoints: []properties.Endpoint{
					{
						Name: "GetConstantByName",
					},
				},
				Timeout: 10,
				Client:  &http.Client{Timeout: time.Second * time.Duration(10)},
			},
		},
		{
			name:        "case 2",
			layer:       "DL",
			serviceName: "Core",
			endpoints:   []properties.Endpoint{},
			expectedService: &properties.Service{
				Layer:     "DL",
				Name:      "Core",
				Ingress:   "http://localhost:8080/mock",
				Path:      "/dl-core",
				Endpoints: []properties.Endpoint{},
				Timeout:   10,
				Client:    &http.Client{Timeout: time.Second * time.Duration(10)},
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			properties.NewProperties(
				properties.WithEnvironment("BETA"),
				properties.WithServices([]*properties.Service{
					{
						Layer:     "DL",
						Name:      "Core",
						Ingress:   "http://localhost:8080/mock",
						Path:      "/dl-core",
						Endpoints: tt.endpoints,
					},
				}...),
			)

			svc := GetServiceByLayerAndName(tt.layer, tt.serviceName)
			assert.NotNil(t, svc)
			assert.Equal(t, tt.expectedService, svc)
		})
	}
}
