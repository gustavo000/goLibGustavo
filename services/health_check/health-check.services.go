package health_check

import (
	"net/http"

	"github.com/gustavo000/goLibGustavo/models/properties"
	"github.com/gustavo000/goLibGustavo/models/rest"
	"github.com/gustavo000/goLibGustavo/pkg/functions"
)

func HealthCheck() {

}

func CheckStatus(w http.ResponseWriter, r *http.Request) *rest.Response {
	responseMap := make(map[string]any)
	responseMap["version"] = properties.GetProperty().GetVersion()
	responseMap["status"] = "UP"
	responseMap["deployedAt"] = properties.GetProperty().GetStartUpTime()
	responseMap["dependencies"] = functions.GetDependencies()
	responseMap["name"] = properties.GetProperty().GetName()
	responseMap["performanceStats"] = getPerformanceStats()
	return functions.GenerateHttpResponse(200, responseMap)
}
