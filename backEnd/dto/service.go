package dto

import (
	"gin/public"
	"github.com/gin-gonic/gin"
)

// ------------------ 服务列表 ------------------

type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"网关" validate:""`                  // 关键词
	PageNum  int    `json:"pageNum" form:"pageNum" comment:"页数" example:"1" validate:"required"`      // 页数
	PageSize int    `json:"pageSize" form:"pageSize" comment:"每页数量" example:"20" validate:"required"` // 每页数量
}

type ServiceItemOutput struct {
	ID          int64  `json:"id" form:"id"`                     // id
	ServiceName string `json:"service_name" form:"service_name"` // 服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc"` // 服务描述
	LoadType    int    `json:"load_type" form:"load_type"`       // 类型
	ServiceAddr string `json:"service_addr" form:"service_addr"` // 服务地址
	Qps         int64  `json:"qps" form:"qps"`                   // qps
	Qpd         int64  `json:"qpd" form:"qpd"`                   // qpd
	TotalNode   int    `json:"totalNode" form:"totalNode"`       // 节点数
}

type ServiceOutput struct {
	Total int64               `json:"total" form:"total" comment:"总数" example:"" validate:""` // 总数
	List  []ServiceItemOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`   // 列表
}

// BindValidParam 绑定到结构体以及校验参数
func (param *ServiceListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param) // 传入向下文以及校验的参数
}

// ------------------ 服务删除 ------------------

type ServiceDeleteInput struct {
	Id int `json:"id" form:"id" comment:"服务ID" example:"1" validate:"required"` // 服务ID
}

// BindValidParam 绑定到结构体以及校验参数
func (param *ServiceDeleteInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param) // 传入向下文以及校验的参数
}

// ------------------ 添加HTTP服务 ------------------

type ServiceAddHttpInput struct {
	ServiceName string `json:"serviceName" form:"serviceName" comment:"服务名称" example:"test_http_service" validate:"required,validServiceName"` // 服务名称
	ServiceDesc string `json:"serviceDesc" form:"serviceDesc" comment:"服务描述" example:"test_http_service" validate:"required,max=255,min=1"`    // 服务描述

	RuleType       int    `json:"ruleType" form:"ruleType" comment:"接入类型" example:"" validate:"max=1,min=0"`                         // 接入类型
	Rule           string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"/testHttpService" validate:"required,validRule"`   // 接入路径
	NeedHttps      int    `json:"needHttps" form:"needHttps" comment:"支持Https" example:"" validate:"max=1,min=0"`                    // 支持Https
	NeedStripUrl   int    `json:"needStripUrl" form:"needStripUrl" comment:"启用StripUrl" example:"" validate:"max=1,min=0"`           // 启用StripUrl
	NeedWebsocket  int    `json:"needWebsocket" form:"needWebsocket" comment:"支持Websocket" example:"" validate:"max=1,min=0"`        // 支持Websocket
	UrlRewrite     string `json:"urlRewrite" form:"urlRewrite" comment:"url重写功能" example:"" validate:"validUrlRewrite"`              // url重写功能
	HeaderTransfor string `json:"headerTransfor" form:"headerTransfor" comment:"header转换" example:"" validate:"validHeaderTransfor"` // header转换

	OpenAuth           int    `json:"openAuth" form:"openAuth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`              // 是否开启权限
	BlackList          string `json:"blackList" form:"blackList" comment:"黑名单" example:"" validate:""`                          // 黑名单
	WhiteList          string `json:"whiteList" form:"whiteList" comment:"白名单" example:"" validate:""`                          // 白名单
	ClientIpFlowLimit  int    `json:"clientIpFlowLimit" form:"clientIpFlowLimit" comment:"客户端限流" example:"" validate:"min=0"`   // 客户端限流
	ServiceIpFlowLimit int    `json:"serviceIpFlowLimit" form:"serviceIpFlowLimit" comment:"服务端限流" example:"" validate:"min=0"` // 服务端限流

	RoundType              int    `json:"roundType" form:"roundType" comment:"轮询方式" example:"" validate:"max=3,min=0"`                         // 轮询方式
	IpList                 string `json:"ipList" form:"ipList" comment:"ip列表" example:"127.0.0.1:80" validate:"required,validIpPortList"`      // ip列表
	WeightList             string `json:"weightList" form:"weightList" comment:"权重列表" example:"50" validate:"required,validWeightList"`        // 权重列表
	UpstreamConnectTimeout int    `json:"upstreamConnectTimeout" form:"upstreamConnectTimeout" comment:"连接超时" example:"" validate:"min=0"`     // 连接超时
	UpstreamHeaderTimeout  int    `json:"upstreamHeaderTimeout" form:"upstreamHeaderTimeout" comment:"获取Header超时" example:"" validate:"min=0"` // 获取Header超时
	UpstreamIdleTimeout    int    `json:"upstreamIdleTimeout" form:"upstreamIdleTimeout" comment:"连接最大空闲时间" example:"" validate:"min=0"`       // 连接最大空闲时间
	UpstreamMaxIdle        int    `json:"upstreamMaxIdle" form:"upstreamMaxIdle" comment:"最大空闲连接数" example:"" validate:"min=0"`                // 最大空闲连接数
}

