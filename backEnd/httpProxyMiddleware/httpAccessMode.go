package httpProxyMiddleware

import (
	"gin/dao"
	"gin/middleware"
	"github.com/gin-gonic/gin"
)

// HttpAccessModeMiddleware 定义接入方式的中间件
// 将我们的请求信息与服务列表做匹配关系
func HttpAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := dao.ServiceMangerHandler.HttpAccessMode(c)
		if err != nil {
			middleware.ResponseError(c, 10001, err)
			c.Abort()
			return
		}
		c.Set("service", service) // 能够匹配上则设置上下问信息
		c.Next()
	}
}
