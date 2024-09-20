package router

import "github.com/casbin/casbin-server/handler"

type Router interface {
	GET(path string, handler HandlerFunc)
	POST(path string, handler HandlerFunc)
	PUT(path string, handler HandlerFunc)
	DELETE(path string, handler HandlerFunc)
	OPTIONS(path string, handler HandlerFunc)
	Serve(addr string) error
}

type HandlerFunc func(handler.Context)
