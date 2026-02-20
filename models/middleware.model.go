package models

import (
	"net/http"
	"time"

	"github.com/gustavo000/goLibGustavo/models/rest"
)

type Middleware struct {
	TimeRequest time.Time
	Name        string
	Service     func(ctx http.Handler, r *http.Request) *rest.Response
	//DefaultHeaders []headers_validator.Header
	//CustomHeaders  []headers_validator.Header
	LogTransaction bool
	SkipHandler    bool
	LogBody        bool
}
