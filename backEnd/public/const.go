package public

const (
	ValidatorKey        = "ValidatorKey"
	TranslatorKey       = "TranslatorKey"
	AdminSessionInfoKey = "AdminSessionInfoKey"
	LoadTypeHttp        = 0
	LoadTypeTcp         = 1
	LoadTypeGrpc        = 2

	HttpRuleTypePrefix = 0
	HttpRuleTypeDomain = 1

	RedisFlowDayKey  = "flow_day_count"
	RedisFlowHourKey = "flow_hour_count"

	FlowTotal         = "flow_total"    // 全站
	FlowServicePrefix = "flow_service_" // 服务
	FlowAppPrefix     = "flow_app_"     // 用户

	JwtSignKey = "my_sign_key"
	JwtExpires = 60 * 60
)

// LoadTypeMap 定义一个字典
var (
	LoadTypeMap = map[int]string{
		LoadTypeHttp: "HTTP",
		LoadTypeTcp:  "TCP",
		LoadTypeGrpc: "GRPC",
	}
)
