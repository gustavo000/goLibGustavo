package routing

import (
	"github.com/gustavo000/goLibGustavo/models/rest"
	"github.com/gustavo000/goLibGustavo/services/health_check"
)

var DefaultRoutes = rest.Routes{
	{
		Method:  "GET",
		Pattern: "/endpoints",
		Controller: rest.Controller{
			Name:        "GetAllEndpoints",
			Service:     GetAllEndpoints,
			SkipHandler: true,
		},
	},
	{
		Method:  "GET",
		Pattern: "/healthcheck",
		Public:  true,
		Controller: rest.Controller{
			Name:        "HealthCheck",
			Service:     health_check.CheckStatus,
			SkipHandler: true,
		},
	},
}

var InternalRoutes = rest.Routes{}
