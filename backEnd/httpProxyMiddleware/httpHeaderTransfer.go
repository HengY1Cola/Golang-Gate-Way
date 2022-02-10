package httpProxyMiddleware

import (
	"gin/dao"
	"gin/middleware"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

// HeaderTransferMiddleware 头的转换
func HeaderTransferMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo 从Gin中获取service信息
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 10002, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*dao.ServiceDetail) // 转型
		// todo 拿到设置Header的值
		for _, item := range strings.Split(serviceDetail.Http.HeaderTransfor, ",") {
			eachOne := strings.Split(item, " ")
			if len(eachOne) != 3 { // 如果不合法的话 直接跳过
				continue
			}
			if eachOne[0] == "add" || eachOne[0] == "edit" {
				c.Request.Header.Set(eachOne[1], eachOne[2])
			}
			if eachOne[0] == "del" {
				c.Request.Header.Del(eachOne[1])
			}
		}
		c.Next()
	}
}
