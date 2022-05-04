package core

import (
	"camellia/core/enums"
	"camellia/core/util"
	"camellia/logger"
	"sync"
	"time"
)

type ConnManager struct {
	connections     map[uint64]*Connection
	lock            sync.Mutex
	stateCheckWheel *util.TimingWheel
}

func NewConnManager() *ConnManager {
	m := &ConnManager{
		connections:     make(map[uint64]*Connection),
		lock:            sync.Mutex{},
		stateCheckWheel: util.NewTimingWheel(10*time.Second, 60),
	}
	m.stateCheckWheel.Start()
	return m
}

func (m *ConnManager) Close() {
	m.stateCheckWheel.Stop()
}

func (m *ConnManager) RegisterToManager(authTimeoutS, readTimeoutM time.Duration, conn *Connection) {
	old := register0(m, conn)
	if old != nil {
		//kick out
		m.RemoveFromManager("kick out", old)
		if err := old.Close("kick out"); err != nil {
			logger.Error("dup conn err, stop old:", old.Id, err.Error())
		}
	}

	m.stateCheckWheel.AfterFunc(authTimeoutS, func() {
		if conn == nil {
			return
		}
		state := conn.Ctx.State
		if state == enums.ConnStateInit || state == enums.ConnStateInAuth {
			conn.Close("auth timeout")
		}
	})
	m.stateCheckWheel.AfterFunc(readTimeoutM, func() { idleScanner(m.stateCheckWheel, readTimeoutM, conn) })
}

func register0(m *ConnManager, conn *Connection) *Connection {
	m.lock.Lock()
	defer m.lock.Unlock()

	old, ok := m.connections[conn.Id]
	if ok && conn.Id == old.Id {
		return nil
	}
	m.connections[conn.Id] = conn
	return old
}

func (m *ConnManager) RemoveFromManager(reason string, conn *Connection) *Connection {
	m.lock.Lock()
	defer m.lock.Unlock()
	if c, ok := m.connections[conn.Id]; ok {
		err := c.Close(reason)
		if err != nil {
			logger.Warningf("close conn [%d], reason:%s, err:%s", conn.Id, reason, err.Error())
		}
		delete(m.connections, conn.Id)
		return c
	}
	return nil
}

func idleScanner(wheel *util.TimingWheel, readTimeoutM time.Duration, c *Connection) {
	if c == nil || c.Ctx.State != enums.ConnStateReady {
		return
	}
	if time.Now().Sub(c.Ctx.LastReadTime) > readTimeoutM {
		c.Close("idle timeout")
	} else {
		wheel.AfterFunc(readTimeoutM, func() { idleScanner(wheel, readTimeoutM, c) })
	}
}
