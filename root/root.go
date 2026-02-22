package root

import (
	"encoding/json"
	"os"

	"github.com/gustavo000/goLibGustavo/init/external_services"
	"github.com/gustavo000/goLibGustavo/init/net_http"
	"github.com/gustavo000/goLibGustavo/models/properties"
	"github.com/gustavo000/goLibGustavo/models/rest"
)

var service = []*rest.Service{
	{
		Name:    "Core",
		Ingress: "http://cust-rtmn-orch-dl-core-{release}-service.{namespace}.svc.cluster.local",
		Path:    "/dl-core",
		Timeout: 30,
	},
}

func InitServer(routes rest.Routes) {
	properties.NewProperties(
		properties.WithServiceName("github.com/gustavo000/goLibGustavo"),
		properties.WithBasePath("/library"),
		properties.WithPort("8080"),
		properties.WithNameSpace("orch"),
		properties.WithRelease("r1"),
		properties.WithVersion(getCurrentVersion()),
		properties.WithEnvironment("PROD"),
		properties.WithServices(service...))
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
