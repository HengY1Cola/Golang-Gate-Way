package public

import (
	"golang.org/x/time/rate"
	"sync"
)

// FlowLimiterHandler 基于桶实现限流
var FlowLimiterHandler *FlowLimiter

type FlowLimiter struct {
	FlowLimiterMap map[string]*FlowLimiterItem
	FlowLimitSlice []*FlowLimiterItem
	Locker         sync.RWMutex
}

type FlowLimiterItem struct {
	ServiceName string
	Limiter     *rate.Limiter
}

func init() {
	FlowLimiterHandler = &FlowLimiter{
		FlowLimiterMap: map[string]*FlowLimiterItem{},
		FlowLimitSlice: []*FlowLimiterItem{},
		Locker:         sync.RWMutex{},
	}
}

func (counter *FlowLimiter) GetLimiter(serverName string, qps float64) (*rate.Limiter, error) {
	// todo 判断存不存在
	for _, item := range counter.FlowLimitSlice {
		if item.ServiceName == serverName {
			return item.Limiter, nil
		}
	}

	// todo 创建
	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	counter.FlowLimitSlice = append(counter.FlowLimitSlice, &FlowLimiterItem{
		ServiceName: serverName,
		Limiter:     newLimiter,
	})
	counter.Locker.Lock()
	defer counter.Locker.Unlock()
	counter.FlowLimiterMap[serverName] = &FlowLimiterItem{
		ServiceName: serverName,
		Limiter:     newLimiter,
	}
	return newLimiter, nil
}
