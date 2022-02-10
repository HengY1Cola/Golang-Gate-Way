package controller

import (
	"gin/dao"
	"gin/dto"
	"gin/middleware"
	"gin/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"time"
)

type DashboardController struct {
}

func DashboardRegister(group *gin.RouterGroup) {
	dashboard := &DashboardController{}
	group.GET("/panelGroupData", dashboard.PanelGroupData)
	group.GET("/flowStat", dashboard.FlowStat)
	group.GET("/serviceStat", dashboard.ServiceStat)
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 首页大盘
// @ID /main/panelGroupData
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.PanelGroupOutput} "success"
// @Router /main/panelGroupData [get]
func (dashboardController *DashboardController) PanelGroupData(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}
	// 统计ServiceNum
	service := &dao.ServiceInfo{}
	_, serviceTotal, err := service.PageList(c, tx, &dto.ServiceListInput{PageNum: 1, PageSize: 1})
	if err != nil {
		middleware.ResponseError(c, 4001, err) // 查询服务总数错误
		return
	}
	// 统计租户总数
	app := &dao.App{}
	_, appTotal, err := app.APPList(c, tx, &dto.APPListInput{PageSize: 1, PageNo: 1})
	if err != nil {
		middleware.ResponseError(c, 4002, err) // 查询租户总数错误
		return
	}
	// todo 统计当前请求总数
	counter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
	if err != nil {
		middleware.ResponseError(c, 4006, err)
		return
	}
	// 计算今日总的请求
	out := &dto.PanelGroupOutput{
		ServiceNum:      serviceTotal,
		AppNum:          appTotal,
		CurrentNum:      counter.QPS,
		TodayRequestNum: counter.TotalCount,
	}

	// 返回结果
	middleware.ResponseSuccess(c, out)
}

// FlowStat godoc
// @Summary 访问统计
// @Description 访问统计
// @Tags 首页大盘
// @ID /main/flowStat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /main/flowStat [get]
func (dashboardController *DashboardController) FlowStat(c *gin.Context) {
	// 开辟对应的空间
	var todayList []int64
	var yesterday [24]int64
	current := time.Now()
	loc, _ := time.LoadLocation("Asia/Chongqing")

	counter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
	if err != nil {
		middleware.ResponseError(c, 4005, err) // 读取统计错误
		return
	}
	// todo 处理今日
	for i := 0; i <= current.Hour(); i++ {
		newTime := time.Date(current.Year(), current.Month(), current.Day(), i, 0, 0, 0, loc)
		res, _ := counter.GetHourData(newTime)
		todayList = append(todayList, res)
	}
	// todo 处理昨日
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

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /main/serviceStat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.ServiceMainStatOutput} "success"
// @Router /main/serviceStat [get]
func (dashboardController *DashboardController) ServiceStat(c *gin.Context) {
	// 数据库查询
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}
	service := &dao.ServiceInfo{}
	list, err := service.GroupByType(c, tx)
	if err != nil {
		middleware.ResponseError(c, 4003, err) // 获取占比失败
		return
	}
	var Legend []string
	// 取出对用的值
	for index, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			middleware.ResponseError(c, 4004, errors.New("不存在类型")) // 不存在类型
		}
		list[index].Name = name
		Legend = append(Legend, name)
	}

	// 返回结果
	out := &dto.ServiceMainStatOutput{
		Legend: Legend,
		Data:   list,
	}
	middleware.ResponseSuccess(c, out)
}
