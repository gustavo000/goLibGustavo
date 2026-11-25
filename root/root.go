package root

import (
	"encoding/json"
	"os"

	"github.com/gustavo000/goLibGustavo/init/external_services"
	"github.com/gustavo000/goLibGustavo/init/net_http"
	"github.com/gustavo000/goLibGustavo/models/properties"
	"github.com/gustavo000/goLibGustavo/models/rest"
)

var DefaultServices = []*properties.Service{
	{
		Name:    "Google",
		Ingress: "http://www.google.com/",
		Path:    "/",
		Timeout: 30,
	},
}
var InternalServices = []*properties.Service{}

// environmentFile mirrors the structure of environment/*.properties.json.
type environmentFile struct {
	Name               string `json:"name"`
	BasePath           string `json:"basePath"`
	Port               string `json:"port"`
	Namespace          string `json:"namespace"`
	Release            string `json:"release"`
	DefaultEnvironment string `json:"defaultEnvironment"`
	PublicRateLimit    int    `json:"publicRateLimit"`
	PrivateRateLimit   int    `json:"privateRateLimit"`
}

func InitServer(routes rest.Routes, services []*properties.Service) {
	envs := getEnviroment()
	properties.NewProperties(
		properties.WithServiceName(envs.Name),
		properties.WithBasePath(envs.BasePath),
		properties.WithPort(envs.Port),
		properties.WithNameSpace(envs.Namespace),
		properties.WithRelease(envs.Release),
		properties.WithVersion(getCurrentVersion()),
		properties.WithEnvironment(envs.DefaultEnvironment),
		properties.WithPublicRateLimit(envs.PublicRateLimit),
		properties.WithPrivateRateLimit(envs.PrivateRateLimit),
		properties.WithServices(services...))
	external_services.GetEndpointsOfServices(properties.GetProperty().GetEnv())
	net_http.StartHttp(routes, properties.GetProperty().GetPort())
}

func getCurrentVersion() string {
	var res map[string]string
	bytes, _ := os.ReadFile("tag_version.json")
	err := json.Unmarshal(bytes, &res)
	if err != nil {
		return ""
	}
	return res["version"]
}

func getEnviroment() environmentFile {
	var res environmentFile
	bytes, _ := os.ReadFile("environment/local.properties.json")
	err := json.Unmarshal(bytes, &res)
	if err != nil {
		return environmentFile{}
	}
	return res
}
