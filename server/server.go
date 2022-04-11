package server

import (
	"camellia/core"
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
	old := core.RegisterToManager(c)
	if old != nil {
		log.Println("dup conn, stop old:", old.Id)
		err := old.Close()
		checkErr(err, "close old conn err", false)
	}

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
