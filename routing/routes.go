package routing

import (
	"github.com/gustavo000/goLibGustavo/models"
	health_check "github.com/gustavo000/goLibGustavo/services/health-check"
)

var controllers = models.Routes{
	{
		Method:  "GET",
		Pattern: "/endpoints",
		Controller: models.Controller{
			Name:        "GetAllEndpoints",
			Service:     GetAllEndpoints,
			SkipHandler: true,
		},
	},
	{
		Method:  "GET",
		Pattern: "/healthcheck",
		Controller: models.Controller{
			Name:        "HealthCheck",
			Service:     health_check.CheckStatus,
			SkipHandler: true,
		},
	},
}
