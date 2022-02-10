package controller

import (
	"fmt"
	"gin/dao"
	"gin/dto"
	"gin/middleware"
	"gin/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"strings"
	"time"
)

type ServiceController struct {
}

// ServiceRegister 总接口负责与路由挂钩
func ServiceRegister(group *gin.RouterGroup) {
	serviceRegister := &ServiceController{}
	group.GET("/serviceList", serviceRegister.ServiceList)
	group.GET("/serviceDelete", serviceRegister.serviceDelete)
	group.GET("/serviceStat", serviceRegister.serviceStat)
	group.GET("/serviceDetail", serviceRegister.serviceDetail)

	group.POST("/serviceAddHttp", serviceRegister.serviceAddHttp)
	group.POST("/serviceUpdateHttp", serviceRegister.serviceUpdateHttp)

	group.POST("/serviceAddGrpc", serviceRegister.serviceAddGrpc)
	group.POST("/serviceUpdateGrpc", serviceRegister.serviceUpdateGrpc)

	group.POST("/serviceAddTcp", serviceRegister.serviceAddTcp)
	group.POST("/serviceUpdateTcp", serviceRegister.serviceUpdateTcp)
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/serviceList
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param pageNum query int true "页数"
// @Param pageSize query int true "每页数量"
// @Success 200 {object} middleware.Response{data=dto.ServiceOutput} "success"
// @Router /service/serviceList [get]
func (serviceController *ServiceController) ServiceList(c *gin.Context) {
	// 拿到参数并且校验
	params := &dto.ServiceListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err) // 服务列表参数错误
		return
	}

	// 开始查询
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}
	serviceInfo := &dao.ServiceInfo{}
	list, total, err := serviceInfo.PageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err) // 获取分页信息错误
		return
	}

	// 根据查询结果构建List
	var outputList []dto.ServiceItemOutput
	for _, listItem := range list {
		// 服务详细信息
		serviceDetail, err := listItem.ServiceDetail(c, tx, &listItem)
		if err != nil {
			middleware.ResponseError(c, 2003, err) // 联查表格组装错误
		}

		// 开始组装
		// http ip + port / 域名
		// tcp grpc ip + port
		serviceAddr := "unKnow"
		clusterIp := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")
		// 针对https
		if serviceDetail.Info.LoadType == public.LoadTypeHttp && serviceDetail.Http.NeedHttps == 1 {
			if serviceDetail.Http.RuleType == public.HttpRuleTypePrefix {
				serviceAddr = clusterIp + ":" + clusterSSLPort + serviceDetail.Http.Rule
			}
			if serviceDetail.Http.RuleType == public.HttpRuleTypeDomain {
				serviceAddr = serviceDetail.Http.Rule
			}
		} else if serviceDetail.Info.LoadType == public.LoadTypeHttp { // 针对http
			if serviceDetail.Http.RuleType == public.HttpRuleTypePrefix {
				serviceAddr = clusterIp + ":" + clusterPort + serviceDetail.Http.Rule
			}
			if serviceDetail.Http.RuleType == public.HttpRuleTypeDomain {
				serviceAddr = serviceDetail.Http.Rule
			}
		}
		// 针对Tcp
		if serviceDetail.Info.LoadType == public.LoadTypeTcp {
			serviceAddr = fmt.Sprintf("%v:%v", clusterIp, serviceDetail.Tcp.Port)
		}
		// 针对Grpc
		if serviceDetail.Info.LoadType == public.LoadTypeGrpc {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIp, serviceDetail.Grpc.Port)
		}

		// 拿到IpList
		ipList := serviceDetail.LoadBalance.GetIPListByModel()

		// todo 拿到QPS等数据
		serviceCounter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + listItem.ServiceName)
		if err != nil {
			middleware.ResponseError(c, 2048, err)
			return
		}
		// 包装输出
		outItem := dto.ServiceItemOutput{
			ID:          int64(listItem.Id),
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			LoadType:    listItem.LoadType,
			Qps:         serviceCounter.QPS,
			Qpd:         serviceCounter.TotalCount,
			TotalNode:   len(ipList),
		}
		outputList = append(outputList, outItem)
	}

	// 包装成Output
	out := &dto.ServiceOutput{
		Total: total,
		List:  outputList,
	}
	middleware.ResponseSuccess(c, out)
}

// serviceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/serviceDelete
// @Accept  json
// @Produce  json
// @Param id query int true "服务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/serviceDelete [get]
func (serviceController *ServiceController) serviceDelete(c *gin.Context) {
	// 拿到参数并且校验
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2004, err) // 服务删除ID不正确
		return
	}

	// 开始查询
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}
	serviceInfo := &dao.ServiceInfo{Id: params.Id}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2005, err) // 根据ID查询信息失败
		return
	}
	serviceInfo.IsDelete = 1 // 更改为软删除
	if err := serviceInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2006, err) // 保存为软删除错误
		return
	}

	// 返回结果
	middleware.ResponseSuccess(c, "删除成功")
}

// serviceAddHttp godoc
// @Summary 添加HTTP服务
// @Description 添加HTTP服务
// @Tags 服务管理
// @ID /service/serviceAddHttp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/serviceAddHttp [post]
func (serviceController *ServiceController) serviceAddHttp(c *gin.Context) {
	// 拿到参数并且校验
	params := &dto.ServiceAddHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2007, err) // 添加HTTP参数不正确
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2023, errors.New("ip列表与权重设置不匹配"))
		return
	}

	// 查询数据库并开启事务
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}
	tx = tx.Begin()
	// 判断服务名称是否被占用
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, errors.New("存在服务名称")) // 存在服务
		return
	}
	// 查询域名/前缀是否被占用
	httpRule := &dao.HttpRule{RuleType: params.RuleType, Rule: params.Rule}
	httpRule, err = httpRule.Find(c, tx, httpRule)
	if err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, errors.New("接入前缀/域名存在")) // 接入前缀/域名存在
		return
	}
	// 校验字段关联性
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		tx.Rollback()
		middleware.ResponseError(c, 2010, errors.New("IP列表与权重数量不一致")) // 接入前缀/域名存在
		return
	}

	// 开始添加
	// 添加总的Info表格
	service := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	err = service.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2011, err) // 添加总的Info数据失败
		return
	}
	// 添加到http表格
	http := &dao.HttpRule{
		ServiceID:      int64(service.Id), // 保存了之后可以直接拿到
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.NeedHttps,
		NeedWebsocket:  params.NeedWebsocket,
		NeedStripUri:   params.NeedStripUrl,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}
	err = http.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2012, err) // 添加http数据失败
		return
	}
	// 添加到权限控制这张表
	access := &dao.AccessControl{
		ServiceID:         int64(service.Id),
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIPFlowLimit: params.ClientIpFlowLimit,
		ServiceFlowLimit:  params.ServiceIpFlowLimit,
	}
	err = access.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2013, err) // 添加到权限控制数据失败
		return
	}
	// 添加到负载均衡这张表格
	loadBalance := &dao.LoadBalance{
		ServiceID:              int64(service.Id),
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}
	err = loadBalance.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2014, err) // 添加到负载均衡数据失败
		return
	}
	tx.Commit()
	// 返回结果
	middleware.ResponseSuccess(c, "添加成功")
}

