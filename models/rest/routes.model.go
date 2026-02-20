package rest

type Routes []Route
type Route struct {
	Method     string
	Pattern    string
	Controller Controller
	Query      []string
}
