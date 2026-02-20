package controllers

import (
	"net/http"
	"time"

	"github.com/gustavo000/goLibGustavo/models/rest"
)

type Controller struct {
	TimeRequest    time.Time
	Name           string
	Service        func(w http.ResponseWriter, r *http.Request) *rest.Response
	LogTransaction bool
	SkipHandler    bool
	LogBody        bool
}