// serviceUpdateHttp godoc
// @Summary 更新HTTP服务
// @Description 更新HTTP服务
// @Tags 服务管理
// @ID /service/serviceUpdateHttp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/serviceUpdateHttp [post]
func (serviceController *ServiceController) serviceUpdateHttp(c *gin.Context) {
	// 拿到参数并且校验
	params := &dto.ServiceUpdateHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2015, err) // 更新HTTP参数不正确
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2023, errors.New("ip列表与权重设置不匹配"))
		return
	}

	// 数据库并开启事务
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}
	// 校验字段关联性
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		middleware.ResponseError(c, 2010, errors.New("IP列表与权重数量不一致")) // 接入前缀/域名存在
		return
	}
	tx = tx.Begin()
	// 拿到总的信息
	serviceInfo := &dao.ServiceInfo{Id: params.Id}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo) // 先将基本信息补全
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2022, errors.New("基本信息不存在")) // 基本信息不存在
		return
	}
	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2016, errors.New("服务不存在")) // 更新的服务并不存在
		return
	}
	if serviceDetail.Info.LoadType != public.LoadTypeHttp {
		middleware.ResponseError(c, 2046, errors.New("类型错误")) // tcp不存在该服务
		return
	}
	// 开始更新每张表
	// 更新Info
	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2023, errors.New("更新Info表格失败")) // 更新Info表格失败
		return
	}
	// 更新http
	httpRule := serviceDetail.Http
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedStripUri = params.NeedStripUrl
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfor = params.HeaderTransfor
	err = httpRule.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2017, errors.New("更新HTTP表格失败")) // 更新HTTP表格失败
		return
	}
	// 更新权限表格
	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientIPFlowLimit = params.ClientIpFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceIpFlowLimit
	err = accessControl.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2018, errors.New("更新权限表格失败")) // 更新权限表格失败
		return
	}
	// 更新负载均衡
	loadBalance := serviceDetail.LoadBalance
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadBalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadBalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	loadBalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	err = loadBalance.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2019, errors.New("更新负载均衡表格失败")) // 更新负载均衡表格失败
		return
	}
	tx.Commit()
	// 返回结果
	middleware.ResponseSuccess(c, "更新成功")
}

// serviceDetail godoc
// @Summary 服务详情
// @Description 服务详情
// @Tags 服务管理
// @ID /service/serviceDetail
// @Accept  json
// @Produce  json
// @Param id query int true "服务ID"
// @Success 200 {object} middleware.Response{data=dao.ServiceDetail} "success"
// @Router /service/serviceDetail [get]
func (serviceController *ServiceController) serviceDetail(c *gin.Context) {
	// 拿到参数并且校验
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2020, err) // 获取服务ID信息参数错误
		return
	}

	// 开始查询
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}
	serviceInfo := &dao.ServiceInfo{Id: params.Id}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2005, err) // 根据ID查询信息失败
		return
	}
	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2021, err) // 获取serviceDetail信息失败
		return
	}

	// 返回结果
	middleware.ResponseSuccess(c, serviceDetail)
}

// serviceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 服务管理
// @ID /service/serviceStat
// @Accept  json
// @Produce  json
// @Param id query int true "服务ID"
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /service/serviceStat [get]
func (serviceController *ServiceController) serviceStat(c *gin.Context) {
	// 拿到参数并且校验
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2020, err) // 获取服务ID信息参数错误
		return
	}

	// 开始查询
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}
	serviceInfo := &dao.ServiceInfo{Id: params.Id}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2005, err) // 根据ID查询信息失败
		return
	}
	_, err = serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2021, err) // 获取serviceDetail信息失败
		return
	}

	// 开辟对应的空间
	var todayList []int64
	var yesterday [24]int64
	counter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + serviceInfo.ServiceName)
	if err != nil {
		middleware.ResponseError(c, 2048, err) // 读取统计错误
		return
	}
	current := time.Now()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	for i := 0; i <= current.Hour(); i++ {
		newTime := time.Date(current.Year(), current.Month(), current.Day(), i, 0, 0, 0, loc)
		res, _ := counter.GetHourData(newTime)
		todayList = append(todayList, res)
	}
	yesterdayTime := current.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i < 24; i++ {
		oldTime := time.Date(yesterdayTime.Year(), yesterdayTime.Month(), yesterdayTime.Day(), i, 0, 0, 0, loc)
		res, _ := counter.GetHourData(oldTime)
		yesterday[i] = res
	}
	// 返回结果
	middleware.ResponseSuccess(c, &dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterday,
	})
}

