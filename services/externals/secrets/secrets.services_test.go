package secrets

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gustavo000/goLibGustavo/models"
	"github.com/gustavo000/goLibGustavo/resources/properties"
)

func Test_GetSecretFromEnvironment(t *testing.T) {
	os.Setenv("ENVIRONMENT", "BETA")

	type SecretResponse struct {
		Id    int      `json:"Id"`
		Name  string   `json:"Name"`
		Value []string `json:"Values"`
	}

	properties.NewProperties(
		properties.WithEnvironment("BETA"),
		properties.WithServices([]*properties.Service{
			{
				Layer:   "Sec",
				Name:    "CorpEnvironment",
				Ingress: "http://localhost/mock",
				Client: &http.Client{
					Transport: models.NewRoundTripForTest(
						func(req *http.Request) *http.Response {
							secretResponse := SecretResponse{
								Id:    1,
								Name:  "test",
								Value: []string{"test"},
							}

							secretResponseJSON, err := json.Marshal(secretResponse)
							assert.NoError(t, err)

							return &http.Response{
								StatusCode: http.StatusOK,
								Status:     "200 " + http.StatusText(http.StatusOK),
								Body:       io.NopCloser(bytes.NewReader(secretResponseJSON)),
								Header:     make(http.Header),
							}
						},
					),
				},
			},
		}...),
	)

	secret := &properties.Secret{
		Filter: "test",
		Name:   "test",
		Value:  "test",
	}

	response := GetSecretFromEnvironment(secret)
	assert.Equal(t, "200 OK", response.Http.Status)
	assert.Equal(t, 200, response.Http.StatusCode)

	gotB, err := io.ReadAll(response.Http.Body)
	assert.NoError(t, err)

	secretResponseExpected := SecretResponse{
		Id:    1,
		Name:  "test",
		Value: []string{"test"},
	}

	secretResponseGot := SecretResponse{}
	err = json.Unmarshal(gotB, &secretResponseGot)
	assert.NoError(t, err)
	assert.Equal(t, secretResponseExpected, secretResponseGot)
}
