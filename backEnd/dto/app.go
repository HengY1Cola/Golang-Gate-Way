package dto

import (
	"gin/public"
	"github.com/gin-gonic/gin"
	"time"
)

// ------------------ 租户列表 ------------------

type APPListInput struct {
	Info     string `json:"info" form:"info" comment:"查找信息" validate:""`                             // 查找信息
	PageSize int    `json:"pageSize" form:"pageSize" comment:"页数" validate:"required,min=1,max=999"` // 页数
	PageNo   int    `json:"pageNo" form:"pageNo" comment:"页码" validate:"required,min=1,max=999"`     // 页码
}

func (params *APPListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type APPListOutput struct {
	List  []APPListItemOutput `json:"list" form:"list" comment:"租户列表"`
	Total int64               `json:"total" form:"total" comment:"租户总数"`
}

type APPListItemOutput struct {
	ID        int64     `json:"id" form:"id" comment:"租户编号"`
	AppID     string    `json:"appId" form:"appId" comment:"租户AppId"`
	Name      string    `json:"name" form:"name" comment:"租户名称"`
	Secret    string    `json:"secret" form:"secret" comment:"密钥"`
	WhiteIPS  string    `json:"whiteIps" form:"whiteIps" comment:"ip白名单，支持前缀匹配"`
	Qpd       int64     `json:"qpd" form:"qpd" comment:"日请求量限制"`
	Qps       int64     `json:"qps" form:"qps" comment:"每秒请求量限制"`
	RealQpd   int64     `json:"realQpd" form:"realQpd" comment:"日请求量限制"`
	RealQps   int64     `json:"realQps" form:"realQps" comment:"每秒请求量限制"`
	UpdatedAt time.Time `json:"createAt" form:"create_at" comment:"添加时间	"`
	CreatedAt time.Time `json:"updateAt" form:"update_at" comment:"更新时间"`
	IsDelete  int8      `json:"isDelete" form:"is_delete" comment:"是否已删除；0：否；1：是"`
}

// ------------------ 租户详情 ------------------

type APPDetailInput struct {
	ID int64 `json:"id" form:"id" comment:"租户ID" validate:"required"`
}

func (params *APPDetailInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// ------------------ 租户添加 ------------------

type APPAddHttpInput struct {
	AppID    string `json:"appId" form:"appId" comment:"租户id" validate:"required"`
	Name     string `json:"name" form:"name" comment:"租户名称" validate:"required"`
	Secret   string `json:"secret" form:"secret" comment:"密钥" validate:""`
	WhiteIPS string `json:"whiteIps" form:"whiteIps" comment:"ip白名单，支持前缀匹配"`
	Qpd      int64  `json:"qpd" form:"qpd" comment:"日请求量限制" validate:""`
	Qps      int64  `json:"qps" form:"qps" comment:"每秒请求量限制" validate:""`
}

func (params *APPAddHttpInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// ------------------ 租户更新 ------------------

type APPUpdateHttpInput struct {
	ID       int64  `json:"id" form:"id" comment:"主键ID" validate:"required"`
	AppID    string `json:"appId" form:"appId" comment:"租户id" validate:""`
	Name     string `json:"name" form:"name"  comment:"租户名称" validate:"required"`
	Secret   string `json:"secret" form:"secret"  comment:"密钥" validate:"required"`
	WhiteIPS string `json:"whiteIps" form:"whiteIps"  comment:"ip白名单，支持前缀匹配"`
	Qpd      int64  `json:"qpd" form:"qpd"  comment:"日请求量限制"`
	Qps      int64  `json:"qps" form:"qps"  comment:"每秒请求量限制"`
}

func (params *APPUpdateHttpInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// ------------------ 租户统计 ------------------

type StatisticsOutput struct {
	Today     []int64 `json:"today" form:"today" comment:"今日统计" validate:"required"`
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日统计" validate:"required"`
}
