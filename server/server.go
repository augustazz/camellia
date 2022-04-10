package server

import (
	"camellia/core"
	"fmt"
	"log"
	"net"
)

var id uint64

type Server struct {
}

func (server *Server) Start() {
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 9090,
	})
	checkErr(err, "listen err", true)

	fmt.Println("start success, port: ", 9090)

	for {
		conn, err := listener.Accept()
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

func (server *Server) stop() error {
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
