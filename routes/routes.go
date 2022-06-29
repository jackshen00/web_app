package routes

import (
	"fmt"
	"net/http"

	"time"
	"web_app/controllers"
	"web_app/logger"
	"web_app/middlewares"
	"web_app/settings"

	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"

	_ "web_app/docs"

	"github.com/gin-contrib/pprof"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUp(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 设置成发布模式
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// 注册业务路由
	v1.POST("/signup", controllers.SignUpHandler)

	v1.POST("/login", controllers.Loginhandler)

	v1.Use(middlewares.JWTAuthMiddleware(), middlewares.RateLimitMiddleware(2*time.Second, 1)) // 应用JWT认证中间件

	//r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	//	// 如果是登录的用户,判断请求头中是否有 有效的JWT  ？
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "ok",
	//	})
	//})

	{
		v1.GET("/community", controllers.CommunityHandler)
		v1.GET("/community/:id", controllers.CommunityDetailHandler)

		v1.POST("/post", controllers.CreatePostHandler)
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
		v1.GET("/posts", controllers.GetPostListHandler)

		// 根据时间或分数获取帖子
		v1.GET("/posts2", controllers.GetPostListHandler2)

		v1.POST("/vote", controllers.PostVoteController)
	}

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.AppConfig.Version)
		fmt.Printf("%#v \n", settings.Conf.AppConfig.Version)
	})

	pprof.Register(r) // 注册pprof路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
