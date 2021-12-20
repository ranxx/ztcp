package pack

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"reflect"

	"github.com/ranxx/ztcp/message"
)

// MsgID ...
// type MsgID uint32

// Packer 打包
type Packer interface {
	GetHeadLength() int64
	UnpackHead([]byte) (message.Messager, error)
	Pack(message.Messager) ([]byte, error)
}

type packer struct {
	msg reflect.Type
}

// DefaultPacker packer
func DefaultPacker(msgType message.Messager) Packer {
	// reflect.New(reflect.TypeOf(msgType))
	return &packer{msg: reflect.TypeOf(msgType).Elem()}
}

func (p *packer) GetHeadLength() int64 {
	return 8
}

func (p *packer) UnpackHead(body []byte) (message.Messager, error) {
	reader := bytes.NewBuffer(body)
	msg := reflect.New(p.msg).Interface().(message.Messager)

	msgid := msg.GetMsgID()
	length := msg.GetDataLength()

	// 解析头
	if err := binary.Read(reader, binary.BigEndian, &msgid); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	msg.SetMsgID(msgid)
	msg.SetDataLength(length)
	return msg, nil
}

func (p *packer) Pack(msg message.Messager) ([]byte, error) {
	writer := bytes.NewBuffer(make([]byte, 0, int(p.GetHeadLength()+int64(msg.GetDataLength()))))
	if err := binary.Write(writer, binary.BigEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	if err := binary.Write(writer, binary.BigEndian, msg.GetDataLength()); err != nil {
		return nil, err
	}
	if err := binary.Write(writer, binary.BigEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

// Package ...
type Package struct {
	MsgID  message.MsgID // 消息号
	Length uint32        // 数据部分的 长度
	Msg    []byte        // 数据部分
}

// HeadLength 消息头长度
func HeadLength() int64 {
	return 4 + 4
}

// NewPackage ...
func NewPackage(msgid message.MsgID, msg []byte) *Package {
	pack := Package{
		MsgID:  msgid,
		Length: uint32(len(msg)),
		Msg:    msg,
	}
	return &pack
}

// Empty 空pack
func Empty() *Package {
	return NewPackage(0, nil)
}

// UnpackHeadBytes 解析 head
func (p *Package) UnpackHeadBytes(body []byte) error {
	reader := bytes.NewBuffer(body)

	// 解析头
	if err := binary.Read(reader, binary.BigEndian, &p.MsgID); err != nil {
		return err
	}

	if err := binary.Read(reader, binary.BigEndian, &p.Length); err != nil {
		return err
	}
	return nil
}

// PackBytes ...
func (p *Package) PackBytes() ([]byte, error) {
	l := int(p.Length) + int(HeadLength())
	buffer := bytes.NewBuffer(make([]byte, 0, l))
	err := p.Pack(buffer)
	return buffer.Bytes(), err
}

// Pack ...
func (p *Package) Pack(writer io.Writer) error {
	if err := binary.Write(writer, binary.BigEndian, &p.MsgID); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.BigEndian, &p.Length); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.BigEndian, &p.Msg); err != nil {
		return err
	}
	return nil
}

// ReadHead ...
func (p *Package) ReadHead(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &p.MsgID); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.BigEndian, &p.Length); err != nil {
		return err
	}
	return nil
}

// // UnpackBytes ...
// func (p *Package) UnpackBytes(body []byte) error {

// 	if reader.Len() < int(p.Length) {
// 		return nil
// 	}

// 	if err := binary.Read(reader, binary.BigEndian, &p.Msg); err != nil {
// 		return err
// 	}

// 	return nil
// }

// // Unpack ...
// func (p *Package) Unpack(reader io.Reader) error {
// 	var err error
// 	err = binary.Read(reader, binary.BigEndian, &p.MsgID)
// 	err = binary.Read(reader, binary.BigEndian, &p.Length)
// 	p.Msg = make([]byte, p.Length)
// 	err = binary.Read(reader, binary.BigEndian, &p.Msg)
// 	return err
// }

// // Reset 重置
// func (p *Package) Reset(msg []byte) *Package {
// 	p.Msg = msg
// 	return p
// }

// // PreLength 前置长度
// func (p *Package) PreLength() int64 {
// 	return 4 + 4
// }

// // ReadLast ...
// func (p *Package) ReadLast(reader io.Reader) error {
// 	var err error
// 	p.Msg = make([]byte, p.Length)
// 	err = binary.Read(reader, binary.BigEndian, &p.Msg)
// 	return err
// }

// func absInt64(in int64) int64 {
// 	if in > 0 {
// 		return in
// 	}
// 	return in * -1
// }

// // IsPackage 是否为package数据
// func (p *Package) IsPackage(inerval time.Duration) bool {
// 	return true
// }

// // UnpackBytes ...
// func (p *Package) UnpackBytes(body []byte) error {
// 	buffer := bytes.NewBuffer(body)
// 	return p.Unpack(buffer)
// }

// SplitFunc split
func SplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// log.Println(atEOF, len(data))
	pkg := Package{}
	if atEOF || len(data) < int(HeadLength()) {
		return
	}

	if err = pkg.UnpackHeadBytes(data); err != nil {
		return
	}

	totalLen := int64(pkg.Length) + HeadLength()
	if totalLen > int64(len(data)) {
		return
	}
	return int(totalLen), data[:totalLen], nil
}

// SplitFuncEdgeTriggered ...
func SplitFuncEdgeTriggered(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// log.Println("client.tcp", atEOF, len(data))
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if !atEOF {
		return len(data), data[:], nil
	}
	return 0, nil, nil
}

// NewScanner ...
func NewScanner(reader io.Reader, split bufio.SplitFunc) *bufio.Scanner {
	scanner := bufio.NewScanner(reader)
	scanner.Split(split)
	scanner.Buffer(nil, 1024*1024*1024)
	return scanner
}
