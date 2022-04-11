package channel

import (
	"camellia/core/datapack"
	pb "camellia/pb_generate"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type DataHandler interface {
	Exec(ctx *ConnContext, pkg datapack.Message)
	Match(ctx *ConnContext, pkg datapack.Message) bool
}

type AuthDataHandler struct {
}

func (h *AuthDataHandler) Exec(ctx *ConnContext, pkg datapack.Message) {

}

type HeadDataHandler struct {}
type TailDataHandler struct {}
type StdDataHandler struct{}


func (h *HeadDataHandler) Exec(ctx *ConnContext, msg datapack.Message) {
	fmt.Println("head in")
}

func (h *HeadDataHandler) Match(ctx *ConnContext, msg datapack.Message) bool {
	return true
}

func (h *TailDataHandler) Exec(ctx *ConnContext, msg datapack.Message) {
	fmt.Println("tail out")
}

func (h *TailDataHandler) Match(ctx *ConnContext, msg datapack.Message) bool {
	return true
}

func (h *StdDataHandler) Exec(ctx *ConnContext, msg datapack.Message) {
	pbMsg := msg.(*datapack.PbMessage)


	var payload pb.SimplePayload
	err := proto.Unmarshal(msg.GetPayload(), &payload)
	if err != nil {
		fmt.Println("unmarshal payload fail, ", err)
		return
	}
	fmt.Println(pbMsg.Header.String(), payload.String())

	if !pbMsg.Header.Ack {
		return
	}

	resp := datapack.PbMessage{
		Header: &pb.Header{
			MsgType: 1,
			MsgId:   pbMsg.Header.MsgId,
		},
		PayloadPb: &pb.SimplePayload{
			Payload: []byte("success"),
		},
	}
	pack := datapack.TcpPackage{}

	ctx.WriteChan<- pack.Pack(&resp)
}

func (h *StdDataHandler) Match(ctx *ConnContext, msg datapack.Message) bool {
	return true
}



