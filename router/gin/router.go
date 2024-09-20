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

type GinRouterGroup struct {
	routerGroup *gin.RouterGroup
}

func NewGroup(routerGroup *gin.RouterGroup) *GinRouterGroup {
	return &GinRouterGroup{routerGroup: routerGroup}
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

func (r *GinRouter) Group(relativePath string, handlers ...router.HandlerFunc) router.RouterGroup {
	return NewGroup(r.engine.Group(relativePath))
}

func (r *GinRouterGroup) GET(path string, handler router.HandlerFunc) {
	r.routerGroup.GET(path, func(c *gin.Context) {
		handler(&ginHandler.GinContext{Ctx: c}) // Convert gin.Context to your custom Context
	})
}

func (r *GinRouterGroup) POST(path string, handler router.HandlerFunc) {
	r.routerGroup.POST(path, func(c *gin.Context) {
		handler(&ginHandler.GinContext{Ctx: c}) // Convert gin.Context to your custom Context
	})
}
func (r *GinRouterGroup) PUT(path string, handler router.HandlerFunc) {
	r.routerGroup.PUT(path, func(c *gin.Context) {
		handler(&ginHandler.GinContext{Ctx: c}) // Convert gin.Context to your custom Context
	})
}
func (r *GinRouterGroup) DELETE(path string, handler router.HandlerFunc) {
	r.routerGroup.DELETE(path, func(c *gin.Context) {
		handler(&ginHandler.GinContext{Ctx: c}) // Convert gin.Context to your custom Context
	})
}
func (r *GinRouterGroup) OPTIONS(path string, handler router.HandlerFunc) {
	r.routerGroup.OPTIONS(path, func(c *gin.Context) {
		handler(&ginHandler.GinContext{Ctx: c}) // Convert gin.Context to your custom Context
	})
}

func (r *GinRouterGroup) Group(relativePath string, handlers ...router.HandlerFunc) router.RouterGroup {
	return NewGroup(r.routerGroup.Group(relativePath))
}