// BindValidParam 绑定到结构体以及校验参数
func (param *ServiceAddHttpInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param) // 传入向下文以及校验的参数
}

// ------------------ 更新HTTP服务 ------------------

type ServiceUpdateHttpInput struct {
	Id          int    `json:"id" form:"id" comment:"服务编号" example:"62" validate:"required,min=1"`                                             // 服务编号
	ServiceName string `json:"serviceName" form:"serviceName" comment:"服务名称" example:"test_http_service" validate:"required,validServiceName"` // 服务名称
	ServiceDesc string `json:"serviceDesc" form:"serviceDesc" comment:"服务描述" example:"我是一条测试" validate:"required,max=255,min=1"`               // 服务描述

	RuleType       int    `json:"ruleType" form:"ruleType" comment:"接入类型" example:"" validate:"max=1,min=0"`                         // 接入类型
	Rule           string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"/testHttpService" validate:"required,validRule"`   // 接入路径
	NeedHttps      int    `json:"needHttps" form:"needHttps" comment:"支持Https" example:"" validate:"max=1,min=0"`                    // 支持Https
	NeedStripUrl   int    `json:"needStripUrl" form:"needStripUrl" comment:"启用StripUrl" example:"" validate:"max=1,min=0"`           // 启用StripUrl
	NeedWebsocket  int    `json:"needWebsocket" form:"needWebsocket" comment:"支持Websocket" example:"" validate:"max=1,min=0"`        // 支持Websocket
	UrlRewrite     string `json:"urlRewrite" form:"urlRewrite" comment:"url重写功能" example:"" validate:"validUrlRewrite"`              // url重写功能
	HeaderTransfor string `json:"headerTransfor" form:"headerTransfor" comment:"header转换" example:"" validate:"validHeaderTransfor"` // header转换

	OpenAuth           int    `json:"openAuth" form:"openAuth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`              // 是否开启权限
	BlackList          string `json:"blackList" form:"blackList" comment:"黑名单" example:"" validate:""`                          // 黑名单
	WhiteList          string `json:"whiteList" form:"whiteList" comment:"白名单" example:"" validate:""`                          // 白名单
	ClientIpFlowLimit  int    `json:"clientIpFlowLimit" form:"clientIpFlowLimit" comment:"客户端限流" example:"" validate:"min=0"`   // 客户端限流
	ServiceIpFlowLimit int    `json:"serviceIpFlowLimit" form:"serviceIpFlowLimit" comment:"服务端限流" example:"" validate:"min=0"` // 服务端限流

	RoundType              int    `json:"roundType" form:"roundType" comment:"轮询方式" example:"" validate:"max=3,min=0"`                         // 轮询方式
	IpList                 string `json:"ipList" form:"ipList" comment:"ip列表" example:"127.0.0.1:80" validate:"required,validIpPortList"`      // ip列表
	WeightList             string `json:"weightList" form:"weightList" comment:"权重列表" example:"50" validate:"required,validWeightList"`        // 权重列表
	UpstreamConnectTimeout int    `json:"upstreamConnectTimeout" form:"upstreamConnectTimeout" comment:"连接超时" example:"" validate:"min=0"`     // 连接超时
	UpstreamHeaderTimeout  int    `json:"upstreamHeaderTimeout" form:"upstreamHeaderTimeout" comment:"获取Header超时" example:"" validate:"min=0"` // 获取Header超时
	UpstreamIdleTimeout    int    `json:"upstreamIdleTimeout" form:"upstreamIdleTimeout" comment:"连接最大空闲时间" example:"" validate:"min=0"`       // 连接最大空闲时间
	UpstreamMaxIdle        int    `json:"upstreamMaxIdle" form:"upstreamMaxIdle" comment:"最大空闲连接数" example:"" validate:"min=0"`                // 最大空闲连接数
}

// BindValidParam 绑定到结构体以及校验参数
func (param *ServiceUpdateHttpInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param) // 传入向下文以及校验的参数
}

// ------------------ 服务统计 ------------------

type ServiceStatOutput struct {
	Today     []int64   `json:"today" form:"today" comment:"今日数据统计" example:"" validate:""`         // 今日数据统计
	Yesterday [24]int64 `json:"yesterday" form:"yesterday" comment:"昨日数据统计" example:"" validate:""` // 昨日数据统计
}

