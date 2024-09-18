package router

import "github.com/casbin/casbin-server/handler"

type Router interface {
	POST(path string, handler HandlerFunc)
	Serve(addr string) error
}

type HandlerFunc func(handler.Context)
