package httpProxyRouter

import (
	"gin/controller"
	"gin/httpProxyMiddleware"
	"gin/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// todo 优化点
	//router := gin.New()
	router := gin.Default()
	router.Use(middlewares...) // 使用原来的中间件
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 提到前面来 如果有自己的use就不会走通用的
	// token路由
	oAuthRouter := router.Group("/oauth")
	oAuthRouter.Use(
		middleware.TranslationMiddleware(),
	)
	{
		controller.OAuthRegister(oAuthRouter)
	}

	router.Use( // 洋葱结构是要讲顺序的
		httpProxyMiddleware.HttpAccessModeMiddleware(),
		httpProxyMiddleware.HttpFlowCountMiddleware(),      // 流量统计
		httpProxyMiddleware.HttpFlowLimiterMiddleware(),    // 限流器
		httpProxyMiddleware.HttpJwtAuth(),                  // jwt验证
		httpProxyMiddleware.HttpJwtFlowCountMiddleware(),   // jwt流量统计
		httpProxyMiddleware.HttpJwtFlowLimiterMiddleware(), // jwt限流
		httpProxyMiddleware.HttpWhiteAndBlackList(),        // 黑白名单
		httpProxyMiddleware.HeaderTransferMiddleware(),     // 控制头部
		httpProxyMiddleware.HttpStripUrlMiddleware(),       // URL的替换
		httpProxyMiddleware.HttpUrlRewriteMiddleware(),     // URL重写(注意： 如果开启了替换可能不生效)
		httpProxyMiddleware.HttpReverseProxyMiddleware(),   // 反向代理
	)

	return router
}
