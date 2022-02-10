package httpProxyMiddleware

import (
	"gin/dao"
	"gin/middleware"
	"gin/reverseProxy"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
)

// HttpReverseProxyMiddleware 定义接入方式的中间件
func HttpReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo 拿到上一个中间件中的service信息
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 10002, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*dao.ServiceDetail) // 转型
		// todo 构建Http相应代理服务
		lb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 10003, err)
			c.Abort()
			return
		}
		trans, err := dao.TransporterHandler.GetTrans(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 10004, err)
			c.Abort()
			return
		}
		proxy := reverseProxy.NewLoadBalanceReverseProxy(c, lb, trans)
		if lib.GetStringConf("proxy.base.debug_mode") == "debug" {
			log.Println("即将进入代理系统：", c.Request.URL.Path)
		}
		proxy.ServeHTTP(c.Writer, c.Request) // 执行下游服务器
	}
}
