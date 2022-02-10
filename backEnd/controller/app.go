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

type APPController struct {
}

//APPRegister admin路由注册
func APPRegister(router *gin.RouterGroup) {
	admin := APPController{}
	router.GET("/appList", admin.appList)
	router.GET("/appDetail", admin.appDetail)
	router.GET("/appStat", admin.appStat)
	router.GET("/appDelete", admin.appDelete)
	router.POST("/appAdd", admin.appAdd)
	router.POST("/appUpdate", admin.appUpdate)
}

// appList godoc
// @Summary 租户列表
// @Description 租户列表
// @Tags 租户管理
// @ID /app/appList
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param pageSize query string true "每页多少条"
// @Param pageNo query string true "页码"
// @Success 200 {object} middleware.Response{data=dto.APPListOutput} "success"
// @Router /app/appList [get]
func (admin *APPController) appList(c *gin.Context) {
	// 验证参数
	params := &dto.APPListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3001, err) // 租户列表传入参数错误
		return
	}

	info := &dao.App{}
	list, total, err := info.APPList(c, lib.GORMDefaultPool, params)
	if err != nil {
		middleware.ResponseError(c, 3002, err) // 获取租户列表信息失败
		return
	}

	var outputList []dto.APPListItemOutput
	for _, item := range list {
		appCounter, err := public.FlowCounterHandler.GetCounter(public.FlowAppPrefix + item.AppID)
		if err != nil {
			middleware.ResponseError(c, 3003, err)
			c.Abort()
			return
		}
		outputList = append(outputList, dto.APPListItemOutput{
			ID:       item.ID,
			AppID:    item.AppID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPS: item.WhiteIPS,
			Qpd:      item.Qpd,
			Qps:      item.Qps,
			RealQpd:  appCounter.TotalCount,
			RealQps:  appCounter.QPS,
		})
	}
	output := dto.APPListOutput{
		List:  outputList,
		Total: total,
	}
	middleware.ResponseSuccess(c, output)
	return
}

// appDetail godoc
// @Summary 租户详情
// @Description 租户详情
// @Tags 租户管理
// @ID /app/appDetail
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=dao.App} "success"
// @Router /app/appDetail [get]
func (admin *APPController) appDetail(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3004, err) // 租户ID参数接受错误
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	detail, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 3005, err) // 查询具体某租户
		return
	}
	middleware.ResponseSuccess(c, detail)
	return
}

// appDelete godoc
// @Summary 租户删除
// @Description 租户删除
// @Tags 租户管理
// @ID /app/appDelete
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/appDelete [get]
func (admin *APPController) appDelete(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3004, err)
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 3005, err) // 查询具体某租户
		return
	}
	info.IsDelete = 1 // 进行软删除
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 3006, err)
		return
	}
	middleware.ResponseSuccess(c, "删除成功")
	return
}

// appAdd godoc
// @Summary 租户添加
// @Description 租户添加
// @Tags 租户管理
// @ID /app/appAdd
// @Accept  json
// @Produce  json
// @Param body body dto.APPAddHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/appAdd [post]
func (admin *APPController) appAdd(c *gin.Context) {
	params := &dto.APPAddHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3007, err) // 添加租户信息错误
		return
	}

	//验证appId是否被占用
	search := &dao.App{
		AppID: params.AppID,
	}
	if _, err := search.Find(c, lib.GORMDefaultPool, search); err == nil {
		middleware.ResponseError(c, 3008, errors.New("租户ID被占用，请重新输入"))
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}
	tx := lib.GORMDefaultPool
	info := &dao.App{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPS: params.WhiteIPS,
		Qps:      params.Qps,
		Qpd:      params.Qpd,
	}
	if err := info.Save(c, tx); err != nil {
		middleware.ResponseError(c, 3009, err) // 添加租户信息行为错误
		return
	}
	middleware.ResponseSuccess(c, "添加成功")
	return
}

// appUpdate godoc
// @Summary 租户更新
// @Description 租户更新
// @Tags 租户管理
// @ID /app/appUpdate
// @Accept  json
// @Produce  json
// @Param body body dto.APPUpdateHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/appUpdate [post]
func (admin *APPController) appUpdate(c *gin.Context) {
	params := &dto.APPUpdateHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3010, err) // 接受默认的原租户信息错误
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 3011, err) // 根据提交的ID找不到
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}
	info.Name = params.Name
	info.Secret = params.Secret
	info.WhiteIPS = params.WhiteIPS
	info.Qps = params.Qps
	info.Qpd = params.Qpd
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 3012, err) // 更新租户信息
		return
	}
	middleware.ResponseSuccess(c, "更新成功")
	return
}

// appStat godoc
// @Summary 租户统计
// @Description 租户统计
// @Tags 租户管理
// @ID /app/appStat
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=dto.StatisticsOutput} "success"
// @Router /app/appStat [get]
func (admin *APPController) appStat(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3004, err)
		return
	}

	search := &dao.App{
		ID: params.ID,
	}
	detail, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 3013, err)
		return
	}

	//今日流量全天小时级访问统计
	var todayStat []int64
	counter, err := public.FlowCounterHandler.GetCounter(public.FlowAppPrefix + detail.AppID)
	if err != nil {
		middleware.ResponseError(c, 3014, err)
		c.Abort()
		return
	}
	currentTime := time.Now()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	for i := 0; i <= time.Now().In(loc).Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, loc)
		hourData, _ := counter.GetHourData(dateTime)
		todayStat = append(todayStat, hourData)
	}

	//昨日流量全天小时级访问统计
	var yesterdayStat []int64
	yesterdayTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterdayTime.Year(), yesterdayTime.Month(), yesterdayTime.Day(), i, 0, 0, 0, loc)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayStat = append(yesterdayStat, hourData)
	}
	stat := dto.StatisticsOutput{
		Today:     todayStat,
		Yesterday: yesterdayStat,
	}
	middleware.ResponseSuccess(c, stat)
	return
}
