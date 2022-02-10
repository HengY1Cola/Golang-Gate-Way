package dto

import (
	"gin/public"
	"github.com/gin-gonic/gin"
	"time"
)

// AdminLoginInput 定义验证登录的结构体
// 注意验证器之间不能有空格
// 设置Tag
// json是转出去的结构体 form是转进来的
type AdminLoginInput struct {
	UserName string `json:"username" form:"username" comment:"用户名" example:"admin" validate:"required,isValidateUserName"` // 这里是文档里面的描述
	PassWord string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`                    // 密码
}

// AdminLoginOutput 定义返回结构体
type AdminLoginOutput struct {
	Token string `json:"token" form:"token" comment:"token" example:"token" validate:""` // token
}

// AdminSessionInfo session结构体
type AdminSessionInfo struct {
	ID        int       `json:"id"`
	UserName  string    `json:"userName"`
	LoginTime time.Time `json:"loginTime"`
}

// BindValidParam 绑定到结构体以及校验参数
func (param *AdminLoginInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param) // 传入向下文以及校验的参数
}
