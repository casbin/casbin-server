package gin

import (
	"net/http"

	"github.com/casbin/casbin-server/dto"
	"github.com/casbin/casbin-server/handler"
	"github.com/gin-gonic/gin"
)

type GinContext struct {
	ctx *gin.Context
}

func (g *GinContext) Bind(v interface{}) error {
	return g.ctx.Bind(v)
}

func (g *GinContext) JSON(statusCode int, v interface{}) error {
	g.ctx.JSON(statusCode, v)
	return nil
}

func (g *GinContext) Param(key string) string {
	return g.ctx.Param(key)
}

func (g *GinContext) QueryParam(key string) string {
	return g.ctx.Query(key)
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