// ------------------ 添加GRPC服务 ------------------

type ServiceAddGrpcInput struct {
	ServiceName string `json:"serviceName" form:"serviceName" comment:"服务名称" example:"test_grpc01" validate:"required,validServiceName"` // 服务名称
	ServiceDesc string `json:"serviceDesc" form:"serviceDesc" comment:"服务描述" example:"我是grpc测试" validate:"required,max=255,min=1"`       // 服务描述
	Port        int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" example:"8888" validate:"required,min=8001,max=8999"`      // 服务描述

	HeaderTransfor     string `json:"headerTransfor" form:"headerTransfor" comment:"metadata转换" example:"" validate:"validHeaderTransfor"` // header转换
	OpenAuth           int    `json:"openAuth" form:"openAuth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`                         // 是否开启权限
	BlackList          string `json:"blackList" form:"blackList" comment:"黑名单" example:"" validate:"validIpList"`                          // 黑名单
	WhiteList          string `json:"whiteList" form:"whiteList" comment:"白名单" example:"" validate:"validIpList"`                          // 白名单
	WhiteHostName      string `json:"whiteHostName" form:"whiteHostName" comment:"白名单主机，逗号间隔" example:"" validate:"validIpList"`           // 白名单主机
	ClientIpFlowLimit  int    `json:"clientIpFlowLimit" form:"clientIpFlowLimit" comment:"客户端限流" example:"" validate:"min=0"`              // 客户端限流
	ServiceIpFlowLimit int    `json:"serviceIpFlowLimit" form:"serviceIpFlowLimit" comment:"服务端限流" example:"" validate:"min=0"`            // 服务端限流

	RoundType  int    `json:"roundType" form:"roundType" comment:"轮询方式" example:"" validate:"max=3,min=0"`                    // 轮询方式
	IpList     string `json:"ipList" form:"ipList" comment:"ip列表" example:"127.0.0.1:80" validate:"required,validIpPortList"` // ip列表
	WeightList string `json:"weightList" form:"weightList" comment:"权重列表" example:"50" validate:"required,validWeightList"`   // 权重列
	ForbidList string `json:"forbidList" form:"forbidList" comment:"禁用IP列表" example:"" validate:"validIpList"`                // 禁用IP列表
}

// BindValidParam 绑定到结构体以及校验参数
func (param *ServiceAddGrpcInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param) // 传入向下文以及校验的参数
}

// ------------------ 更新GRPC服务 ------------------

type ServiceUpdateGrpcInput struct {
	Id          int    `json:"id" form:"id" comment:"服务编号" example:"62" validate:"required,min=1"`                                              // 服务编号
	ServiceName string `json:"serviceName" form:"serviceName" comment:"服务名称" example:"test_grpc_update01" validate:"required,validServiceName"` // 服务名称
	ServiceDesc string `json:"serviceDesc" form:"serviceDesc" comment:"服务描述" example:"我是更新测试" validate:"required,max=255,min=1"`                // 服务描述
	Port        int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" example:"8889" validate:"required,min=8001,max=8999"`             // 端口

	HeaderTransfor     string `json:"headerTransfor" form:"headerTransfor" comment:"metadata转换" example:"" validate:"validHeaderTransfor"` // header转换
	OpenAuth           int    `json:"openAuth" form:"openAuth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`                         // 是否开启权限
	BlackList          string `json:"blackList" form:"blackList" comment:"黑名单" example:"" validate:"validIpList"`                          // 黑名单
	WhiteList          string `json:"whiteList" form:"whiteList" comment:"白名单" example:"" validate:"validIpList"`                          // 白名单
	WhiteHostName      string `json:"whiteHostName" form:"whiteHostName" comment:"白名单主机，逗号间隔" example:"" validate:"validIpList"`           // 白名单主机
	ClientIpFlowLimit  int    `json:"clientIpFlowLimit" form:"clientIpFlowLimit" comment:"客户端限流" example:"" validate:"min=0"`              // 客户端限流
	ServiceIpFlowLimit int    `json:"serviceIpFlowLimit" form:"serviceIpFlowLimit" comment:"服务端限流" example:"" validate:"min=0"`            // 服务端限流

	RoundType  int    `json:"roundType" form:"roundType" comment:"轮询方式" example:"" validate:"max=3,min=0"`                    // 轮询方式
	IpList     string `json:"ipList" form:"ipList" comment:"ip列表" example:"127.0.0.1:80" validate:"required,validIpPortList"` // ip列表
	WeightList string `json:"weightList" form:"weightList" comment:"权重列表" example:"50" validate:"required,validWeightList"`   // 权重列
	ForbidList string `json:"forbidList" form:"forbidList" comment:"禁用IP列表" example:"" validate:"validIpList"`                // 禁用IP列表
}

