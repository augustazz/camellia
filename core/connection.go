package core

import (
	"camellia/core/channel"
	"camellia/core/datapack"
	"camellia/core/enums"
	"fmt"
	"io"
	"net"
	"time"
)


type Connection struct {
	Id   uint64
	conn *net.Conn

	recvChan   chan datapack.Message
	writeChan chan []byte

	Ctx *channel.ConnContext
}


func NewConnection(id uint64, conn *net.Conn) *Connection {
	c := &Connection{
		Id:        id,
		conn:      conn,
		recvChan:   make(chan datapack.Message, 512),
		writeChan: make(chan []byte, 512),
	}
	c.Ctx = &channel.ConnContext{WriteChan: c.writeChan, State: enums.ConnStateInit}

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

		message := make([]byte, pack.MsgLen())
		_, err = io.ReadFull(*c.conn, message)
		if err != nil {
			fmt.Println("unpack header err", err)
			continue
		}
		pack.UnPackFrameData(message)
		c.recvChan <- pack.GetMessage()
	}

	c.Close()
}

func (c *Connection) startMsgHandler() {
	for {
		select {
		case msg := <-c.recvChan:
			c.Ctx.Head.Fire(c.Ctx, msg)
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
