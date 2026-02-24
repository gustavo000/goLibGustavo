package auxiliar_url

/*
import (
	"fmt"
	"github.com/gustavo000/goLibGustavo/models"
	http_layer "github.com/gustavo000/goLibGustavo/pkg/http_layer"
	"github.com/gustavo000/goLibGustavo/resources/properties"
	"net/http"
	"strings"
	"time"
)

func GetClientBy(helperUrl models.Url) *http_layer.HttpClient {
	var url string
	service := GetServiceByLayerAndName(helperUrl.LayerName, helperUrl.ServiceName)
	url = SetIngress(url, helperUrl, service)
	url += service.Path
	url = SetEndpoint(url, helperUrl, service)
	url = SetPathParams(url, helperUrl)
	url = SetQueryParams(url, helperUrl)

	if helperUrl.ForceAtEnd != "" {
		url += helperUrl.ForceAtEnd
	}
	return http_layer.GetHttpClient(service, url)
}

func SetIngress(url string, helperUrl models.Url, service *properties.Service) string {
	if properties.GetProperty().IsLocal() && !strings.Contains(service.Layer, "EXTERNAL") {
		url += GetUrlByEnv()
	} else {
		url += service.Ingress
	}
	url = strings.ReplaceAll(url, "{release}", properties.GetProperty().GetRelease())
	url = strings.ReplaceAll(url, "{namespace}", properties.GetProperty().GetNamespace())
	return url
}

func SetEndpoint(url string, helperUrl models.Url, service *properties.Service) string {
	url += helperUrl.EndpointForce
	return url
}

func SetPathParams(url string, helperUrl models.Url) string {
	if len(helperUrl.ParserParams) > 0 {
		for _, param := range helperUrl.ParserParams {
			url = strings.ReplaceAll(url, param.Key, param.Value)
		}
	}
	return url
}

func SetQueryParams(url string, helperUrl models.Url) string {
	if helperUrl.QueryForce != "" {
		url += "?" + helperUrl.QueryForce
	} else if len(helperUrl.QueryParams) > 0 {
		url += "?"
		for _, query := range helperUrl.QueryParams {
			url += fmt.Sprintf("%s=%s&", query.Key, query.Value)
		}
		url = strings.TrimSuffix(url, "&")
	}
	return url
}

func GetUrlByEnv() string {
	env := properties.GetProperty().GetEnv()
	url := "https://{env}-policyorch.ftc-cx.tech/{release}"
	switch env {
	case "TEST":
		url = strings.ReplaceAll(url, "{env}", "uat")
	case "PROD":
		url = strings.ReplaceAll(url, "{env}", "prd")
	default:
		url = strings.ReplaceAll(url, "{env}", strings.ToLower(env))
	}
	return url
}

func GetServiceByLayerAndName(layerName string, serviceName string) *properties.Service {
	var serviceResult *properties.Service
	for _, service := range properties.GetProperty().External.Services {
		if strings.ToLower(service.Layer) == strings.ToLower(layerName) && strings.ToLower(service.Name) == strings.ToLower(serviceName) {
			serviceResult = service
			break
		}
	}
	if serviceResult.Client == nil {
		if serviceResult.Timeout == 0 {
			serviceResult.Timeout = 10
		}
		timeout := time.Second * time.Duration(serviceResult.Timeout)
		serviceResult.Client = &http.Client{Timeout: timeout}
	}
	return serviceResult
}
*/
