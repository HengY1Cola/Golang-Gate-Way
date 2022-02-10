package load_balance

// LoadBalance 定义接口
type LoadBalance interface {
	Add(...string) error
	Get(string) (string, error)
	//后期服务发现补充
	Update()
}
