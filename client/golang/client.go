package main

import (
	"camellia/client/golang/handler"
	"camellia/core"
	"camellia/core/datapack"
	"camellia/core/enums"
	"camellia/core/event"
	pb "camellia/pb_generate"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp4", "127.0.0.1:9090")
	if err != nil {
		log.Fatalln(err)
	}

	event.Initialize()
	c := core.NewConnection(0, &conn)
	//init and add handlerContext
	c.Ctx.InitHandlerContext(handler.ClientAuthHandlerFunc)

	go write(c)

	c.ReadLoop()
}

func write(conn *core.Connection) {
	counter := uint64(0)
	for {
		if conn.Ctx.State == enums.ConnStateClosed {
			break
		}
		if conn.Ctx.State != enums.ConnStateReady {
			time.Sleep(time.Second)
			continue
		}
		msg := datapack.PbMessage{
			Header: &pb.Header{
				MsgType: pb.MsgType_MsgTypePropUpload,
				MsgId:   counter,
				Ack: true,
			},
			PayloadPb: &pb.PropUpload{
				Props: map[string]string{
					"version": "1.0.0",
				},
			},
		}
		pack := datapack.TcpPackage{}
		conn.Push(pack.Pack(&msg))
		counter++
		time.Sleep(time.Second * 2)
	}
}

func checkErr(err error, ifErr string) {
	if err != nil {
		log.Fatal(ifErr, err)
	}
}
