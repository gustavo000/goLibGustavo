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

func InitServer(routes rest.Routes, services []*properties.Service) {
	envs := getEnviroment()
	properties.NewProperties(
		properties.WithServiceName(envs["name"]),
		properties.WithBasePath(envs["basePath"]),
		properties.WithPort(envs["port"]),
		properties.WithNameSpace(envs["namespace"]),
		properties.WithRelease(envs["release"]),
		properties.WithVersion(getCurrentVersion()),
		properties.WithEnvironment(envs["defaultEnvironment"]),
		properties.WithServices(services...))
	external_services.GetEndpointsOfServices(properties.GetProperty().GetEnv())
	net_http.StartHttp(routes)
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

func getEnviroment() map[string]string {
	var res map[string]string
	bytes, _ := os.ReadFile("environment/dev.properties.json")
	err := json.Unmarshal(bytes, &res)
	if err != nil {
		return nil
	}
	return res
}
