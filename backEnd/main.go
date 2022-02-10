package main

import (
	"flag"
	"gin/dao"
	"gin/grpcProxyRouter"
	"gin/httpProxyRouter"
	"gin/router"
	"gin/tcpProxyRouter"
	"github.com/e421083458/golang_common/lib"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// endPoint dashboard后台管理 server代理服务

var (
	endpoint = flag.String("endpoint", "", "input endpoint dashboard or server")
)

func main() {
	// 解析参数
	flag.Parse()
	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *endpoint == "dashboard" { //原有的逻辑
		// 如果configPath为空的话 从命令行的-config中获取
		lib.InitModule("./conf/dev/", []string{"base", "mysql", "redis"})
		defer lib.Destroy()
		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	} else if *endpoint == "server" { // 代理服务
		lib.InitModule("./conf/dev/", []string{"base", "mysql", "redis"})
		defer lib.Destroy()
		log.Println("配置加载完毕， 代理服务器已经启动")

		dao.ServiceMangerHandler.LoadOnce() // 加载服务列表
		dao.AppManagerHandler.LoadOnce()    // 加载服务列表
		// 加入代理逻辑
		// todo 同时开启Http与Https
		go func() {
			httpProxyRouter.HttpServerRun()
		}()
		go func() {
			httpProxyRouter.HttpsServerRun()
		}()
		// todo 启动Tcp与Grpc
		go func() {
			tcpProxyRouter.TcpServerRun()
		}()
		go func() {
			grpcProxyRouter.GrpcServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		tcpProxyRouter.TcpServerStop()
		grpcProxyRouter.GrpcServerStop()
		httpProxyRouter.HttpServerStop()
		httpProxyRouter.HttpsServerStop()
	} else {
		os.Exit(3)
	}
}
