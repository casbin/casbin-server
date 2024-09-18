package handler

type HttpHandler interface {
	Enforce(c Context)
}

type Context interface {
	Bind(interface{}) error
	JSON(int, interface{}) error
	Param(string) string
	QueryParam(string) string
}
