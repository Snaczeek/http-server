package router

import (
	"snaczek-server/coreutils"
)

type HandelerFunction func(req coreutils.Request) coreutils.Respone

type Router struct {
	routes map[string]map[string]HandelerFunction // path -> method -> handler
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]HandelerFunction),
	}
}

func (r *Router) RegisterRoute(method string, path string, handler HandelerFunction) {
	if r.routes[path] == nil {
		r.routes[path] = make(map[string]HandelerFunction)
	}
	r.routes[path][method] = handler
}

func (r *Router) Route(req coreutils.Request) coreutils.Respone {
	path, ok := r.routes[req.Path]
	if !ok {
		return coreutils.Respone {
			Status_code: 404,
			Headers: map[string]string{"Content-Type": "text/plain"},
			Body: []byte("404 Not Found\n"),
		} 
	}

	handler, ok := path[req.Method]
	if !ok {
		return coreutils.Respone {
			Status_code: 405,
			Headers: map[string]string{"Content-Type": "text/plain"},
			Body: []byte("405 Method Not Allowed\n"),
		} 
	}

	return handler(req)
}

