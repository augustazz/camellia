package enums

import "errors"

//--------- datapack.Message error------
var (
	//msg validate fail
	MsgValidateErrEmpty         = errors.New("msg is empty")
	MsgValidateErrHeaderEmpty   = errors.New("msg header is empty")
	MsgValidateErrMsgTypeEmpty  = errors.New("msg header msgType is empty")
	MsgValidateErrUserInfoEmpty = errors.New("msg header userInfo is empty")
)

var ()
