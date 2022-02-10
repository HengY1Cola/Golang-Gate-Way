package dto

import (
	"gin/public"
	"github.com/gin-gonic/gin"
	"time"
)

type AdminInfoOut struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	LoginTime    time.Time `json:"loginTime"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}

type ChangePwdInput struct {
	OldPassWord string `json:"oldPassWord" form:"oldPassWord" comment:"旧密码" example:"123456" validate:"required"`                    // 旧密码
	PassWord    string `json:"password" form:"password" comment:"新密码" example:"Nbxx12345678" validate:"required,isConformPwdFormat"` // 新密码
}

// BindValidParam 绑定到结构体以及校验参数
func (param *ChangePwdInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param) // 传入向下文以及校验的参数
}
