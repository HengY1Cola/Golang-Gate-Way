package httpProxyMiddleware

import (
	"errors"
	"gin/dao"
	"gin/middleware"
	"gin/public"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

// HttpJwtAuth Jwt验证
func HttpJwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo 拿到上一个中间件中的service信息
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 10002, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*dao.ServiceDetail) // 转型
		// todo 解码jwt信息
		if serviceDetail.AccessControl.OpenAuth == 1 {
			matched := false
			token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "") // 一定要把空格加上
			if token != "" {
				log.Println(token)
				decode, err := public.JwtDecode(token)
				if err != nil {
					middleware.ResponseError(c, 10014, err)
					c.Abort()
					return
				}
				appList := dao.AppManagerHandler.GetAppList()
				for _, each := range appList {
					if each.AppID == decode.Issuer { // 匹配到了
						c.Set("appDetail", each) // 将Info放到上下文中
						matched = true
						break
					}
				}
			}
			if !matched {
				middleware.ResponseError(c, 10015, errors.New("权限未校验成功"))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
