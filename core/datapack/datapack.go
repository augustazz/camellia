package datapack

import (
	"bytes"
	pb "camellia/pb_generate"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"log"
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

	header  []byte
	payload []byte
}

type Message interface {
	SerializeFlag() uint8
	SerializeHeader() []byte
	SerializePayload() []byte

	DeserializeHeader([]byte)
	DeserializePayload([]byte)
}

type PbMessage struct {
	Header  *pb.Header
	Payload *pb.Payload
}

func NewPbMessage() *PbMessage {
	return &PbMessage{
		Header:  &pb.Header{},
		Payload: &pb.Payload{},
	}
}

func (m *PbMessage) SerializeFlag() uint8 {
	return 0
}

func (m *PbMessage) SerializeHeader() []byte {
	h, err := proto.Marshal(m.Header)
	checkErr(err, "marshal header err")
	return h
}

func (m *PbMessage) SerializePayload() []byte {
	p, err := proto.Marshal(m.Payload)
	checkErr(err, "marshal payload err")
	return p
}

func (m *PbMessage) DeserializeHeader(b []byte) {
	err := proto.Unmarshal(b, m.Header)
	checkErr(err, "unmarshal header err")
}

func (m *PbMessage) DeserializePayload(b []byte) {
	err := proto.Unmarshal(b, m.Payload)
	checkErr(err, "unmarshal payload err")
}

func (pkg *TcpPackage) PreReadData() {
	pkg.header = make([]byte, pkg.headerLen)
	pkg.payload = make([]byte, pkg.payloadLen)
}

func (pkg *TcpPackage) MsgLen() uint64 {
	return uint64(pkg.headerLen) + pkg.payloadLen
}

func (pkg *TcpPackage) Pack(data Message) []byte {
	var buf bytes.Buffer
	var err error
	err = binary.Write(&buf, binary.BigEndian, uint16(pb.Constants_Magic))
	if err != nil {
		log.Fatal("write err", err)
	}

	//flag
	var flag = data.SerializeFlag() // |
	err = binary.Write(&buf, binary.BigEndian, flag)
	if err != nil {
		log.Fatal("write err", err)
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
	pkg.header = data[:pkg.headerLen]
	pkg.payload = data[pkg.headerLen:]
	//var err error
	//buf := bytes.NewReader(data)
	//
	//err = binary.Read(buf, binary.BigEndian, pkg.header)
	//checkErr(err, "write payload err")
	//
	//err = binary.Read(buf, binary.BigEndian, pkg.payload)
	//checkErr(err, "write payload err")
}

func (pkg *TcpPackage) GetMessage() Message {
	//var err error
	//flag := pkg.flag
	//pb
	dp := NewPbMessage()
	dp.DeserializeHeader(pkg.header)
	dp.DeserializePayload(pkg.payload)

	return dp
}

func checkErr(err error, ifErr string) {
	if err != nil {
		log.Fatal(ifErr, err)
	}
}
