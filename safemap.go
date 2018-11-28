package x

// ref: https://gist.github.com/fatih/6206844

import (
	"sync"
)

// SafeMap is concurrent security map
type SafeMap struct {
	lock *sync.RWMutex
	sm   map[interface{}]interface{}
}

// NewSafeMap get a new concurrent security map
func NewSafeMap() *SafeMap {
	return &SafeMap{
		lock: new(sync.RWMutex),
		sm:   make(map[interface{}]interface{}),
	}
}

// Get used to get a value by key
func (m *SafeMap) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.sm[k]; ok {
		return val
	}
	return nil
}

// Set used to set value with key
func (m *SafeMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.sm[k]; !ok {
		m.sm[k] = v
	} else if val != v {
		m.sm[k] = v
	} else {
		return false
	}
	return true
}

// IsExists determine whether k exists
func (m *SafeMap) IsExists(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.sm[k]; !ok {
		return false
	}
	return true
}

// Delete used to delete a key
func (m *SafeMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.sm, k)
}
