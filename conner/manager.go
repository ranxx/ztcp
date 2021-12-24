package conner

import "sync"

// Manager ...
type Manager struct {
	conns  map[int64]Conner
	rwlock *sync.RWMutex
}

// NewManager ...
func NewManager() *Manager {
	return &Manager{
		conns:  map[int64]Conner{},
		rwlock: new(sync.RWMutex),
	}
}

// AddConn ...
func (m *Manager) AddConn(c Conner) *Manager {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	m.conns[c.ID()] = c
	return m
}

// Get ...
func (m *Manager) Get(id int64) Conner {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()
	return m.conns[id]
}

// Close 关闭
func (m *Manager) Close() {
	m.rwlock.RLock()
	for _, v := range m.conns {
		v.Close()
	}
	m.rwlock.RUnlock()
	// 全部删除
	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	m.conns = map[int64]Conner{}
}

// Del 删除
func (m *Manager) Del(id int64) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	delete(m.conns, id)
}
