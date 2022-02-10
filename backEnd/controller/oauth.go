package controller

import (
	"encoding/base64"
	"gin/dao"
	"gin/dto"
	"gin/middleware"
	"gin/public"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
	"time"
)

// OAuthController 创建管理员登录结构体
type OAuthController struct {
}

func OAuthRegister(group *gin.RouterGroup) {
	oAuth := &OAuthController{}
	group.POST("/tokens", oAuth.Tokens)
}

// Tokens godoc
// @Summary 获取Token
// @Description 获取Token
// @Tags OAuth
// @ID /oauth/tokens
// @Accept  json
// @Produce  json
// @Param body body dto.TokensInput true "body"
// @Success 200 {object} middleware.Response{data=dto.TokensOutput} "success"
// @Router /oauth/tokens [post]
func (OAuthController *OAuthController) Tokens(c *gin.Context) {
	// 拿到参数并且校验
	params := &dto.TokensInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 5001, err)
		return
	}

	// todo 获取Authorization
	splits := strings.Split(c.GetHeader("Authorization"), " ") // 拿到Header
	if len(splits) != 2 {
		middleware.ResponseError(c, 5001, errors.New("格式错误"))
		return
	}
	decodeString, err := base64.StdEncoding.DecodeString(splits[1]) // 解密
	if err != nil {
		middleware.ResponseError(c, 5002, err)
		return
	}
	splits = strings.Split(string(decodeString), ":")
	appId, secret := splits[0], splits[1]
	// todo 生成appList
	appList := dao.AppManagerHandler.GetAppList()
	// todo 匹配信息 匹配到则jwt加密
	for _, each := range appList {
		if each.AppID == appId && each.Secret == secret {
			loc, _ := time.LoadLocation("Asia/Chongqing")
			claims := jwt.StandardClaims{
				ExpiresAt: time.Now().Add(public.JwtExpires * time.Second).In(loc).Unix(),
				Issuer:    each.AppID,
			}
			token, err := public.JwtEncode(claims)
			if err != nil {
				middleware.ResponseError(c, 5003, err)
				return
			}
			output := &dto.TokensOutput{
				AccessToken: token,
				ExpiresIn:   public.JwtExpires,
				TokenType:   "Bearer",
				Scope:       "read_write",
			}
			middleware.ResponseSuccess(c, output)
			return
		}
	}
	middleware.ResponseError(c, 5004, errors.New("不存在信息"))

	return
}
