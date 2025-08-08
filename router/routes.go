package router

import (
	"snaczek-server/handlers"
)

// RegisterAllRoutes adds all your routes to the given Router
func RegisterAllRoutes(r *Router) {
	r.RegisterRoute("GET", "/hello", handlers.HelloHandler)

}
