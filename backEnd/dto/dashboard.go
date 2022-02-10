package dto

// ------------------ 大盘指标 ------------------

type PanelGroupOutput struct {
	ServiceNum      int64 `json:"serviceNum" form:"serviceNum" comment:"服务总数" `           // 服务总数
	AppNum          int64 `json:"appNum" form:"appNum" comment:"租户总数" `                   // 租户总数
	CurrentNum      int64 `json:"currentNum" form:"currentNum" comment:"实时请求总数" `         // 实时请求总数
	TodayRequestNum int64 `json:"todayRequestNum" form:"todayRequestNum" comment:"今日请求" ` // 今日请求
}

// ------------------ 服务统计占比 ------------------

type ServiceMainStatOutput struct {
	Legend []string   `json:"legend" form:"legend" comment:"总的数据" ` // 总的数据
	Data   []StatItem `json:"data" form:"data" comment:"每份数据" `     // 每份数据
}

type StatItem struct {
	Name     string `json:"name" form:"name" comment:"名称" `
	LoadType int    `json:"load_type" form:"load_type" comment:"对应的类型" `
	Value    int64  `json:"value" form:"value" comment:"占比" `
}
