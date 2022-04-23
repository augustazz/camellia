package event

import (
	"camellia/core/enums"
	"fmt"
)

type EventType uint16

const (
	EventTypeConnActive EventType = iota
	EventTypeConnStatusChanged
)

type ConnStatusChanged struct {
	ConnId  uint64
	Current enums.ConnState
	Before  enums.ConnState
	Err     error
	Msg     string
}

func (c *ConnStatusChanged) String() string {
	var errMsg string
	if c.Err != nil {
		errMsg = c.Err.Error()
	}
	return fmt.Sprintf("conn [%d] state change from: %d to %d, msg: %s, err: %s", c.ConnId, c.Before, c.Current, c.Msg, errMsg)
}
