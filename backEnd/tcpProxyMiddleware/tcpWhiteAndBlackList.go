package tcpProxyMiddleware

import (
	"fmt"
	"gin/dao"
	"gin/public"
	"strings"
)

// TcpWhiteAndBlackList 黑白名单管理
func TcpWhiteAndBlackList() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		// todo 拿到信息并处理
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}
		// todo 进行名单限制
		var ipWhiteList []string
		var ipBlackList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			ipWhiteList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.BlackList != "" {
			ipBlackList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && ipWhiteList[0] != "" {
			if !public.InStringSlice(ipWhiteList, clientIP) {
				c.conn.Write([]byte(fmt.Sprintf("%s not in white ip list", clientIP)))
				c.Abort()
				return
			}
		} else if serviceDetail.AccessControl.OpenAuth == 1 && ipWhiteList[0] == "" {
			if ipBlackList[0] != "" {
				if public.InStringSlice(ipBlackList, clientIP) {
					c.conn.Write([]byte(fmt.Sprintf("%s in black ip list", clientIP)))
					c.Abort()
					return
				}
			}
		}
		c.Next()
	}
}
