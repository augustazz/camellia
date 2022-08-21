package main

import (
	"context"
	"github.com/augustazz/camellia/client/golang/handler"
	"github.com/augustazz/camellia/config"
	"github.com/augustazz/camellia/constants"
	"github.com/augustazz/camellia/core"
	"github.com/augustazz/camellia/core/datapack"
	"github.com/augustazz/camellia/core/event"
	"github.com/augustazz/camellia/logger"
	pb "github.com/augustazz/camellia/pb_generate"
	"net"
	"time"
)

func main() {
	//init logger
	//init logger
	//setup logger
	ctx := context.Background()

	conf := config.LogConfig{
		Debug: true,
		Path:  "./logs/client",
	}
	logger.SetupLogger(ctx, "camellia-client", conf)

	logger.Info("start tcp dial")

	conn, err := net.Dial("tcp4", "127.0.0.1:9090")
	if err != nil {
		logger.Fatal(err)
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
		if conn.Ctx.State == constants.ConnStateClosed {
			break
		}
		if conn.Ctx.State != constants.ConnStateReady {
			time.Sleep(time.Second)
			continue
		}
		msg := datapack.PbMessage{
			HeaderPb: &pb.Header{
				MsgType: pb.MsgType_PropUpload,
				Src:     pb.Endpoint_Client,
				Dest:    pb.Endpoint_ServerThing,
				MsgId:   counter,
				UserInfo: &pb.UserInfo{
					Uid: "100023",
					Did: "DT39485",
				},
				Ack: true,
			},
			PayloadPb: &pb.PropUploadMessage{
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
