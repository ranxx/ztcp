package encoding

import (
	"github.com/ranxx/ztcp/message"
)

// Marshaler 序列化
type Marshaler interface {
	Marshal(message.MsgID, interface{}) ([]byte, error)
}

// Marshal marshal
type Marshal func(message.MsgID, interface{}) ([]byte, error)

// Marshal marshal
func (w Marshal) Marshal(id message.MsgID, v interface{}) ([]byte, error) {
	return w(id, v)
}

// WrapMarshalerWithoutMsgID 包裹 序列化，不包含 msgid
type WrapMarshalerWithoutMsgID func(interface{}) ([]byte, error)

// Marshal 序列化
func (w WrapMarshalerWithoutMsgID) Marshal(id message.MsgID, v interface{}) ([]byte, error) {
	return w(v)
}

type marshaler struct {
	marshal          Marshaler                   // 消息序列化
	specialMarshaler map[message.MsgID]Marshaler // 特效消息体 序列化
}

// NewMarshaler ...
func NewMarshaler(marshal Marshaler, specialMarshalers ...map[message.MsgID]Marshal) Marshaler {
	m := map[message.MsgID]Marshaler{}
	for _, specialMarshaler := range specialMarshalers {
		for k, v := range specialMarshaler {
			m[k] = Marshal(v)
		}
	}
	return &marshaler{
		marshal:          marshal,
		specialMarshaler: m,
	}
}

func (m *marshaler) Marshal(t message.MsgID, v interface{}) ([]byte, error) {
	// TODO: 是否强制marshaler

	// 先判断v的类型
	if b, ok := v.([]byte); ok {
		return b, nil
	}

	// 判断
	marshal := m.specialMarshaler[t]
	if marshal == nil {
		marshal = m.marshal
	}
	return marshal.Marshal(t, v)
}
