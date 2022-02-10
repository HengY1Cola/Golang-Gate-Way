package httpProxyMiddleware

import (
	"gin/dao"
	"gin/middleware"
	"gin/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// HttpJwtFlowCountMiddleware 开启了权限的流量统计
func HttpJwtFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo 拿到中间件中的service信息
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 10002, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*dao.ServiceDetail) // 转型
		// todo 拿到app信息
		appInterface, ok := c.Get("appDetail")
		if serviceDetail.AccessControl.OpenAuth == 1 && !ok { // 开启了但是没有拿到
			middleware.ResponseError(c, 10016, errors.New("服务器错误"))
			c.Abort()
			return
		} else if !ok { // 没找到的其他情况
			c.Next()
			return
		}
		appDetail := appInterface.(*dao.App)
		// todo 记录
		appCounter, err := public.FlowCounterHandler.GetCounter(public.FlowAppPrefix + appDetail.AppID)
		if err != nil {
			middleware.ResponseError(c, 10017, err)
			c.Abort()
			return
		}
		appCounter.Increase() // 增加一个
		// todo 判断是否进行限制
		if appDetail.Qpd > 0 && appCounter.TotalCount > appDetail.Qpd {
			middleware.ResponseError(c, 10018, errors.New("用户超过日请求限制"))
			c.Abort()
			return
		}
		c.Next()
	}
}
