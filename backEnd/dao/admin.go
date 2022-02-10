package dao

import (
	"gin/dto"
	"gin/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

// 在数据库文件中我们已经连接上了数据库了
// 基本逻辑 结构体的声明 数据库表名字 以及查询

type Admin struct {
	Id        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName  string    `json:"user_name" gorm:"column:user_name" description:"管理员用户名"`
	Salt      string    `json:"salt" gorm:"column:salt" description:"盐值"`
	Password  string    `json:"password" gorm:"column:password" description:"密码"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}

// Find 查询方法
func (t *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	admin := &Admin{}
	err := tx.WithContext(c).Where(search).Find(admin).Error
	if err != nil {
		return nil, err
	}
	return admin, nil
}

// Save 更新方法
func (t *Admin) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(t).Error // 把自己本身保存下来
}

// --------------------逻辑区域---------------------

// LoginAndCheck 拿到登录的用户名与密码进行加盐判断
func (t *Admin) LoginAndCheck(c *gin.Context, tx *gorm.DB, params *dto.AdminLoginInput) (*Admin, error) {
	// 拿到用户名admin的信息
	adminInfo, err := t.Find(c, tx, &Admin{
		UserName: params.UserName,
		IsDelete: 0,
	})
	if err != nil {
		return nil, errors.New("用户信息不存在")
	}

	// 拿到盐值与密码进行sha256的加密
	saltPassWord := public.GenSaltPassword(adminInfo.Salt, params.PassWord)

	// 与数据库查询是否相等
	if saltPassWord != adminInfo.Password {
		return nil, errors.New("密码错误，请重新输入")
	}
	return adminInfo, nil
}
