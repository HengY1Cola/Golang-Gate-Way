package grpcProxyRouter

import (
	"fmt"
	"gin/dao"
	"gin/grpcProxyMiddleware"
	"gin/reverseProxy"
	"github.com/e421083458/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

var grpcServerList []*warpGrpcServer

type warpGrpcServer struct {
	Addr string
	*grpc.Server
}

func GrpcServerRun() {
	serviceList, _ := dao.ServiceMangerHandler.GetGrpcServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *dao.ServiceDetail) {
			// todo 准备工作
			addr := fmt.Sprintf(":%d", serviceDetail.Grpc.Port)
			log.Println("addr", addr)
			rb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf(" [INFO] GetLoadBalancer %v err:%v\n", addr, err)
				return
			}
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf(" [INFO] GrpcListen %v err:%v\n", addr, err)
			}
			// todo 设置对应
			grpcHandler := reverseProxy.NewGrpcLoadBalanceHandler(rb)
			s := grpc.NewServer(
				grpc.ChainStreamInterceptor(
					grpcProxyMiddleware.GrpcFlowCountMiddleware(serviceDetail),
					grpcProxyMiddleware.GrpcFlowLimitMiddleware(serviceDetail),
					grpcProxyMiddleware.GrpcJwtAuthTokenMiddleware(serviceDetail),
					grpcProxyMiddleware.GrpcJwtFlowCountMiddleware(serviceDetail),
					grpcProxyMiddleware.GrpcJwtFlowLimitMiddleware(serviceDetail),
					grpcProxyMiddleware.GrpcWhiteAndBlackListMiddleware(serviceDetail),
					grpcProxyMiddleware.GrpcHeaderTransferMiddleware(serviceDetail),
				),
				grpc.CustomCodec(proxy.Codec()),
				grpc.UnknownServiceHandler(grpcHandler))
			err = s.Serve(lis)
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf(" [INFO] grpc_proxy_run %v err:%v\n", addr, err)
			}
			// todo 添加到slice
			grpcServerList = append(grpcServerList, &warpGrpcServer{
				Addr:   addr,
				Server: s,
			})
			log.Printf(" [INFO] grpc_proxy_run %v\n", addr)
		}(tempItem)
	}
}

func GrpcServerStop() {
	for _, grpcServer := range grpcServerList {
		grpcServer.GracefulStop()
		log.Printf(" [INFO] grpc_proxy_stop %v stopped\n", grpcServer.Addr)
	}
}
