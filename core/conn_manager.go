package core

import (
	"camellia/core/event"
	"sync"
)

var connections = struct {
	cache map[uint64]*Connection
	lock  sync.Mutex
}{
	cache: make(map[uint64]*Connection),
	lock:  sync.Mutex{},
}


func (c *Connection) ConnActive() *Connection {
	old := RegisterToManager(c)
	event.PostEvent(event.EventTypeConnActive, "")
	return old
}


func RegisterToManager(conn *Connection) *Connection {
	connections.lock.Lock()
	defer connections.lock.Unlock()

	old, ok := connections.cache[conn.Id]
	if ok && conn.Id == old.Id {
		return nil
	}
	connections.cache[conn.Id] = conn
	return old
}


