package encoding

import (
	"reflect"

	"github.com/ranxx/ztcp/message"
)

// Unmarshaler 反序列化
type Unmarshaler interface {
	Unmarshal(message.MsgID, []byte) (interface{}, error)
}

// Unmarshal 反序列化
type Unmarshal func(message.MsgID, []byte) (interface{}, error)

// Unmarshal unmarshal
func (w Unmarshal) Unmarshal(id message.MsgID, b []byte) (interface{}, error) {
	return w(id, b)
}

// WrapTypeUnmarshal wrap type
// Type 必须为 point
func WrapTypeUnmarshal(Type interface{}, unmarshal func([]byte, interface{}) error) Unmarshaler {
	v := reflect.TypeOf(Type)
	return Unmarshal(func(id message.MsgID, b []byte) (interface{}, error) {
		tmp := reflect.New(v)
		e := unmarshal(b, tmp)
		return tmp, e
	})
}

type unmarshaler struct {
	unmarshaler        Unmarshaler                   // 消息反序列化
	specialUnmarshaler map[message.MsgID]Unmarshaler // 特效消息体 反序列化
}

// NewUnmarshaler unmarshaler
func NewUnmarshaler(unmarshal Unmarshaler, specialUnmarshalers ...map[message.MsgID]Unmarshal) Unmarshaler {
	m := map[message.MsgID]Unmarshaler{}
	for _, specialUnmarshaler := range specialUnmarshalers {
		for k, v := range specialUnmarshaler {
			m[k] = Unmarshal(v)
		}
	}
	return &unmarshaler{
		unmarshaler:        unmarshal,
		specialUnmarshaler: m,
	}
}

func (u *unmarshaler) Unmarshal(t message.MsgID, b []byte) (interface{}, error) {
	unmarshal := u.specialUnmarshaler[t]
	if unmarshal == nil {
		unmarshal = u.unmarshaler
	}
	return unmarshal.Unmarshal(t, b)
}
