package dao

import (
	"gin/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type ServiceInfo struct {
	Id          int       `json:"id" gorm:"primary_key" description:"自增主键"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"负载类型 0=http 1=tcp 2=grpc"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	CreatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	UpdatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete    int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

// Find 查询方法
func (t *ServiceInfo) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	service := &ServiceInfo{}
	err := tx.WithContext(c).Where(search).First(service).Error
	if err != nil {
		return nil, err
	}
	return service, nil
}

// Save 更新方法
func (t *ServiceInfo) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(t).Error // 把自己本身保存下来
}

// --------------------逻辑区域---------------------

// PageList 分页查询方法
func (t *ServiceInfo) PageList(c *gin.Context, tx *gorm.DB, param *dto.ServiceListInput) ([]ServiceInfo, int64, error) {
	var total int64 = 0                            // 定义总数
	var list []ServiceInfo                         // 接受查询结果
	offset := (param.PageNum - 1) * param.PageSize // 定义偏移量

	query := tx.WithContext(c)                              // 新建一个查询
	query = query.Table(t.TableName()).Where("is_delete=0") // 通用的查询
	// 自己构建一个进行模糊查询
	if param.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(param.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}

// ServiceDetail 联查表格组装
func (t *ServiceInfo) ServiceDetail(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	// 挨个查询
	httpRule := &HttpRule{ServiceID: int64(search.Id)}
	httpRule, err := httpRule.Find(c, tx, httpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	tcpRule := &TcpRule{ServiceID: int64(search.Id)}
	tcpRule, err = tcpRule.Find(c, tx, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	grpcRule := &GrpcRule{ServiceID: int64(search.Id)}
	grpcRule, err = grpcRule.Find(c, tx, grpcRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	accessControl := &AccessControl{ServiceID: int64(search.Id)}
	accessControl, err = accessControl.Find(c, tx, accessControl)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	loadBalance := &LoadBalance{ServiceID: int64(search.Id)}
	loadBalance, err = loadBalance.Find(c, tx, loadBalance)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 开始组装
	detail := &ServiceDetail{
		Info:          search,
		Http:          httpRule,
		Tcp:           tcpRule,
		Grpc:          grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	return detail, nil
}

// GroupByType 服务占比
func (t *ServiceInfo) GroupByType(c *gin.Context, tx *gorm.DB) ([]dto.StatItem, error) {
	var whole []dto.StatItem
	query := tx.WithContext(c)
	err := query.Table(t.TableName()).Where("is_delete=0").Select("load_type, count(*) as value").Group("load_type").Scan(&whole).Error
	if err != nil {
		return nil, err
	}
	return whole, nil
}
