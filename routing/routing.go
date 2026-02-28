package routing

import (
	"net/http"

	"github.com/gustavo000/goLibGustavo/models/rest"
)

var loadedRoutes rest.Routes

type HttpOption func(server *http.Server)

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
func GetAllEndpoints(r *http.Request, res *rest.Response) *rest.Response {
	return res.WithStatus(http.StatusOK).WithBody(GetEndpoints())
}
