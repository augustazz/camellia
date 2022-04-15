package server

import (
	"camellia/core"
	"camellia/core/channel"
	"camellia/core/datapack"
	"camellia/core/event"
	"camellia/core/util"
	pb "camellia/pb_generate"
	"fmt"
	"log"
	"net"
)

var id uint64

type Server struct {
	Port int
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", s.Port))

	//listener, err := net.ListenTCP("tcp4", &net.TCPAddr{
	//	IP:   net.IPv4(127, 0, 0, 1),
	//	Port: s.Port,
	//})
	checkErr(err, "listen err", true)
	fmt.Println("start success, port: ", s.Port)
	event.Initialize()

	for {
		var conn net.Conn
		conn, err = listener.Accept()
		checkErr(err, "listen err", false)
		dealConn(&conn)
		id++
	}
}

func dealConn(conn *net.Conn) {
	clientId := id
	c := core.NewConnection(clientId, conn)
	//init and add handlerContext
	c.Ctx.InitHandlerContext(channel.AuthHandlerFunc, channel.DispatchHandlerFunc)

	old := c.ConnActive()
	if old != nil {
		log.Println("dup conn, stop old:", old.Id)
		//kick out
		if err := old.Close(); err != nil {
			log.Println("dup conn err, stop old:", old.Id, err.Error())
		}
	}

	//resp auth msg
	msg := datapack.NewPbMessage()
	msg.Header.MsgType = pb.MsgType_MsgTypeAuthLaunch
	s := util.RandBytes(64)
	msg.PayloadPb = &pb.SimplePayload{
		Content: s,
	}
	c.Ctx.RandomStr = string(s)
	c.Ctx.WriteChan <- (&datapack.TcpPackage{}).Pack(msg)

	go c.ReadLoop()
}

func (s *Server) stop() error {
	return nil
}

func checkErr(err error, ifErr string, throws bool) {
	if err == nil {
		return
	}
	if throws {
		log.Fatal(ifErr, err.Error())
	}
	log.Println(ifErr, err)
}
