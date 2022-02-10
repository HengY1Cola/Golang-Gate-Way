package httpProxyMiddleware

import (
	"gin/dao"
	"gin/middleware"
	"gin/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// HttpJwtFlowLimiterMiddleware Jwt限流器
func HttpJwtFlowLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo 拿到中间件中的service信息
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 10002, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*dao.ServiceDetail) // 转型
		if serviceDetail.AccessControl.OpenAuth == 1 {
			appInterface, ok := c.Get("appDetail")
			if !ok {
				middleware.ResponseError(c, 10019, errors.New("服务器错误"))
				c.Abort()
				return
			}
			appDetail := appInterface.(*dao.App) // 转型
			// todo 对请求者的频率做要求
			if appDetail.Qps > 0 {
				Limiter, err := public.FlowLimiterHandler.GetLimiter(public.FlowAppPrefix+appDetail.AppID+"_"+c.ClientIP(), float64(appDetail.Qps))
				if err != nil {
					middleware.ResponseError(c, 10020, err)
					c.Abort()
					return
				}
				if !Limiter.Allow() {
					middleware.ResponseError(c, 10021, errors.New("进行限流，请稍后再试"))
					c.Abort()
					return
				}
			}
		}
		c.Next()
	}
}
