package routing

import (
	"net/http"

	"github.com/gustavo000/goLibGustavo/models/rest"
	"github.com/gustavo000/goLibGustavo/pkg/functions"
)

var loadedRoutes rest.Routes

func GetEndpoints() []rest.Endpoint {
	var endpoints []rest.Endpoint
	for _, route := range loadedRoutes {
		endpoint := rest.Endpoint{
			Name:   route.Controller.Name,
			Path:   route.Pattern,
			Method: route.Method,
		}
		endpoints = append(endpoints, endpoint)
	}
	return endpoints
}
func GetAllEndpoints(w http.ResponseWriter, r *http.Request) *rest.Response {
	return functions.GenerateHttpResponse(200, GetEndpoints())
}
