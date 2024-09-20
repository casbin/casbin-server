package gin

import ( // Import the handler package that contains Context and HandlerFunc

	ginHandler "github.com/casbin/casbin-server/handler/gin"
	"github.com/casbin/casbin-server/router"
	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	engine *gin.Engine
}

func New() *GinRouter {
	return &GinRouter{engine: gin.Default()}
}

func (r *GinRouter) GET(path string, handler router.HandlerFunc) {
	r.engine.GET(path, func(c *gin.Context) {
		handler(&ginHandler.GinContext{Ctx: c}) // Convert gin.Context to your custom Context
	})
}

func (r *GinRouter) POST(path string, handler router.HandlerFunc) {
	r.engine.POST(path, func(c *gin.Context) {
		handler(&ginHandler.GinContext{Ctx: c}) // Convert gin.Context to your custom Context
	})
}
func (r *GinRouter) PUT(path string, handler router.HandlerFunc) {
	r.engine.PUT(path, func(c *gin.Context) {
		handler(&ginHandler.GinContext{Ctx: c}) // Convert gin.Context to your custom Context
	})
}
func (r *GinRouter) DELETE(path string, handler router.HandlerFunc) {
	r.engine.DELETE(path, func(c *gin.Context) {
		handler(&ginHandler.GinContext{Ctx: c}) // Convert gin.Context to your custom Context
	})
}
func (r *GinRouter) OPTIONS(path string, handler router.HandlerFunc) {
	r.engine.OPTIONS(path, func(c *gin.Context) {
		handler(&ginHandler.GinContext{Ctx: c}) // Convert gin.Context to your custom Context
	})
}
func (r *GinRouter) Serve(addr string) error {
	return r.engine.Run(addr)
}
