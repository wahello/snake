package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger" //nolint: goimports
	"github.com/swaggo/gin-swagger/swaggerFiles"

	// import swagger handler
	_ "github.com/1024casts/snake/api/http" // docs is generated by Swag CLI, you have to import it.

	"github.com/1024casts/snake/api"
	"github.com/1024casts/snake/internal/handler/v1/user"
	mw "github.com/1024casts/snake/internal/middleware"
	"github.com/1024casts/snake/pkg/middleware"
)

// Load loads the middlewares, routes, handlers.
func NewRouter() *gin.Engine {
	g := gin.New()
	// 使用中间件
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middleware.Logging())
	g.Use(middleware.RequestID())
	g.Use(middleware.Prom(nil))
	g.Use(middleware.Trace())
	g.Use(mw.Translations())

	// 404 Handler.
	g.NoRoute(api.RouteNotFound)
	g.NoMethod(api.RouteNotFound)

	// 静态资源，主要是图片
	g.Static("/static", "./static")

	// swagger api docs
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// pprof router 性能分析路由
	// 默认关闭，开发环境下可以打开
	// 访问方式: HOST/debug/pprof
	// 通过 HOST/debug/pprof/profile 生成profile
	// 查看分析图 go tool pprof -http=:5000 profile
	// see: https://github.com/gin-contrib/pprof
	// pprof.Register(g)

	// HealthCheck 健康检查路由
	g.GET("/health", api.HealthCheck)
	// metrics router 可以在 prometheus 中进行监控
	// 通过 grafana 可视化查看 prometheus 的监控数据，使用插件6671查看
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 认证相关路由
	g.POST("/v1/register", user.Register)
	g.POST("/v1/login", user.Login)
	g.POST("/v1/login/phone", user.PhoneLogin)
	g.GET("/v1/vcode", user.VCode)

	// 用户
	g.GET("/v1/users/:id", user.Get)

	u := g.Group("/v1/users")
	u.Use(middleware.JWT())
	{
		u.PUT("/:id", user.Update)
		u.POST("/follow", user.Follow)
		u.GET("/:id/following", user.FollowList)
		u.GET("/:id/followers", user.FollowerList)
	}

	return g
}
