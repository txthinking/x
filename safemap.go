package ant

import (
    "sync"
)

type SafeMap struct {
    lock *sync.RWMutex
    sm   map[interface{}]interface{}
}

func NewSafeMap() *SafeMap {
    return &SafeMap{
        lock: new(sync.RWMutex),
        sm:   make(map[interface{}]interface{}),
    }
}

func (m *SafeMap) Get(k interface{}) interface{} {
    m.lock.RLock()
    defer m.lock.RUnlock()
    if val, ok := m.sm[k]; ok {
        return val
    }
    return nil
}

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

func (m *SafeMap) IsExists(k interface{}) bool {
    m.lock.RLock()
    defer m.lock.RUnlock()
    if _, ok := m.sm[k]; !ok {
        return false
    }
    return true
}

func (m *SafeMap) Delete(k interface{}) {
    m.lock.Lock()
    defer m.lock.Unlock()
    delete(m.sm, k)
}

