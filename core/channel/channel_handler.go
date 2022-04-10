package channel

import (
	"camellia/core/datapack"
	pb "camellia/pb_generate"
	"fmt"
)

type DataHandler interface {
	Exec(ctx *ConnContext, pkg *datapack.TcpPackage)
}

type AuthDataHandler struct {
}

func (h *AuthDataHandler) Exec(ctx *ConnContext, pkg *datapack.TcpPackage) {

}

type HeadDataHandler struct {}
type TailDataHandler struct {}
type StdDataHandler struct{}


func (h *HeadDataHandler) Exec(ctx *ConnContext, pkg *datapack.TcpPackage) {
	fmt.Println("head in")
}

func (h *TailDataHandler) Exec(ctx *ConnContext, pkg *datapack.TcpPackage) {
	fmt.Println("tail out")
}

func (h *StdDataHandler) Exec(ctx *ConnContext, pkg *datapack.TcpPackage) {
	msg := pkg.GetMessage()

	if msg.SerializeFlag() == 0 {
		pbMsg := msg.(*datapack.PbMessage)
		fmt.Println(pbMsg.Header.String(), pbMsg.Payload.String())
		if !pbMsg.Header.Ack {
			return
		}

		resp := datapack.PbMessage{
			Header: &pb.Header{
				MsgType: 1,
				MsgId:   pbMsg.Header.MsgId,
			},
			Payload: &pb.Payload{
				Payload: []byte("success"),
			},
		}
		pack := datapack.TcpPackage{}

		ctx.WriteChan<- pack.Pack(&resp)
	}
}



