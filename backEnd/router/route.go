package router

import (
	"gin/controller"
	"gin/docs"
	"gin/middleware"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// programatically set swagger info
	docs.SwaggerInfo.Title = lib.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = lib.GetStringConf("base.swagger.desc")
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = lib.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = lib.GetStringConf("base.swagger.base_path")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 开始注册内容
	store, err := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret")) // 将session存储到Redis做一个持久化
	if err != nil {
		log.Fatalf("sessions.NewRedisStore Err: %v", err) // 如果连接失败则服务器都没启动
	}

	// 登录路由
	adminLoginRouter := router.Group("/admin_login")
	adminLoginRouter.Use(
		sessions.Sessions("mySession", store), // 存储session
		middleware.RecoveryMiddleware(),       // 捕获所有panic，并且返回错误信息
		middleware.RequestLog(),               // 请求输出日志,经过这个接口的都会记录到日志中
		middleware.TranslationMiddleware(),    // 翻译
	)
	{ //注册到子路由中
		controller.AdminLoginRegister(adminLoginRouter)
	}

	// 登录成功后路由
	adminRouter := router.Group("/admin")
	adminRouter.Use(
		sessions.Sessions("mySession", store), // 存储session
		middleware.RecoveryMiddleware(),       // 捕获所有panic，并且返回错误信息
		middleware.RequestLog(),               // 请求输出日志,经过这个接口的都会记录到日志中
		middleware.SessionAuthMiddleware(),    // session中间件的校验
		middleware.TranslationMiddleware(),    // 翻译
	)
	{ //注册到子路由中
		controller.AdminRegister(adminRouter)
	}

	// 服务路由
	serviceRouter := router.Group("/service")
	serviceRouter.Use(
		sessions.Sessions("mySession", store), // 存储session
		middleware.RecoveryMiddleware(),       // 捕获所有panic，并且返回错误信息
		middleware.RequestLog(),               // 请求输出日志,经过这个接口的都会记录到日志中
		middleware.SessionAuthMiddleware(),    // session中间件的校验
		middleware.TranslationMiddleware(),    // 翻译
	)
	{ //注册到子路由中
		controller.ServiceRegister(serviceRouter)
	}

	// 租户路由
	appRouter := router.Group("/app")
	appRouter.Use(
		sessions.Sessions("mySession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware())
	{
		controller.APPRegister(appRouter)
	}

	// 首页大盘
	mainRouter := router.Group("/main")
	mainRouter.Use(
		sessions.Sessions("mySession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware())
	{
		controller.DashboardRegister(mainRouter)
	}
	return router
}
