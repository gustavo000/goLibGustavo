package routing

import (
	"github.com/gustavo000/goLibGustavo/models/rest"
	health_check "github.com/gustavo000/goLibGustavo/services/health_check"
)

var controllers = rest.Routes{
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
		Controller: rest.Controller{
			Name:        "HealthCheck",
			Service:     health_check.CheckStatus,
			SkipHandler: true,
		},
	},
}
