package router

import (
	"embed"
	"github.com/alyzers/alyzers/pkg/ctx"
	httpx "github.com/alyzers/alyzers/pkg/http"
	"github.com/alyzers/alyzers/pkg/http/interceptor"
	"github.com/alyzers/alyzers/pkg/http/ws"
	"github.com/alyzers/alyzers/pkg/version"
	"github.com/cnlesscode/gotool/gintool"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/9/8 15:48
 * @file: router.go
 * @description: setup router
 *  		     internal api router, use by web
 */

type Router struct {
	Http *httpx.Http
	Ctx  *ctx.Context
}

//go:embed static
var web embed.FS

func NewRouter(httpConf *httpx.Http, ctx *ctx.Context) *Router {
	return &Router{
		Http: httpConf,
		Ctx:  ctx,
	}
}

func (rt *Router) Router(log *zap.Logger) *gin.Engine {

	gin.SetMode(rt.Http.Mode)

	r := gin.New()

	r.Use(
		// cors
		interceptor.CorsInterceptor(),
		// recover
		interceptor.ExceptionInterceptor,
		// response
		interceptor.UnifiedResponseInterceptor(),
	)

	// r.Use(interceptor.AuthorizationInterceptor(rt.Http.Auth.SecretKey, rt.Http.Auth))

	// web static resource
	if rt.Http.UseFileAssets {
		r.Use(static.Serve("/", static.EmbedFolder(web, "static")))
		r.NoRoute(func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/")
		})
	}

	if rt.Http.AccessLog {
		r.Use(httpx.AccessLogFormat(log))
	}

	if rt.Http.PProf {
		pprof.Register(r, "/debug/pprof")
	}

	if rt.Http.ExposeMetrics {
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, version.GetVersion())
	})

	// alyzers router, internal api router
	engine := r.Group(rt.Http.InternalContextPath)
	{
		// ws
		engine.POST("/ws", ws.Handle)

		// core
		rt.routerGroup(engine)
	}

	return r
}

func (rt *Router) routerGroup(r *gin.RouterGroup) {

	auth := interceptor.AuthorizationInterceptor(
		rt.Http.Auth.SecretKey,
		rt.Http.Auth.RedisKeyPrefix,
		*rt.Ctx.GetRedis(),
	)

	// user
	rt.userRouter(r, auth)

	// auth
	rt.authRouter(r, auth)
}

func queryInt(r *gin.Context, key string) int {
	value, ok := gintool.QueryInt(r, key)
	if !ok {
		return 0
	}
	return value
}
