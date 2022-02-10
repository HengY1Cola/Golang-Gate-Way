package httpProxyMiddleware

import (
	"gin/dao"
	"gin/middleware"
	"gin/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// HttpFlowCountMiddleware 流量统计
func HttpFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo 拿到中间件中的service信息
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 10002, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*dao.ServiceDetail) // 转型
		// todo 拿到信息开始处理
		// 分为 1. 总体 2. 服务 3. 用户
		// 总体
		totalCounter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
		if err != nil {
			middleware.ResponseError(c, 10009, err)
			c.Abort()
			return
		}
		totalCounter.Increase() // +1
		// 服务
		serviceCounter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			middleware.ResponseError(c, 10009, err)
			c.Abort()
			return
		}
		serviceCounter.Increase() // +1
		c.Next()
	}
}
