package core

import (
	"camellia/core/channel"
	"camellia/core/datapack"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

var connections = struct {
	cache map[uint64]*Connection
	lock  sync.Mutex
}{
	cache: make(map[uint64]*Connection),
	lock:  sync.Mutex{},
}

type Connection struct {
	Id   uint64
	conn *net.Conn

	recvChan   chan *datapack.TcpPackage
	writeChan chan []byte

	Ctx *channel.ConnContext
}


func NewConnection(id uint64, conn *net.Conn) *Connection {
	c := &Connection{
		Id:        id,
		conn:      conn,
		recvChan:   make(chan *datapack.TcpPackage, 512),
		writeChan: make(chan []byte, 512),
	}
	c.Ctx = &channel.ConnContext{WriteChan: c.writeChan}
	c.Ctx.InitHandlerContext()

	go c.startWriteHandler()
	go c.startMsgHandler()
	return c
}

func (c *Connection) Close() error {
	connections.lock.Lock()
	defer connections.lock.Unlock()

	delete(connections.cache, c.Id)

	fmt.Println("close conn: ", c.Id)
	return (*c.conn).Close()
}

func (c *Connection) ReadLoop() {
	for {
		frameHeader := make([]byte, datapack.FIXED_HEADER_LEN)
		_, err := io.ReadFull(*c.conn, frameHeader)
		if err != nil {
			fmt.Println("read err", err)
			if err == io.EOF {
				break
			}
			continue
		}

		pack := datapack.TcpPackage{}
		err = pack.UnPackFrameHeader(frameHeader)
		if err != nil {
			fmt.Println("unpack header err", err)
			continue
		}

		pack.PreReadData() //初始化header和payload []byte ...
		message := make([]byte, pack.MsgLen())
		_, err = io.ReadFull(*c.conn, message)
		if err != nil {
			fmt.Println("unpack header err", err)
			continue
		}
		pack.UnPackFrameData(message)
		c.recvChan <- &pack
	}

	c.Close()
}

func (c *Connection) startMsgHandler() {
	for {
		select {
		case pkg := <-c.recvChan:
			c.Ctx.Head.Fire(c.Ctx, pkg)
		default:
			time.Sleep(time.Second)
		}
	}
}

func (c *Connection) startWriteHandler() {
	for {
		select {
		case msg := <-c.writeChan:
			n, _ := (*c.conn).Write(msg)
			fmt.Println("send msg byte: ", n)
		}
	}
}



func (c *Connection) Push(msg []byte) {
	c.writeChan<- msg
}

func RegisterToManager(conn *Connection) *Connection {
	connections.lock.Lock()
	defer connections.lock.Unlock()

	old, ok := connections.cache[conn.Id]
	if ok && conn.Id == old.Id {
		return nil
	}
	connections.cache[conn.Id] = conn
	return old
}
