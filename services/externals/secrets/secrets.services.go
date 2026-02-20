package secrets

import (
	"github.com/google/uuid"
	"github.com/gustavo000/goLibGustavo/models"
	"github.com/gustavo000/goLibGustavo/models/rest"
	headers_request "github.com/gustavo000/goLibGustavo/pkg/headers_request"
	"github.com/gustavo000/goLibGustavo/resources/properties"
	auxiliar_url "github.com/gustavo000/goLibGustavo/services/externals/auxiliar_url"
)

func GetSecretFromEnvironment(secret *properties.Secret) *rest.Response {
	requestHeaders := headers_request.GetBaseRequestHeaders()
	requestHeaders = requestHeaders.ModifyHeaders(
		headers_request.RequestHeader{Key: "X-Environment", ForceValue: properties.GetProperty().GetEnv()},
		headers_request.RequestHeader{Key: "X-Chref", ForceValue: "F_COM"},
		headers_request.RequestHeader{Key: "X-Cmref", ForceValue: "F_COM"},
		headers_request.RequestHeader{Key: "X-Country", ForceValue: "CL"},
		headers_request.RequestHeader{Key: "X-Request-Id", ForceValue: uuid.NewString()},
		headers_request.RequestHeader{Key: "X-Trace-Id", ForceValue: uuid.NewString()},
	)
	requestHeaders = requestHeaders.AddHeaders(headers_request.RequestHeader{Key: "X-Secret", ForceValue: secret.Filter})
	serviceName := "CorpEnvironment"
	helperUrl := models.Url{
		LayerName:     "Sec",
		ServiceName:   serviceName,
		EndpointForce: "/secrets",
	}
	client := auxiliar_url.GetClientBy(helperUrl)
	return client.WithHeaders(requestHeaders).MakeGet()
}
