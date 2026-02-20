package rest

import "net/http"

type Service struct {
	Name       string
	Path       string
	Endpoints  []Endpoint
	Timeout    int64
	Client     *http.Client
	IsExternal bool
	Middleware bool
	Ingress    string
	Layer      string
}

type Endpoint struct {
	Name   string
	Path   string
	Method string
}

type ServiceJson struct {
	Path       string         `json:"path"`
	Endpoints  []EndpointJson `json:"endpoints"`
	Timeout    int64          `json:"timeout"`
	Client     *http.Client   `json:"client"`
	IsExternal bool           `json:"isExternal"`
	Layer      string         `json:"filter"`
	Name       string         `json:"name"`
	Ingress    string         `json:"ingress"`
}

type EndpointJson struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
}
