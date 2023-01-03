package utils

import "sync"

type GetSet[value interface{}] struct {
	value value
	mutex sync.RWMutex
}

func (g GetSet[value]) Get() value {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.value
}

func (g GetSet[value]) Set(v value) {
	g.mutex.Lock()
	g.value = v
	g.mutex.Unlock()
}
