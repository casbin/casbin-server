package gin

import (
	"net/http"

	"github.com/casbin/casbin-server/dto"
	"github.com/casbin/casbin-server/handler"
	"github.com/gin-gonic/gin"
)

type GinContext struct {
	Ctx *gin.Context
}

func (g *GinContext) Bind(v interface{}) error {
	return g.Ctx.Bind(v)
}

func (g *GinContext) JSON(statusCode int, v interface{}) error {
	g.Ctx.JSON(statusCode, v)
	return nil
}

func (g *GinContext) Param(key string) string {
	return g.Ctx.Param(key)
}

func (g *GinContext) QueryParam(key string) string {
	return g.Ctx.Query(key)
}

type httpHandler struct{}

func NewHttpHandler() handler.HttpHandler {
	return &httpHandler{}
}

func (h *httpHandler) Enforce(c handler.Context) {
	response := dto.EnforceResponse{
		Allowed: true,
	}
	c.JSON(http.StatusOK, response)
}
