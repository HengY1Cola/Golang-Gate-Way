package tcpProxyRouter

import (
	"context"
	"fmt"
	"gin/dao"
	"gin/reverseProxy"
	"gin/tcpProxyMiddleware"
	"gin/tcpserver"
	"log"
	"net"
)

var tcpServerList []*tcpserver.TcpServer

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

// TcpServerRun 启动Tcp服务
func TcpServerRun() {
	serviceList, _ := dao.ServiceMangerHandler.GetTcpServiceList()
	for _, each := range serviceList {
		temp := each
		log.Printf(" [INFO] TcpProxyRun:%v\n", temp.Tcp.Port)
		go func(serviceInfo *dao.ServiceDetail) { // 开启协程
			rb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceInfo)
			if err != nil {
				log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", fmt.Sprintf(":%d", serviceInfo.Tcp.Port), err)
				return
			}
			// todo 进入中间件
			router := tcpProxyMiddleware.NewTcpSliceRouter()
			router.Group("/").Use(
				tcpProxyMiddleware.TCPFlowCountMiddleware(),
				tcpProxyMiddleware.TCPFlowLimitMiddleware(),
				tcpProxyMiddleware.TcpWhiteAndBlackList(),
			)

			// todo 加入到List并开启服务
			routerHandler := tcpProxyMiddleware.NewTcpSliceRouterHandler(
				func(c *tcpProxyMiddleware.TcpSliceRouterContext) tcpserver.TCPHandler {
					return reverseProxy.NewTcpLoadBalanceReverseProxy(c, rb)
				}, router)
			baseCtx := context.WithValue(context.Background(), "service", serviceInfo)
			tcpServer := &tcpserver.TcpServer{
				Addr:    fmt.Sprintf(":%d", serviceInfo.Tcp.Port), // 定义端口
				Handler: routerHandler,
				BaseCtx: baseCtx,
			}
			tcpServerList = append(tcpServerList, tcpServer)
			err = tcpServer.ListenAndServe() // 开始监听
			if err != nil && err != tcpserver.ErrServerClosed {
				log.Fatalf(" [INFO] TcpProxyRun:%v err:%v\n", temp.Tcp.Port, err)
			}
		}(temp)
	}

}

// TcpServerStop 关闭Tcp服务
func TcpServerStop() {
	for _, each := range tcpServerList {
		each.Close()
		log.Printf(" [INFO] TcpProxyStop stopped %v\n", each.Addr)
	}
}
