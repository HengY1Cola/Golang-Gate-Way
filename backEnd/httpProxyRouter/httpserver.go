package httpProxyRouter

import (
	"context"
	"gin/certFile"
	"gin/middleware"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler  *http.Server
	HttpsSrvHandler *http.Server
)

// ----------------------------------- Http -----------------------------------

func HttpServerRun() {
	gin.SetMode(lib.GetStringConf("proxy.base.debug_mode"))
	r := InitRouter( // 复用中间件
		middleware.RecoveryMiddleware(), // 捕获所有panic，并且返回错误信息
		middleware.RequestLog(),         // 请求输出日志,经过这个接口的都会记录到日志中)
	) // 初始化中间件
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.http.max_header_bytes")),
	}
	log.Printf(" [INFO] HttpProxyRun:%s\n", lib.GetStringConf("proxy.http.addr"))
	if err := HttpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] HttpProxyRun:%s err:%v\n", lib.GetStringConf("proxy.http.addr"), err)
	}
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpProxyStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpProxyStop stopped\n")
}

// ----------------------------------- Https -----------------------------------

func HttpsServerRun() {
	gin.SetMode(lib.GetStringConf("proxy.base.debug_mode"))
	r := InitRouter(
		middleware.RecoveryMiddleware(), // 捕获所有panic，并且返回错误信息
		middleware.RequestLog(),         // 请求输出日志,经过这个接口的都会记录到日志中)
	) // 初始化中间件
	HttpsSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.https.max_header_bytes")),
	}
	log.Printf(" [INFO] HttpsProxyRun:%s\n", lib.GetStringConf("proxy.https.addr"))
	if err := HttpsSrvHandler.ListenAndServeTLS(certFile.Path("server.crt"), certFile.Path("server.key")); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] HttpsProxyRun:%s err:%v\n", lib.GetStringConf("proxy.https.addr"), err)
	}
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpsProxyStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpsProxyStop stopped\n")
}
