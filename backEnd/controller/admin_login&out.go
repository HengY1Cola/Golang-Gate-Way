package controller

import (
	"encoding/json"
	"gin/dao"
	"gin/dto"
	"gin/middleware"
	"gin/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

// AdminLoginController 创建管理员登录结构体
type AdminLoginController struct {
}

// AdminLoginRegister 总接口负责与路由挂钩
// 在controller中写一个接口的话接受的上游的Group
func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}       //定义好然后注册到Group中
	group.POST("/login", adminLogin.AdminLogin) // 这个是登录下发session的作用
	group.GET("/logout", adminLogin.AdminLoginOut)
}

// AdminLogin godoc
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (adminLogin *AdminLoginController) AdminLogin(c *gin.Context) { // 传入Gin上下文的指针

	// 拿到参数并且校验
	params := &dto.AdminLoginInput{} // 到dto中注册输入输出
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 1001, err) // 调用中间件
		return
	}

	// 验证的业务逻辑
	// 1. 拿到用户名这里固定死的为admin
	// 2. 拿到盐值与密码进行sha256的加密
	// 3. 与数据库查询是否相等
	tx, err := lib.GetGormPool("default") // 拿到数据库的连接 default为配置文件中的list.default
	if err != nil {
		middleware.ResponseError(c, 10000, err) // 去取数据库连接失败
		return
	}
	admin := &dao.Admin{}
	admin, err = admin.LoginAndCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 1002, err) // 验证错误不通过
		return
	}

	// 进行session设置
	// 1. 将session对应的结构体信息json化
	// 2. 存入到session中
	jsonInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	marshal, err := json.Marshal(jsonInfo)
	if err != nil {
		middleware.ResponseError(c, 1003, err) // json格式的判断
		return
	}
	sess := sessions.Default(c)
	sess.Set(public.AdminSessionInfoKey, string(marshal)) // 设置
	err = sess.Save()
	if err != nil {
		middleware.ResponseError(c, 1004, err) // session写入redis错误
		return
	} // 存储到Redis中

	// 定义输出
	output := &dto.AdminLoginOutput{Token: admin.UserName} // 通过就是admin 没有通过就是空
	middleware.ResponseSuccess(c, output)
}

// AdminLoginOut godoc
// @Summary 管理员退出登录
// @Description 管理员退出登录
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [get]
func (adminLogin *AdminLoginController) AdminLoginOut(c *gin.Context) {
	// 通过连接的上下文删除对应的key
	sess := sessions.Default(c)
	sess.Delete(public.AdminSessionInfoKey)
	sess.Save()

	// 返回结果
	middleware.ResponseSuccess(c, "")
}
