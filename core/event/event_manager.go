package event

import (
	"camellia/logger"
)

var Inst *EventManager

type EventManager struct {
	taskBuffer chan eventTask
	events     map[EventType]func(EventType, interface{}) //事件处理集合 type-func
}

type eventTask struct {
	e      EventType
	f      func(EventType, interface{})
	param  interface{}
	result interface{}
}

func (t eventTask) apply() {
	if t.f != nil {
		t.f(t.e, t.param)
	}
}

func Initialize() *EventManager {
	Inst = &EventManager{
		taskBuffer: make(chan eventTask, 512),
		events:     make(map[EventType]func(EventType, interface{})),
	}

	//注册事件处理函数
	Inst.RegisterEventFunction(map[EventType]func(EventType, interface{}){
		EventTypeConnActive:        ConnActiveEventFunc,
		EventTypeConnStatusChanged: ConnStateChangeEventFunc,
	})

	go func() {
		for task := range Inst.taskBuffer {
			task.apply()
			//task.result
		}
	}()

	return Inst
}

func (m *EventManager) RegisterEventFunction(eventFunc map[EventType]func(EventType, interface{})) {
	for k, v := range eventFunc {
		m.events[k] = v
	}
}

func PostEvent(e EventType, param interface{}) {
	if Inst == nil {
		logger.Warning("event manager instance nil")
		return
	}
	if f, ok := Inst.events[e]; ok {
		t := eventTask{
			e:     e,
			param: param,
			f:     f,
		}
		Inst.taskBuffer <- t
	}
}

func ConnActiveEventFunc(e EventType, f interface{}) {

}

func ConnStateChangeEventFunc(e EventType, f interface{}) {
	if d, ok := f.(ConnStatusChanged); ok {
		logger.Info("conn state changed: ", d)
	}

}
