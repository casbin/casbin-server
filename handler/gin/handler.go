package gin

import (
	"context"
	"log"
	"net/http"

	"github.com/casbin/casbin-server/dto"
	"github.com/casbin/casbin-server/handler"
	"github.com/casbin/casbin-server/proto"
	server "github.com/casbin/casbin-server/server"
	"github.com/gin-gonic/gin"
)

type GinContext struct {
	Ctx *gin.Context
}

func (g *GinContext) Bind(v interface{}) error {
	return g.Ctx.Bind(v)
}

func (g *GinContext) ShouldBind(v interface{}) error {
	return g.Ctx.ShouldBind(v)
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

type httpHandler struct {
	server *server.Server
}

func NewHttpHandler(
	server *server.Server,
) handler.HttpHandler {
	return &httpHandler{
		server: server,
	}
}

func (h *httpHandler) Enforce(c handler.Context) {
	var response dto.EnforceResponse
	var request dto.EnforceRequest
	ctx := context.Background()
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Create New Enforcer if input handler == -1
	if request.GetEnforcerHandler() == -1 {
		e, err := h.server.NewEnforcer(ctx, &proto.NewEnforcerRequest{AdapterHandle: -1})
		if err != nil {
			log.Println("Error at calling NewEnforcer in Enforce handler:", err.Error())
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		request.EnforcerHandler = e.GetHandler()
	}
	allowed, err := h.server.Enforce(
		ctx,
		&proto.EnforceRequest{
			EnforcerHandler: request.EnforcerHandler,
			Params:          request.Params,
		},
	)
	if err != nil {
		log.Println("Error at calling Enforce in Enforce handler:", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response.BoolReply = *allowed

	if response.GetRes() {
		c.JSON(http.StatusOK, response)
		return
	} else {
		c.JSON(http.StatusForbidden, response)
		return
	}
}
