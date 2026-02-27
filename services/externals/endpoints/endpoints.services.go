package endpoints

import (
	"github.com/google/uuid"
	"github.com/gustavo000/goLibGustavo/models/properties"
	"github.com/gustavo000/goLibGustavo/models/rest"
	"github.com/gustavo000/goLibGustavo/pkg/headers_request"
	"github.com/gustavo000/goLibGustavo/services/externals/auxiliar_url"
)

func GetEndpoints(service *properties.Service, environment string) *properties.Service {
	requestHeaders := headers_request.GetBaseRequestHeaders()
	requestHeaders.ModifyHeaders(
		headers_request.RequestHeader{Key: "X-Environment", ForceValue: environment},
		headers_request.RequestHeader{Key: "X-Chref", ForceValue: "F_COM"},
		headers_request.RequestHeader{Key: "X-Cmref", ForceValue: "F_COM"},
		headers_request.RequestHeader{Key: "X-Country", ForceValue: "CL"},
		headers_request.RequestHeader{Key: "X-Request-Id", ForceValue: uuid.NewString()},
		headers_request.RequestHeader{Key: "X-Trace-Id", ForceValue: uuid.NewString()},
		headers_request.RequestHeader{Key: "Content-Type", ForceValue: "application/json"},
	)
	helperUrl := rest.Url{
		LayerName:     service.Layer,
		ServiceName:   service.Name,
		EndpointForce: "/endpoints",
	}
	client := auxiliar_url.GetClientBy(helperUrl)
	httpResponse := client.WithHeaders(requestHeaders).MakeGet()
	var endpoints []properties.Endpoint
	if !httpResponse.IsException() {
		_ = httpResponse.GetObject(&endpoints)
		service.Endpoints = endpoints
	}
	return service
}
