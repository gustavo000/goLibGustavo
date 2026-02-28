package rest

import (
	"net/http"
	"time"
)

type Controller struct {
	TimeRequest time.Time
	Name        string
	Service     func(r *http.Request, res *Response) *Response
	//DefaultHeaders []headers_validator.Header
	//CustomHeaders  []headers_validator.Header
	LogTransaction bool
	SkipHandler    bool
	LogBody        bool
}
