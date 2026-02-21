package routing

import (
	"net/http"

	"github.com/gustavo000/goLibGustavo/models/rest"
	"github.com/gustavo000/goLibGustavo/pkg/functions"
)

var loadedRoutes rest.Routes

type HttpOption func(server *http.Server)

func WithBasePathAndRoutes(basePath string, routes rest.Routes) HttpOption {
	loadedRoutes = append(loadedRoutes, routes...)
	for _, route := range DefaultRoutes {
		if route.Pattern != "/healthcheck" && route.Pattern != "/endpoints" && route.Pattern != "/metrics" {
			loadedRoutes = append(loadedRoutes, route)
		}
	}
	return func(app *http.Server) {
		baseAPI := app.Addr(basePath)
		app.Addr = basePath
		baseAPI.Use(handlers.OnRequestHandler())
		{
			for _, route := range DefaultRoutes {
				baseAPI.Get(route.Pattern, route.Controller.Handler)
			}

			for _, route := range routes {
				switch route.Method {
				case "GET":
					baseAPI.Get(route.Pattern, route.Controller.Handler)
				case "POST":
					baseAPI.Post(route.Pattern, route.Controller.Handler)
				case "PUT":
					baseAPI.Put(route.Pattern, route.Controller.Handler)
				case "PATCH":
					baseAPI.Patch(route.Pattern, route.Controller.Handler)
				case "DELETE":
					baseAPI.Delete(route.Pattern, route.Controller.Handler)
				}
			}
		}
	}
}

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
