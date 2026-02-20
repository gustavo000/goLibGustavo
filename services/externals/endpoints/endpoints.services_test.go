package endpoints

import (
	c "context"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
	"github.com/gustavo000/goLibGustavo/resources/properties"
)

func getMockAppExternalService() *iris.Application {
	appMock := iris.New()
	appMock.Get("/mock/endpoints", func(ctx iris.Context) {
		header := ctx.Request().Header
		environment := header.Get("X-Environment")
		if environment == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		ctx.JSON([]map[string]interface{}{
			{
				"Name":   "test",
				"Path":   "test",
				"Method": "GET",
			},
		})
	})

	return appMock
}

func Test_GetSecretFromEnvironment(t *testing.T) {
	os.Setenv("ENVIRONMENT", "BETA")

	properties.NewProperties(
		properties.WithEnvironment("test"),
		properties.WithServices([]*properties.Service{
			{
				Layer:   "Sec",
				Name:    "CorpEnvironment",
				Ingress: "http://localhost:50070/mock",
			},
		}...),
	)

	appMock := getMockAppExternalService()
	go func() {
		err := appMock.Listen(":50070")
		assert.Equal(t, "http: Server closed", err.Error())
	}()

	defer func() {
		err := appMock.Shutdown(c.Background())
		assert.NoError(t, err)
	}()

	for {
		_, err := http.Get("http://localhost:50070")
		if err != nil && strings.Contains(err.Error(), "connection refused") {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}

	svc := &properties.Service{
		Layer: "Sec",
		Name:  "CorpEnvironment",
	}

	expectedEndpoint := properties.Endpoint{
		Name:   "test",
		Path:   "test",
		Method: "GET",
	}

	service := GetEndpoints(svc, "test")
	assert.NotNil(t, service)
	assert.Equal(t, expectedEndpoint, service.Endpoints[0])
}
