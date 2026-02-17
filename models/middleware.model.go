package models

import (
	"net/http"
	"time"
)

type Middleware struct {
	TimeRequest time.Time
	Name        string
	Service     func(ctx http.Handler, r *http.Request) *Response
	//DefaultHeaders []headers_validator.Header
	//CustomHeaders  []headers_validator.Header
	LogTransaction bool
	SkipHandler    bool
	LogBody        bool
}
