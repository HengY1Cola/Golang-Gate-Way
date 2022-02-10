package httpProxyMiddleware

import (
	"gin/dao"
	"gin/middleware"
	"gin/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

// HttpStripUrlMiddleware 去除部分URL
func HttpStripUrlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo 从Gin中获取service信息
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 10002, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*dao.ServiceDetail) // 转型
		// todo 拿信息开始替换
		if serviceDetail.Http.RuleType == public.HttpRuleTypePrefix && serviceDetail.Http.NeedStripUri == 1 { // 前缀的情况才能替换
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, serviceDetail.Http.Rule, "", 1)
		}
		c.Next()
	}
}
