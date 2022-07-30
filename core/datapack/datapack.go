package datapack

import (
	"bytes"
	"camellia/logger"
	pb "camellia/pb_generate"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
)

const (
	FIXED_HEADER_LEN = 15 //fixed header len

)

//TcpPackage tcp data
/*
|mg|flag|header len|pkg len|var header|payload|
|2byte|1byte|4byte|8byte|var|var|

flag
|预留|序列化类型
|4bit|4bit
*/
type TcpPackage struct {
	magic      uint16
	flag       uint8
	headerLen  uint32
	payloadLen uint64
	dataPack   []byte
	//header  []byte
	//payload []byte
}

//func (pkg *TcpPackage) PreReadData() {
//	pkg.header = make([]byte, pkg.headerLen)
//	pkg.payload = make([]byte, pkg.payloadLen)
//}

func (pkg *TcpPackage) MsgLen() uint64 {
	return uint64(pkg.headerLen) + pkg.payloadLen
}

func (pkg *TcpPackage) Pack(data Message) []byte {
	var buf bytes.Buffer
	var err error
	err = binary.Write(&buf, binary.BigEndian, uint16(pb.Constants_Magic))
	if err != nil {
		logger.Fatal("write err", err)
	}

	//flag
	var flag = data.SerializeFlag() // |
	err = binary.Write(&buf, binary.BigEndian, flag)
	if err != nil {
		logger.Fatal("write err", err)
	}

	header := data.SerializeHeader()
	payload := data.SerializePayload()

	err = binary.Write(&buf, binary.BigEndian, int32(len(header)))
	checkErr(err, "write header len err")
	err = binary.Write(&buf, binary.BigEndian, int64(len(payload)))
	checkErr(err, "write payload len err")

	err = binary.Write(&buf, binary.BigEndian, header)
	checkErr(err, "write header err")

	err = binary.Write(&buf, binary.BigEndian, payload)
	checkErr(err, "write payload err")

	return buf.Bytes()
}

//func (pkg *TcpPackage) UnPack(data []byte) {
//	var err error
//	buf := bytes.NewReader(data)
//	err = binary.Read(buf, binary.BigEndian, pkg.magic)
//	checkErr(&err, "write payload err")
//
//	err = binary.Read(buf, binary.BigEndian, pkg.headerLen)
//	checkErr(&err, "write payload err")
//
//	err = binary.Read(buf, binary.BigEndian, pkg.payloadLen)
//	checkErr(&err, "write payload err")
//}

func (pkg *TcpPackage) UnPackFrameHeader(data []byte) error {
	var err error
	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.BigEndian, &pkg.magic)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.BigEndian, &pkg.flag)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.BigEndian, &pkg.headerLen)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.BigEndian, &pkg.payloadLen)
	if err != nil {
		return err
	}
	return nil
}

func (pkg *TcpPackage) UnPackFrameData(data []byte) {
	pkg.dataPack = data

	//var err error
	//buf := bytes.NewReader(data)
	//
	//err = binary.Read(buf, binary.BigEndian, pkg.header)
	//checkErr(err, "write payload err")
	//
	//err = binary.Read(buf, binary.BigEndian, pkg.payload)
	//checkErr(err, "write payload err")
}

//GetMessage TcpPackage --> Message
func (pkg *TcpPackage) GetMessage() Message {
	msg := NewPbMessageHeader()
	msg.Original = pkg.dataPack

	msg.DeserializeHeader(pkg.dataPack[:pkg.headerLen])
	msg.PayloadOffset = pkg.headerLen
	return msg
}

//------------Message-----------

//Message 消息数据
type Message interface {
	SerializeFlag() uint8
	SerializeHeader() []byte
	SerializePayload() []byte

	DeserializeHeader([]byte)

	GetOriginal() []byte
	GetHeader() *pb.Header
	GetPayload() []byte

	//GetMessageHeader()
	//GetMessagePayload()
	//DeserializePayload([]byte)
}

type PbMessage struct {
	Original      []byte
	PayloadOffset uint32

	HeaderPb  *pb.Header
	PayloadPb proto.Message
}

func NewPbMessageHeader() *PbMessage {
	return &PbMessage{
		HeaderPb: &pb.Header{},
	}
}

func NewPbMessage() *PbMessage {
	return &PbMessage{
		HeaderPb: &pb.Header{},
	}
}

func NewPbMessageWithEndpoint(src, dest pb.Endpoint) *PbMessage {
	return &PbMessage{
		HeaderPb: &pb.Header{
			Src:  src,
			Dest: dest,
		},
	}
}

func (m *PbMessage) SerializeFlag() uint8 {
	return uint8(pb.SerializeFlag_PbSerial)
}

func (m *PbMessage) SerializeHeader() []byte {
	h, err := proto.Marshal(m.HeaderPb)
	checkErr(err, "marshal header err")
	return h
}

func (m *PbMessage) SerializePayload() []byte {
	p, err := proto.Marshal(m.PayloadPb)
	checkErr(err, "marshal payload err")
	return p
}

func (m *PbMessage) DeserializeHeader(b []byte) {
	err := proto.Unmarshal(b, m.HeaderPb)
	checkErr(err, "unmarshal header err")
}

func (m *PbMessage) GetHeader() *pb.Header {
	return m.HeaderPb
}

func (m *PbMessage) GetPayload() []byte {
	return m.Original[m.PayloadOffset:]
}

func (m *PbMessage) GetOriginal() []byte {
	return m.Original
}

func checkErr(err error, ifErr string) {
	if err != nil {
		logger.Fatal(ifErr, err)
	}
}
