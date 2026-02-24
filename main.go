package main

import (
	"github.com/gustavo000/goLibGustavo/database"
	"github.com/gustavo000/goLibGustavo/root"
	"github.com/gustavo000/goLibGustavo/routing"
)

func main() {
	database.ConnectDb("g", "g", "localhost", "5432", "inventory")
	allRoutes := append(routing.DefaultRoutes, routing.InternalRoutes...)
	root.InitServer(allRoutes)

}
