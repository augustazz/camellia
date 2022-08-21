package server

import (
	"context"
	"fmt"
	"github.com/augustazz/camellia/core"
	"github.com/augustazz/camellia/core/channel"
	"github.com/augustazz/camellia/core/datapack"
	"github.com/augustazz/camellia/core/event"
	"github.com/augustazz/camellia/logger"
	pb "github.com/augustazz/camellia/pb_generate"
	"github.com/augustazz/camellia/util"
	"net"
	"time"
)

var id uint64

type Server struct {
	connManager *core.ConnManager

	Ctx         context.Context
	Port        int
	AuthTimeout time.Duration //认证超时时间 秒
	IdleTimeout time.Duration //连接无数据传输超时时间 分钟
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", s.Port))

	//listener, err := net.ListenTCP("tcp4", &net.TCPAddr{
	//	IP:   net.IPv4(127, 0, 0, 1),
	//	Port: s.Port,
	//})
	checkErr(err, "listen err", true)
	logger.Info("start success, port: ", s.Port)
	event.Initialize()
	s.connManager = core.NewConnManager()

	s.listenConn(listener)
}

func (s *Server) Close() {
	s.connManager.Close()
}

func (s *Server) listenConn(listener net.Listener) {
	for {
		var conn net.Conn
		var err error
		conn, err = listener.Accept()
		checkErr(err, "listen err", false)
		c := newConnection(&conn)
		s.connManager.RegisterToManager(s.AuthTimeout, s.IdleTimeout, c)

		//resp auth msg
		msg := datapack.NewPbMessageWithEndpoint(pb.Endpoint_ServerConnCenter, pb.Endpoint_Client)
		msg.HeaderPb.MsgType = pb.MsgType_AuthLaunch
		randomStr := util.RandBytes(64)
		msg.PayloadPb = &pb.SimpleMessage{
			Content: randomStr,
		}
		c.Ctx.RandomStr = string(randomStr)
		c.Ctx.WriteChan <- (&datapack.TcpPackage{}).Pack(msg)

		go c.ReadLoop()

		id++
	}
}

func newConnection(conn *net.Conn) *core.Connection {
	clientId := id
	c := core.NewConnection(clientId, conn)
	//init and add handlerContext
	c.Ctx.InitHandlerContext(channel.AuthHandlerFunc, channel.DispatchHandlerFunc, channel.AckHandlerFunc)

	//post event
	event.PostEvent(event.EventTypeConnActive, "")

	return c
	//old := c.ConnActive()
	//if old != nil {
	//	logger.Info("dup conn, stop old:", old.Id)
	//	//kick out
	//	if err := old.Close("kick out"); err != nil {
	//		logger.Error("dup conn err, stop old:", old.Id, err.Error())
	//	}
	//}
}

func (s *Server) stop() error {
	return nil
}

func checkErr(err error, ifErr string, throws bool) {
	if err == nil {
		return
	}
	if throws {
		logger.Fatal(ifErr, err.Error())
	}
	logger.Info(ifErr, err)
}