// serviceAddGrpc godoc
// @Summary 添加GRPC服务
// @Description 添加GRPC服务
// @Tags 服务管理
// @ID /service/serviceAddGrpc
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddGrpcInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/serviceAddGrpc [post]
func (serviceController *ServiceController) serviceAddGrpc(c *gin.Context) {
	// 拿到参数并且校验
	params := &dto.ServiceAddGrpcInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2028, err) // 添加GRPC参数不正确
		return
	}

	// 查询数据库并开启事务
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}

	// 判断服务名称是否被占用
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err == nil {
		middleware.ResponseError(c, 2008, errors.New("存在服务名称")) // 存在服务
		return
	}
	// 判断端口是否被占用
	tcpRuleSearch := &dao.TcpRule{
		Port: params.Port,
	}
	_, err = tcpRuleSearch.Find(c, tx, tcpRuleSearch)
	if err == nil {
		middleware.ResponseError(c, 2022, errors.New("端口被占用")) // 端口被占用
		return
	}
	grpcRuleSearch := &dao.GrpcRule{
		Port: params.Port,
	}
	_, err = grpcRuleSearch.Find(c, tx, grpcRuleSearch)
	if err == nil {
		middleware.ResponseError(c, 2022, errors.New("端口被占用")) // 端口被占用
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2023, errors.New("ip列表与权重设置不匹配"))
		return
	}

	// 开始添加
	tx = tx.Begin()
	info := &dao.ServiceInfo{
		LoadType:    public.LoadTypeGrpc,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2024, err) // grpc入Info失败
		return
	}

	grpcRule := &dao.GrpcRule{
		ServiceID:      int64(info.Id),
		Port:           params.Port,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2025, err) // grpc入Grpc分表失败
		return
	}

	loadBalance := &dao.LoadBalance{
		ServiceID:  int64(info.Id),
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2026, err) // grpc入负载均衡失败
		return
	}

	accessControl := &dao.AccessControl{
		ServiceID:         int64(info.Id),
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIpFlowLimit,
		ServiceFlowLimit:  params.ServiceIpFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2027, err) // grpc入权限失败
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "grpc添加成功")
	return
}

// serviceUpdateGrpc godoc
// @Summary grpc服务更新
// @Description grpc服务更新
// @Tags 服务管理
// @ID /service/serviceUpdateGrpc
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateGrpcInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/serviceUpdateGrpc [post]
func (serviceController *ServiceController) serviceUpdateGrpc(c *gin.Context) {
	params := &dto.ServiceUpdateGrpcInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2029, err) // 更新GRPC参数不正确
		return
	}

	// 查询数据库并开启事务
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2023, errors.New("ip列表与权重设置不匹配"))
		return
	}

	// 拿到信息
	service := &dao.ServiceInfo{
		Id: params.Id,
	}
	service, err = service.Find(c, tx, service)
	if service.LoadType != public.LoadTypeGrpc {
		middleware.ResponseError(c, 2046, errors.New("类型错误")) // 不存在该服务
		return
	}
	detail, err := service.ServiceDetail(c, lib.GORMDefaultPool, service)
	if err != nil {
		middleware.ResponseError(c, 2030, err)
		return
	}

	tx = tx.Begin()
	// 更新Info表
	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	info.ServiceName = params.ServiceName
	info.LoadType = public.LoadTypeGrpc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2031, err) // grpc更新Info表错误
		return
	}
	// 更新负载均衡
	loadBalance := &dao.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = int64(info.Id)
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2032, err) // grpc更新负载均衡错误
		return
	}
	// 更新grpc分表
	grpcRule := &dao.GrpcRule{}
	if detail.Grpc != nil {
		grpcRule = detail.Grpc
	}
	grpcRule.ServiceID = int64(info.Id)
	grpcRule.Port = params.Port
	grpcRule.HeaderTransfor = params.HeaderTransfor
	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2033, err) // 更新grpc分表错误
		return
	}
	// 更新权限
	accessControl := &dao.AccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = int64(info.Id)
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIpFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceIpFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2034, err) // 更新grpc权限错误
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "更新成功")
	return
}

