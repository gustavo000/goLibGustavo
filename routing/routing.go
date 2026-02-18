package routing

import (
	"github.com/gustavo000/goLibGustavo/models"
	"github.com/gustavo000/goLibGustavo/pkg/functions"
)

func GetEndpoints() {

}

func GetAllEndpoints() *models.Response {
	return functions.GenerateHttpResponse(200, GetEndpoints())
}
