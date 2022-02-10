package httpProxyMiddleware

import (
	"gin/dao"
	"gin/middleware"
	"gin/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// HttpFlowLimiterMiddleware 限流器
func HttpFlowLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo 拿到中间件中的service信息
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 10002, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*dao.ServiceDetail) // 转型
		// todo 针对服务进行限流
		if serviceDetail.AccessControl.ServiceFlowLimit > 0 {
			qps := serviceDetail.AccessControl.ServiceFlowLimit
			serviceLimiter, err := public.FlowLimiterHandler.GetLimiter(public.FlowServicePrefix+serviceDetail.Info.ServiceName, float64(qps))
			if err != nil {
				middleware.ResponseError(c, 10010, err)
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				middleware.ResponseError(c, 10011, errors.New("进行限流，请稍后再试"))
				c.Abort()
				return
			}
		}
		// todo 对请求地址做限流(但是次数并没有)
		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			qps := serviceDetail.AccessControl.ClientIPFlowLimit
			ipLimiter, err := public.FlowLimiterHandler.GetLimiter(public.FlowServicePrefix+serviceDetail.Info.ServiceName+c.ClientIP(), float64(qps))
			if err != nil {
				middleware.ResponseError(c, 10012, err)
				c.Abort()
				return
			}
			if !ipLimiter.Allow() {
				middleware.ResponseError(c, 10013, errors.New("进行限流，请稍后再试"))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
