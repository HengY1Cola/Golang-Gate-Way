package httpProxyMiddleware

import (
	"errors"
	"gin/dao"
	"gin/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"regexp"
	"strings"
)

// HttpUrlRewriteMiddleware Url重写
func HttpUrlRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo 从Gin中获取service信息
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 10002, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*dao.ServiceDetail) // 转型
		// todo 拿到规则开始处理
		for _, each := range strings.Split(serviceDetail.Http.UrlRewrite, ",") {
			eachOne := strings.Split(each, " ") // 第一个为正则 第二个为替换后
			if len(eachOne) != 2 {
				continue
			}
			regexp, err := regexp.Compile(eachOne[0])
			if err != nil {
				log.Println("正则匹配错误")
				continue
			}
			replacePath := regexp.ReplaceAll([]byte(c.Request.URL.Path), []byte(eachOne[1]))
			c.Request.URL.Path = string(replacePath)
		}
		c.Next()
	}
}
