package constants

import "errors"

//--------- datapack.Message error------
var (
	//msg validate fail
	MsgValidateErrEmpty         = errors.New("msg is empty")
	MsgValidateErrHeaderEmpty   = errors.New("msg header is empty")
	MsgValidateErrMsgTypeEmpty  = errors.New("msg header msgType is empty")
	MsgValidateErrUserInfoEmpty = errors.New("msg header userInfo is empty")

	MsgProcessorNotFound = errors.New("msg processor not found")
	ServicesNotAvailable = errors.New("services not available")

	ConfigServicesNotFound = errors.New("config services not found")
)
