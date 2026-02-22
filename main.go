package main

import (
	"github.com/gustavo000/goLibGustavo/root"
	"github.com/gustavo000/goLibGustavo/routing"
)

func main() {
	allRoutes := append(routing.DefaultRoutes, routing.InternalRoutes...)
	root.InitServer(allRoutes)

}
