package dao

import (
	"fmt"
	"gin/dto"
	"gin/public"
	"log"
	"net/http/httptest"
	"strings"
	"sync"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"基本信息"`
	Http          *HttpRule      `json:"http" description:"http_rule"`
	Tcp           *TcpRule       `json:"tcp" description:"tcp_rule"`
	Grpc          *GrpcRule      `json:"grpc" description:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"loadBalance" description:"loadBalance"`
	AccessControl *AccessControl `json:"accessControl" description:"accessControl"`
}

// ------------------------ 为代理模块下的定义 ------------------------

// ServiceMangerHandler 对外暴露拿起来就可以调用对应的方法了
var ServiceMangerHandler *ServiceManger

func init() {
	ServiceMangerHandler = NewServiceManager() // 当使用ServiceMangerHandler时候就会
}

type ServiceManger struct {
	ServiceMap   map[string]*ServiceDetail
	ServiceSlice []*ServiceDetail
	Locker       sync.RWMutex
	init         sync.Once
	err          error
}

// NewServiceManager 暴露一个创造方法
func NewServiceManager() *ServiceManger {
	return &ServiceManger{
		ServiceMap:   map[string]*ServiceDetail{},
		ServiceSlice: []*ServiceDetail{},
		Locker:       sync.RWMutex{},
		init:         sync.Once{},
	}
}

// LoadOnce 一次性加载服务加载到内存
func (s *ServiceManger) LoadOnce() error {
	s.init.Do(func() { // 只会执行一次，内含锁的机制
		serviceInfo := &ServiceInfo{}
		tx, err := lib.GetGormPool("default")                 // 拿到tx
		c, _ := gin.CreateTestContext(httptest.NewRecorder()) // 拿到c
		if err != nil {
			s.err = err
			return
		}
		params := &dto.ServiceListInput{ //  拿到params
			PageSize: 99999, // 去查找的服务数量
			PageNum:  1,
		}
		list, _, err := serviceInfo.PageList(c, tx, params) // list []ServiceInfo
		if err != nil {
			s.err = err
			return
		}
		s.Locker.Lock() // 加锁防止在遍历的时候出现内存溢出（别人同时在写数据库）
		defer s.Locker.Unlock()
		for _, listItem := range list { // 遍历把ServiceMap与ServiceSlice填充上
			tempItem := listItem
			// 服务详细信息
			serviceDetail, err := tempItem.ServiceDetail(c, tx, &tempItem)
			if err != nil {
				s.err = err
			}
			s.ServiceMap[listItem.ServiceName] = serviceDetail
			s.ServiceSlice = append(s.ServiceSlice, serviceDetail)
			log.Println(serviceDetail.Info.ServiceName)
			log.Println(serviceDetail.LoadBalance.GetIPListByModel())
		}
	})
	return s.err
}

// HttpAccessMode Http模式下
func (s *ServiceManger) HttpAccessMode(c *gin.Context) (*ServiceDetail, error) {
	// 1.前缀匹配 /abc => ServiceSlice.rule <= c.Request.URL.Path
	// 2.域名匹配 www.test.com => ServiceSlice.rule <= c.Request.Host
	host := ""
	if strings.Index(c.Request.Host, ":") != -1 {
		host = c.Request.Host[0:strings.Index(c.Request.Host, ":")] // 防止端口污染
	} else {
		host = c.Request.Host
	}
	path := c.Request.URL.Path

	// 开始服务循环匹配
	// * 可以自己在这里添加自己的逻辑
	for _, serviceItem := range s.ServiceSlice {
		if serviceItem.Info.LoadType != public.LoadTypeHttp { // 如果不是Http
			continue
		}
		// 如或是Http中的
		if serviceItem.Http.RuleType == public.HttpRuleTypeDomain { // 进入域名匹配
			// todo 分情况
			// 如果 serviceItem.Http.Rule => xxx.com/service  则会优先匹配到对应的端口下
			// 如果 serviceItem.Http.Rule => xxx.com 全部转到基本端口下
			// 如果 匹配不到 则 再见
			if serviceItem.Http.Rule == fmt.Sprintf("%v/%v", host, strings.Split(path, "/")[1]) {
				return serviceItem, nil
			} else if serviceItem.Http.Rule == host {
				return serviceItem, nil
			}
		}
		if serviceItem.Http.RuleType == public.HttpRuleTypePrefix { // 进入前缀匹配
			pathMatch := "/" + strings.Split(path, "/")[1] // 匹配路由
			if pathMatch == serviceItem.Http.Rule {        // 匹配上了
				return serviceItem, nil
			}
		}
	}
	return nil, errors.New("no matched service")
}

// GetTcpServiceList 获取Tcp服务列表
func (s *ServiceManger) GetTcpServiceList() ([]*ServiceDetail, error) {
	var list []*ServiceDetail
	for _, item := range s.ServiceSlice {
		temp := item
		if temp.Info.LoadType == public.LoadTypeTcp {
			list = append(list, item)
		}
	}
	return list, nil
}

// GetGrpcServiceList 获取Grpc服务列表
func (s *ServiceManger) GetGrpcServiceList() ([]*ServiceDetail, error) {
	var list []*ServiceDetail
	for _, item := range s.ServiceSlice {
		temp := item
		if temp.Info.LoadType == public.LoadTypeGrpc {
			list = append(list, item)
		}
	}
	return list, nil
}
