package httpProxyMiddleware

import (
	"errors"
	"fmt"
	"gin/dao"
	"gin/middleware"
	"gin/public"
	"github.com/gin-gonic/gin"
	"strings"
)

// HttpWhiteAndBlackList 黑白名单管理
func HttpWhiteAndBlackList() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo 拿到上一个中间件中的service信息
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 10002, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*dao.ServiceDetail) // 转型
		// todo 拿到名单进行管理
		whiteList := strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		if serviceDetail.AccessControl.OpenAuth == 1 && whiteList[0] != "" { // 进入白名单匹配
			if !public.InStringSlice(whiteList, c.ClientIP()) {
				middleware.ResponseError(c, 10007, errors.New(fmt.Sprintf("%s 不存在与白名单中", c.ClientIP())))
				c.Abort()
				return
			}
		} else if serviceDetail.AccessControl.OpenAuth == 1 && whiteList[0] == "" { // 开启了校验并且没有设置黑名单
			blackList := strings.Split(serviceDetail.AccessControl.BlackList, ",")
			if blackList[0] != "" {
				if public.InStringSlice(blackList, c.ClientIP()) { // 如果在黑名单中
					middleware.ResponseError(c, 10008, errors.New(fmt.Sprintf("%s 被拉黑了", c.ClientIP())))
					c.Abort()
					return
				}
			}
		}
		c.Next()
	}
}
