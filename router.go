package srv

import (
	"context"
	"time"
)

type Route struct {
	handler Handler
}

func (r *Route) Handle(handler Handler) *Route {
	r.handler = handler
	return r
}

func (r *Route) HandleFunc(f func(ResponseWriter, *Request)) *Route {
	return r.Handle(HandlerFunc(f))
}

type Router struct {
	// TODO: MaxConnections is the number of incoming connections the server can
	// handle.  It is used to calculate the connection pool size.
	MaxConnections int32
	// TODO: Basecontext is unused in the router implementation. It was added to
	// mimic the server struct.
	BaseContext context.Context
	// TODO: ReadTimeout is unused in the router implementation.  It was added to
	// mimic the server struct.
	ReadTimeout time.Duration
	// routes are available to support multiple processing piplines, though they are
	// not currently supported, they will be in the future.  This would need to also
	// include either routes based on ports (multiple listeners) or routes based on
	// data (single reader and pass full data to each route).
	routes []*Route
	// middlewares contains a list of all of the registered middleware handlers.
	middlewares []middleware
}

func NewRouter() *Router {
	return &Router{}
}

func (rt *Router) NewRoute() *Route {
	route := &Route{}
	rt.routes = append(rt.routes, route)
	return route
}

func (rt *Router) Handle(handler Handler) *Route {
	return rt.NewRoute().Handle(handler)
}

func (rt *Router) HandleFunc(f func(w ResponseWriter, r *Request)) *Route {
	return rt.NewRoute().HandleFunc(f)
}

func (rt *Router) Use(mwf ...MiddlewareFunc) {
	for _, fn := range mwf {
		rt.middlewares = append(rt.middlewares, fn)
	}
}

func (rt *Router) ServeTCP(w ResponseWriter, r *Request) {
	var handler Handler
	for _, route := range rt.routes {
		handler = route.handler
		for i := len(rt.middlewares) - 1; i >= 0; i-- {
			handler = rt.middlewares[i].Middleware(handler)
		}
	}
	if handler == nil {
		panic("a handler was not defined for the router")
	}
	handler.ServeTCP(w, r)
}
