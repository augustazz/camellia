package constants

type ConnState uint16

const (
	ConnStateInit ConnState = iota
	ConnStateInAuth
	ConnStateReady
	ConnStateClosed
)
