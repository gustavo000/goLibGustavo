package rest

type ParserParam struct {
	Key   string
	Value string
}

type QueryParam struct {
	Key   string
	Value string
}

type Url struct {
	IngressForce  string
	LayerName     string
	ServiceName   string
	EndpointName  string
	EndpointForce string
	QueryForce    string
	ForceAtEnd    string
	ParserParams  []ParserParam
	QueryParams   []QueryParam
	IsLocal       bool
}
