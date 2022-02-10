package public

import (
	"sync"
	"time"
)

var FlowCounterHandler *FlowCounter

type FlowCounter struct {
	RedisFlowCountMap   map[string]*RedisFlowCountService
	RedisFlowCountSlice []*RedisFlowCountService
	Locker              sync.RWMutex
}

func init() {
	FlowCounterHandler = &FlowCounter{
		RedisFlowCountMap:   map[string]*RedisFlowCountService{},
		RedisFlowCountSlice: []*RedisFlowCountService{},
		Locker:              sync.RWMutex{},
	}
}

func (counter *FlowCounter) GetCounter(serverName string) (*RedisFlowCountService, error) {
	// todo 判断存不存在
	for _, item := range counter.RedisFlowCountSlice {
		if item.AppID == serverName {
			return item, nil
		}
	}

	// todo 不存在就创建一个
	newCounter := NewRedisFlowCountService(serverName, 1*time.Second)
	counter.RedisFlowCountSlice = append(counter.RedisFlowCountSlice, newCounter)
	counter.Locker.Lock()
	defer counter.Locker.Unlock()
	counter.RedisFlowCountMap[serverName] = newCounter
	return newCounter, nil
}
