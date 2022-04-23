package main

import (
	"camellia/config"
	"camellia/core"
	"camellia/core/datapack"
	"camellia/core/enums"
	"camellia/core/event"
	"camellia/logger"
	pb "camellia/pb_generate"
	"net"
	"time"
)

func main() {
	//init logger
	//init logger
	//setup logger
	conf := config.LogConfig{
		Debug: true,
		Path: "./logs/client",

	}
	logger.SetupLogger("camellia-client", conf)


	logger.Info("start tcp dial")

	conn, err := net.Dial("tcp4", "127.0.0.1:9090")
	if err != nil {
		logger.Fatal(err)
	}

	event.Initialize()
	c := core.NewConnection(0, &conn)
	//init and add handlerContext
	c.Ctx.InitHandlerContext(/*handler.ClientAuthHandlerFunc*/)

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
			HeaderPb: &pb.Header{
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
