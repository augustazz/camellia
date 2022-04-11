package main

import (
	"camellia/core"
	"camellia/core/datapack"
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

	c := core.NewConnection(0, &conn)

	go write(c)

	c.ReadLoop()
}

func write(conn *core.Connection) {
	counter := uint64(0)
	for {
		msg := datapack.PbMessage{
			Header: &pb.Header{
				MsgType: pb.MsgType_MsgTypeAuthUp,
				MsgId:   counter,
				Ack: true,
			},
			PayloadPb: &pb.AuthReq{
				Sig: "auth secret sign",
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
