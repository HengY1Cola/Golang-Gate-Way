package dao

import (
	"fmt"
	"gin/public"
	"gin/reverseProxy/load_balance"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type LoadBalance struct {
	ID            int64  `json:"id" gorm:"primary_key"`
	ServiceID     int64  `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	CheckMethod   int    `json:"check_method" gorm:"column:check_method" description:"检查方法 tcpchk=检测端口是否握手成功	"`
	CheckTimeout  int    `json:"check_timeout" gorm:"column:check_timeout" description:"check超时时间	"`
	CheckInterval int    `json:"check_interval" gorm:"column:check_interval" description:"检查间隔, 单位s		"`
	RoundType     int    `json:"round_type" gorm:"column:round_type" description:"轮询方式 round/weight_round/random/ip_hash"`
	IpList        string `json:"ip_list" gorm:"column:ip_list" description:"ip列表"`
	WeightList    string `json:"weight_list" gorm:"column:weight_list" description:"权重列表"`
	ForbidList    string `json:"forbid_list" gorm:"column:forbid_list" description:"禁用ip列表"`

	UpstreamConnectTimeout int `json:"upstream_connect_timeout" gorm:"column:upstream_connect_timeout" description:"下游建立连接超时, 单位s"`
	UpstreamHeaderTimeout  int `json:"upstream_header_timeout" gorm:"column:upstream_header_timeout" description:"下游获取header超时, 单位s	"`
	UpstreamIdleTimeout    int `json:"upstream_idle_timeout" gorm:"column:upstream_idle_timeout" description:"下游链接最大空闲时间, 单位s	"`
	UpstreamMaxIdle        int `json:"upstream_max_idle" gorm:"column:upstream_max_idle" description:"下游最大空闲链接数"`
}

func (t *LoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (t *LoadBalance) Find(c *gin.Context, tx *gorm.DB, search *LoadBalance) (*LoadBalance, error) {
	model := &LoadBalance{}
	err := tx.WithContext(c).Where(search).Find(model).Error
	return model, err
}

func (t *LoadBalance) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// --------------------逻辑区域---------------------

// GetIPListByModel 根据Model返回当前Ip列表
func (t *LoadBalance) GetIPListByModel() []string {
	return strings.Split(t.IpList, ",")
}

func (t *LoadBalance) GetWeightListByModel() []string {
	return strings.Split(t.WeightList, ",")
}

var LoadBalancerHandler *LoadBalancer

// LoadBalancer 为每一个服务创建负载均衡管理器
type LoadBalancer struct {
	LoadBalanceMap   map[string]*LoadBalancerItem // 多的时候直接取
	LoadBalanceSlice []*LoadBalancerItem          // 少的时候直接遍历
	Locker           sync.RWMutex
}

type LoadBalancerItem struct {
	LoadBalance load_balance.LoadBalance
	ServiceName string
}

func init() {
	LoadBalancerHandler = &LoadBalancer{
		LoadBalanceMap:   map[string]*LoadBalancerItem{},
		LoadBalanceSlice: []*LoadBalancerItem{},
		Locker:           sync.RWMutex{},
	}
}

func (lbr *LoadBalancer) GetLoadBalancer(service *ServiceDetail) (load_balance.LoadBalance, error) {
	// todo 先验证下存不存在
	for _, item := range lbr.LoadBalanceSlice {
		if item.ServiceName == service.Info.ServiceName {
			return item.LoadBalance, nil
		}
	}
	// todo 取前缀
	schema := "http://"
	if service.Http.NeedHttps == 1 {
		schema = "https://"
	}
	if service.Info.LoadType == public.LoadTypeTcp || service.Info.LoadType == public.LoadTypeGrpc {
		schema = ""
	}
	// todo 取ip/权重
	ipList := service.LoadBalance.GetIPListByModel()
	weightList := service.LoadBalance.GetWeightListByModel()
	conf := map[string]string{}
	for index, value := range ipList {
		conf[value] = weightList[index]
	}
	// todo 构造负载均衡并添加
	mConf, err := load_balance.NewLoadBalanceCheckConf(fmt.Sprintf("%s%s", schema, "%s"), conf)
	if err != nil {
		return nil, err
	}
	lb := load_balance.LoadBanlanceFactorWithConf(load_balance.LbType(service.LoadBalance.RoundType), mConf)
	lbr.LoadBalanceSlice = append(lbr.LoadBalanceSlice, &LoadBalancerItem{
		LoadBalance: lb,
		ServiceName: service.Info.ServiceName,
	})
	lbr.Locker.Lock()
	defer lbr.Locker.Unlock()
	lbr.LoadBalanceMap[service.Info.ServiceName] = &LoadBalancerItem{
		LoadBalance: lb,
		ServiceName: service.Info.ServiceName,
	}
	return lb, nil
}

var TransporterHandler *Transporter

// Transporter 为每一个服务创建负载均衡管理器
type Transporter struct {
	TransporterMap   map[string]*TransporterItem
	TransporterSlice []*TransporterItem
	Locker           sync.RWMutex
}

// TransporterItem 加一个名字好找
type TransporterItem struct {
	Trans       *http.Transport
	ServiceName string
}

func init() {
	TransporterHandler = &Transporter{
		TransporterMap:   map[string]*TransporterItem{},
		TransporterSlice: []*TransporterItem{},
		Locker:           sync.RWMutex{},
	}
}

// GetTrans 为每个单独的服务创建
func (t *Transporter) GetTrans(service *ServiceDetail) (*http.Transport, error) {
	for _, transItem := range t.TransporterSlice {
		if transItem.ServiceName == service.Info.ServiceName {
			return transItem.Trans, nil
		}
	}

	//todo 优化点5
	if service.LoadBalance.UpstreamConnectTimeout == 0 {
		service.LoadBalance.UpstreamConnectTimeout = 30
	}
	if service.LoadBalance.UpstreamMaxIdle == 0 {
		service.LoadBalance.UpstreamMaxIdle = 100
	}
	if service.LoadBalance.UpstreamIdleTimeout == 0 {
		service.LoadBalance.UpstreamIdleTimeout = 90
	}
	if service.LoadBalance.UpstreamHeaderTimeout == 0 {
		service.LoadBalance.UpstreamHeaderTimeout = 30
	}
	trans := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(service.LoadBalance.UpstreamConnectTimeout) * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          service.LoadBalance.UpstreamMaxIdle,
		IdleConnTimeout:       time.Duration(service.LoadBalance.UpstreamIdleTimeout) * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: time.Duration(service.LoadBalance.UpstreamHeaderTimeout) * time.Second,
	}

	//save to map and slice
	transItem := &TransporterItem{
		Trans:       trans,
		ServiceName: service.Info.ServiceName,
	}
	t.TransporterSlice = append(t.TransporterSlice, transItem)
	t.Locker.Lock()
	defer t.Locker.Unlock()
	t.TransporterMap[service.Info.ServiceName] = transItem
	return trans, nil
}
