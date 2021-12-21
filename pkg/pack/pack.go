package pack

import (
	"bytes"
	"encoding/binary"

	"github.com/ranxx/ztcp/pkg/message"
)

// Packer 打包
type Packer interface {
	GetHeadLength() int64
	UnpackHead([]byte) (message.Messager, error)
	Pack(message.Messager) ([]byte, error)
}

// DefaultPack 打包
func DefaultPack(msg message.Messager) ([]byte, error) {
	return DefaultPacker(message.DefaultMessager).Pack(msg)
}

type packer struct {
	gen message.GenMessage
}

// DefaultPacker packer
func DefaultPacker(gen message.GenMessage) Packer {
	if gen == nil {
		gen = message.DefaultMessager
	}
	return &packer{gen: gen}
}

func (p *packer) GetHeadLength() int64 {
	return 8
}

func (p *packer) UnpackHead(body []byte) (message.Messager, error) {
	reader := bytes.NewBuffer(body)

	msg := p.gen(message.MsgID(0), nil)

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