// BindValidParam 绑定到结构体以及校验参数
func (param *ServiceUpdateGrpcInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param) // 传入向下文以及校验的参数
}

// ------------------ 添加Tcp服务 ------------------

type ServiceAddTcpInput struct {
	ServiceName string `json:"serviceName" form:"serviceName" comment:"服务名称" example:"test_tcp01" validate:"required,validServiceName"` // 服务名称
	ServiceDesc string `json:"serviceDesc" form:"serviceDesc" comment:"服务描述" example:"我是一条测试" validate:"required,max=255,min=1"`        // 服务描述
	Port        int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" example:"8777" validate:"required,min=8001,max=8999"`     // 端口

	HeaderTransfor     string `json:"headerTransfor" form:"headerTransfor" comment:"header转换" example:"" validate:""`            // header转换
	OpenAuth           int    `json:"openAuth" form:"openAuth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`               // 是否开启权限
	BlackList          string `json:"blackList" form:"blackList" comment:"黑名单" example:"" validate:"validIpList"`                // 黑名单
	WhiteList          string `json:"whiteList" form:"whiteList" comment:"白名单" example:"" validate:"validIpList"`                // 白名单
	WhiteHostName      string `json:"whiteHostName" form:"whiteHostName" comment:"白名单主机，逗号间隔" example:"" validate:"validIpList"` // 白名单主机
	ClientIpFlowLimit  int    `json:"clientIpFlowLimit" form:"clientIpFlowLimit" comment:"客户端限流" example:"" validate:"min=0"`    // 客户端限流
	ServiceIpFlowLimit int    `json:"serviceIpFlowLimit" form:"serviceIpFlowLimit" comment:"服务端限流" example:"" validate:"min=0"`  // 服务端限流

	RoundType  int    `json:"roundType" form:"roundType" comment:"轮询方式" example:"" validate:"max=3,min=0"`                    // 轮询方式
	IpList     string `json:"ipList" form:"ipList" comment:"ip列表" example:"127.0.0.1:80" validate:"required,validIpPortList"` // ip列表
	WeightList string `json:"weightList" form:"weightList" comment:"权重列表" example:"50" validate:"required,validWeightList"`   // 权重列
	ForbidList string `json:"forbidList" form:"forbidList" comment:"禁用IP列表" example:"" validate:"validIpList"`                // 禁用IP列表
}

func (params *ServiceAddTcpInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// ------------------ 更新Tcp服务 ------------------

type ServiceUpdateTcpInput struct {
	Id                 int    `json:"id" form:"id" comment:"服务编号" example:"62" validate:"required,min=1"`                                             // 服务编号
	ServiceName        string `json:"serviceName" form:"serviceName" comment:"服务名称" example:"test_tcp_update01" validate:"required,validServiceName"` // 服务名称
	ServiceDesc        string `json:"serviceDesc" form:"serviceDesc" comment:"服务描述" example:"我是更新测试" validate:"required,max=255,min=1"`               // 服务描述
	Port               int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" example:"8778" validate:"required,min=8001,max=8999"`            // 端口
	OpenAuth           int    `json:"openAuth" form:"openAuth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`                                    // 是否开启权限
	BlackList          string `json:"blackList" form:"blackList" comment:"黑名单" example:"" validate:"validIpList"`                                     // 黑名单
	WhiteList          string `json:"whiteList" form:"whiteList" comment:"白名单" example:"" validate:"validIpList"`                                     // 白名单
	WhiteHostName      string `json:"whiteHostName" form:"whiteHostName" comment:"白名单主机，逗号间隔" example:"" validate:"validIpList"`                      // 白名单主机
	ClientIpFlowLimit  int    `json:"clientIpFlowLimit" form:"clientIpFlowLimit" comment:"客户端限流" example:"" validate:"min=0"`                         // 客户端限流
	ServiceIpFlowLimit int    `json:"serviceIpFlowLimit" form:"serviceIpFlowLimit" comment:"服务端限流" example:"" validate:"min=0"`                       // 服务端限流

	RoundType  int    `json:"roundType" form:"roundType" comment:"轮询方式" example:"" validate:"max=3,min=0"`                    // 轮询方式
	IpList     string `json:"ipList" form:"ipList" comment:"ip列表" example:"127.0.0.1:80" validate:"required,validIpPortList"` // ip列表
	WeightList string `json:"weightList" form:"weightList" comment:"权重列表" example:"50" validate:"required,validWeightList"`   // 权重列
	ForbidList string `json:"forbidList" form:"forbidList" comment:"禁用IP列表" example:"" validate:"validIpList"`                // 禁用IP列表
}

// BindValidParam 绑定到结构体以及校验参数
func (param *ServiceUpdateTcpInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param) // 传入向下文以及校验的参数
}
