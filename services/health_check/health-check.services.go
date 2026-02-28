package health_check

import (
	"net/http"

	"github.com/gustavo000/goLibGustavo/models/properties"
	"github.com/gustavo000/goLibGustavo/models/rest"
)

func HealthCheck() {

}

type HealthHandler struct{}

func NewLHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
func CheckStatus(r *http.Request, res *rest.Response) *rest.Response {
	responseMap := make(map[string]any)
	responseMap["version"] = properties.GetProperty().GetVersion()
	responseMap["status"] = "UP"
	responseMap["deployedAt"] = properties.GetProperty().GetStartUpTime()
	responseMap["name"] = properties.GetProperty().GetName()
	return res.WithStatus(http.StatusOK).WithBody(responseMap)
}
