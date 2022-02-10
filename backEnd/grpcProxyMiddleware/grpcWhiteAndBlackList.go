package grpcProxyMiddleware

import (
	"fmt"
	"gin/dao"
	"gin/public"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"strings"
)

// GrpcWhiteAndBlackListMiddleware 匹配接入方式 基于请求信息
func GrpcWhiteAndBlackListMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// todo 处理列表
		var ipWhiteList []string
		var ipBlackList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			ipWhiteList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.BlackList != "" {
			ipBlackList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}
		// todo 判断处理
		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}
		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]
		if serviceDetail.AccessControl.OpenAuth == 1 && len(ipWhiteList) > 0 {
			if !public.InStringSlice(ipWhiteList, clientIP) {
				return errors.New(fmt.Sprintf("%s not in white ip list", clientIP))
			}
		} else if serviceDetail.AccessControl.OpenAuth == 1 && len(ipWhiteList) == 0 {
			if len(ipBlackList) != 0 {
				if public.InStringSlice(ipBlackList, clientIP) {
					return errors.New(fmt.Sprintf("%s in black ip list", clientIP))
				}
			}
		}
		// todo 收尾
		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		log.Println("我来到了这里3")
		return nil
	}
}