// serviceAddTcp godoc
// @Summary 添加Tcp服务
// @Description 添加Tcp服务
// @Tags 服务管理
// @ID /service/serviceAddTcp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddTcpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/serviceAddTcp [post]
func (serviceController *ServiceController) serviceAddTcp(c *gin.Context) {
	// 拿到参数并且校验
	params := &dto.ServiceAddTcpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2035, err) // 添加TCP参数不正确
		return
	}

	// 查询数据库并开启事务
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}

	// 判断服务名称是否被占用
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err == nil {
		middleware.ResponseError(c, 2008, errors.New("存在服务名称")) // 存在服务
		return
	}
	// 判断端口是否被占用
	tcpRuleSearch := &dao.TcpRule{
		Port: params.Port,
	}
	_, err = tcpRuleSearch.Find(c, tx, tcpRuleSearch)
	if err == nil {
		middleware.ResponseError(c, 2022, errors.New("端口被占用")) // 端口被占用
		return
	}
	grpcRuleSearch := &dao.GrpcRule{
		Port: params.Port,
	}
	_, err = grpcRuleSearch.Find(c, tx, grpcRuleSearch)
	if err == nil {
		middleware.ResponseError(c, 2022, errors.New("端口被占用")) // 端口被占用
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2023, errors.New("ip列表与权重设置不匹配"))
		return
	}

	// 开始添加
	tx = tx.Begin()
	info := &dao.ServiceInfo{
		LoadType:    public.LoadTypeTcp,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2036, err) // tcp入Info失败
		return
	}
	tcpRule := &dao.TcpRule{
		ServiceID: int64(info.Id),
		Port:      params.Port,
	}
	if err := tcpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2037, err) // tcp入tcp分表失败
		return
	}
	loadBalance := &dao.LoadBalance{
		ServiceID:  int64(info.Id),
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2038, err) // tcp负载均衡失败
		return
	}
	accessControl := &dao.AccessControl{
		ServiceID:         int64(info.Id),
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIpFlowLimit,
		ServiceFlowLimit:  params.ServiceIpFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2039, err) // tcp入权限失败
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "tcp添加成功")
	return
}

// serviceUpdateTcp godoc
// @Summary tcp服务更新
// @Description tcp服务更新
// @Tags 服务管理
// @ID /service/serviceUpdateTcp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateTcpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/serviceUpdateTcp [post]
func (serviceController *ServiceController) serviceUpdateTcp(c *gin.Context) {
	params := &dto.ServiceUpdateTcpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2040, err) // 更新Tcp参数不正确
		return
	}
	// 查询数据库并开启事务
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2023, errors.New("ip列表与权重设置不匹配"))
		return
	}

	// 连接数据库
	service := &dao.ServiceInfo{
		Id: params.Id,
	}
	service, err = service.Find(c, tx, service)
	if err != nil {
		middleware.ResponseError(c, 2041, err) // tcp不存在该服务
		return
	}
	detail, err := service.ServiceDetail(c, lib.GORMDefaultPool, service)
	if err != nil {
		middleware.ResponseError(c, 2047, err)
		return
	}

	// 判断下要修改的是不是tcp
	log.Println(detail)
	if detail.Info.LoadType != public.LoadTypeTcp {
		middleware.ResponseError(c, 2046, errors.New("类型错误")) // tcp不存在该服务
		return
	}
	// 更新Info表
	tx = tx.Begin()
	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2042, err) // tcp更新Info表错误
		return
	}
	// 更新负载均衡
	loadBalance := &dao.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = int64(info.Id)
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2043, err) // tcp更新负载均衡错误
		return
	}
	// 更新tcp分表
	grpcRule := &dao.GrpcRule{}
	if detail.Grpc != nil {
		grpcRule = detail.Grpc
	}
	grpcRule.ServiceID = int64(info.Id)
	grpcRule.Port = params.Port
	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2044, err) // 更新tcp分表错误
		return
	}
	// 更新权限
	accessControl := &dao.AccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = int64(info.Id)
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIpFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceIpFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2045, err) // 更新tcp权限错误
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "更新成功")
	return
}
