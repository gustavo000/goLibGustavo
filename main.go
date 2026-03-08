package main

import (
	"github.com/gustavo000/goLibGustavo/root"
	"github.com/gustavo000/goLibGustavo/routing"
)

func main() {
	//database.ConnectDb("g", "g", "localhost", "5432", "inventory")

	root.InitServer(append(routing.DefaultRoutes, routing.InternalRoutes...), append(root.DefaultServices, root.InternalServices...))

}
