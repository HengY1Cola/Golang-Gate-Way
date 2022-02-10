package controller

import (
	"encoding/json"
	"errors"
	"gin/dao"
	"gin/dto"
	"gin/middleware"
	"gin/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AdminController 创建管理员登录结构体
type AdminController struct {
}

// AdminRegister 总接口负责与路由挂钩
func AdminRegister(group *gin.RouterGroup) {
	adminLogin := &AdminController{}
	group.GET("/adminInfo", adminLogin.AdminInfo) // 这个是拿到session返回信息
	group.POST("/changePwd", adminLogin.changePwd)
}

// AdminInfo godoc
// @Summary 获取管理员信息
// @Description 获取管理员信息
// @Tags 管理员接口
// @ID /admin/adminInfo
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOut} "success"
// @Router /admin/adminInfo [get]
func (adminController *AdminController) AdminInfo(c *gin.Context) {
	// 读取session值的获取json转换为结构体
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminInfo := &dto.AdminSessionInfo{}
	err := json.Unmarshal([]byte(sessInfo.(string)), adminInfo)
	if err != nil {
		middleware.ResponseError(c, 1005, err) // json反序列化错误
		return
	}
	// 取出数据封装返回
	out := &dto.AdminInfoOut{
		ID:           adminInfo.ID,
		Name:         adminInfo.UserName,
		LoginTime:    adminInfo.LoginTime,
		Avatar:       "https://inews.gtimg.com/newsapp_bt/0/13392595208/1000",
		Introduction: "",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}

// AdminInfo godoc
// @Summary 更改密码
// @Description 更改密码
// @Tags 管理员接口
// @ID /admin/changePwd
// @Accept  json
// @Produce  json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/changePwd [post]
func (adminController *AdminController) changePwd(c *gin.Context) {
	// 拿到参数并且校验
	params := &dto.ChangePwdInput{} // 到dto中注册输入输出
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 1006, err) // 传入的密码不合格
		return
	}

	// 1. 读取session信息
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminInfo := &dto.AdminSessionInfo{}
	err := json.Unmarshal([]byte(sessInfo.(string)), adminInfo)
	if err != nil {
		middleware.ResponseError(c, 1005, err) // json反序列化错误
		return
	}

	// 2. 利用id来查询验证原来密码是否正确以及是否与原来密码一样
	adminStorageInfo := &dao.Admin{}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}
	adminStorageInfo, err = adminStorageInfo.Find(c, tx, &dao.Admin{
		Id:       adminInfo.ID,
		UserName: adminInfo.UserName,
		IsDelete: 0,
	})
	if err != nil {
		middleware.ResponseError(c, 1007, err) // 根据session查询不到信息
		return
	}
	if adminStorageInfo.Password != public.GenSaltPassword(adminStorageInfo.Salt, params.OldPassWord) {
		middleware.ResponseError(c, 1008, errors.New("与旧的密码不匹配"))
		return
	}
	if adminStorageInfo.Password == public.GenSaltPassword(adminStorageInfo.Salt, params.PassWord) {
		middleware.ResponseError(c, 1009, errors.New("与原来密码相同"))
		return
	}

	// 3. 更改加盐密码
	saltPassword := public.GenSaltPassword(adminStorageInfo.Salt, params.PassWord)
	// 4. 保存为新的密码
	adminStorageInfo.Password = saltPassword
	err = adminStorageInfo.Save(c, tx)
	if err != nil {
		middleware.ResponseError(c, 1010, err) // 新的账号密码保存错误
	} else {
		middleware.ResponseSuccess(c, "更改密码成功")
	}
}
