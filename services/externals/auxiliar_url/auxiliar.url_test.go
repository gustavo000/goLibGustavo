package auxiliar_url

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

	properties.NewProperties(
		properties.WithEnvironment("BETA"),
		properties.WithServices([]*properties.Service{
			{
				Layer:   "Sec",
				Name:    "CorpEnvironment",
				Ingress: "http://localhost:8080/mock",
			},
		}...),
	)

	testTable := []struct {
		name               string
		layer              string
		serviceName        string
		forceAtEnd         string
		httpClientExpected *http_layer.HttpClient
	}{
		{
			name:        "test with external",
			layer:       "Sec",
			serviceName: "CorpEnvironment",
			httpClientExpected: http_layer.GetHttpClient(&properties.Service{
				Layer:   "Sec",
				Name:    "CorpEnvironment",
				Ingress: "http://localhost:8080/mock",
				Timeout: 10,
				Client:  &http.Client{Timeout: time.Second * time.Duration(10)},
			}, "http://localhost:8080/mock"),
		},
		{
			name:        "test with external",
			layer:       "Sec",
			serviceName: "CorpEnvironment",
			forceAtEnd:  "/test",
			httpClientExpected: http_layer.GetHttpClient(&properties.Service{
				Layer:   "Sec",
				Name:    "CorpEnvironment",
				Ingress: "http://localhost:8080/mock",
				Timeout: 10,
				Client:  &http.Client{Timeout: time.Second * time.Duration(10)},
			}, "http://localhost:8080/mock/test"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			client := GetClientBy(models.Url{
				LayerName:   tt.layer,
				ServiceName: tt.serviceName,
				ForceAtEnd:  tt.forceAtEnd,
			})
			assert.Equal(t, tt.httpClientExpected, client)
		})
	}
}

func Test_SetIngress(t *testing.T) {
	testTable := []struct {
		name        string
		url         string
		isLocal     bool
		layer       string
		expectedUrl string
	}{
		{
			name:        "test with local",
			isLocal:     true,
			expectedUrl: "https://local-policyorch.ftc-cx.tech/v1",
		},
		{
			name:        "test with local",
			url:         "http://localhost:8080/{release}/{namespace}",
			layer:       "test",
			expectedUrl: "http://localhost:8080/v1/test/ingress",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.isLocal {
				os.Setenv("ENVIRONMENT", "TEST")
			} else {
				os.Unsetenv("ENVIRONMENT")
			}

			properties.NewProperties(
				properties.WithEnvironment("LOCAL"),
				properties.WithRelease("v1"),
				properties.WithNameSpace("test"),
			)

			url := SetIngress(tt.url,
				models.Url{},
				&properties.Service{
					Layer:   tt.layer,
					Ingress: "/ingress",
				})

			assert.Equal(t, tt.expectedUrl, url)
		})
	}
}

func Test_SetEndpoint(t *testing.T) {
	url := SetEndpoint("http://localhost:8080/mock", models.Url{EndpointForce: "/force"}, nil)
	assert.Equal(t, "http://localhost:8080/mock/force", url)
}

func Test_SetPathParams(t *testing.T) {
	expectedUrl := "https://local/v1"
	url := SetPathParams("https://local/{release}", models.Url{ParserParams: []models.ParserParam{{Key: "{release}", Value: "v1"}}})
	assert.Equal(t, expectedUrl, url)
}

func Test_SetQueryParams(t *testing.T) {
	tableTest := []struct {
		name        string
		url         string
		queryForce  string
		queryParams []models.QueryParam
		expectedUrl string
	}{
		{
			name:        "test with query force",
			url:         "https://uat-policyorch.ftc-cx.tech/{release}",
			queryForce:  "force=true",
			expectedUrl: "https://uat-policyorch.ftc-cx.tech/{release}?force=true",
		},
		{
			name:        "test with query params",
			url:         "https://uat-policyorch.ftc-cx.tech/{release}",
			queryParams: []models.QueryParam{{Key: "force", Value: "true"}},
			expectedUrl: "https://uat-policyorch.ftc-cx.tech/{release}?force=true",
		},
	}

	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedUrl, SetQueryParams(tt.url, models.Url{QueryForce: tt.queryForce, QueryParams: tt.queryParams}))
		})
	}
}

func Test_GetUrlByEnv(t *testing.T) {
	os.Unsetenv("ENVIRONMENT")
	tableTest := []struct {
		name        string
		env         string
		expectedUrl string
	}{
		{
			name:        "test with env dev",
			env:         "TEST",
			expectedUrl: "https://uat-policyorch.ftc-cx.tech/{release}",
		},
		{
			name:        "test with env prod",
			env:         "PROD",
			expectedUrl: "https://prd-policyorch.ftc-cx.tech/{release}",
		},
		{
			name:        "test with empty env",
			env:         "",
			expectedUrl: "https://-policyorch.ftc-cx.tech/{release}",
		},
	}

	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			properties.NewProperties(
				properties.WithEnvironment(tt.env),
			)

			assert.Equal(t, tt.expectedUrl, GetUrlByEnv())
		})
	}
}

func Test_GetServiceByLayerAndName(t *testing.T) {
	os.Setenv("ENVIRONMENT", "BETA")

	properties.NewProperties(
		properties.WithEnvironment("BETA"),
		properties.WithServices([]*properties.Service{
			{
				Layer:   "Sec",
				Name:    "CorpEnvironment",
				Ingress: "http://localhost:8080/mock",
			},
		}...),
	)

	serviceResiltExpected := &properties.Service{
		Layer:   "Sec",
		Name:    "CorpEnvironment",
		Ingress: "http://localhost:8080/mock",
		Timeout: 10,
		Client:  &http.Client{Timeout: time.Second * time.Duration(10)},
	}

	serviceResilt := GetServiceByLayerAndName("Sec", "CorpEnvironment")
	assert.Equal(t, serviceResiltExpected, serviceResilt)
}
