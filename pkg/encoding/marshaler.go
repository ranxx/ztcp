package encoding

import (
	"encoding/json"
	"fmt"
	"reflect"

	gogo "github.com/gogo/protobuf/proto"
	"github.com/ranxx/ztcp/pkg/message"
)

// Marshaler 序列化
type Marshaler interface {
	// 如果 v 有实现 ProvideMarshaler 方法，则会优先调用 Bytes
	Marshal(id message.MsgID, v interface{}) ([]byte, error)
}

// ProvideMarshaler 提供的 marshaler
type ProvideMarshaler interface {
	Bytes() ([]byte, error)
}

// Marshal marshal
type Marshal func(message.MsgID, interface{}) ([]byte, error)

// Marshal marshal
func (w Marshal) Marshal(id message.MsgID, v interface{}) ([]byte, error) {
	return w(id, v)
}

// DefaultMarshal 支持 ProvideMarshaler 接口
//
// 支持 基础类型(array, slice, struct, map)的 json序列化
//
// 不支持的将会 返回 unknown type
func DefaultMarshal(id message.MsgID, v interface{}) ([]byte, error) {
	if bytes, ok := v.(ProvideMarshaler); ok {
		return bytes.Bytes()
	}

	// gogo proto
	if gproto, ok := v.(gogo.Message); ok {
		return gogo.Marshal(gproto)
	}

	if b, ok := v.([]byte); ok {
		return b, nil
	}

	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Array, reflect.Slice, reflect.Struct, reflect.Map:
		return json.Marshal(v)
	}

	return nil, fmt.Errorf("unknown type")
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

// NewMarshaler 如果 marshal 为 nil, 则会使用 Default
func NewMarshaler(marshal Marshaler, specialMarshalers ...map[message.MsgID]Marshal) Marshaler {
	m := map[message.MsgID]Marshaler{}
	for _, specialMarshaler := range specialMarshalers {
		for k, v := range specialMarshaler {
			m[k] = v
		}
	}

	if marshal == nil {
		marshal = Marshal(DefaultMarshal)
	}

	return &marshaler{
		marshal:          marshal,
		specialMarshaler: m,
	}
}

func (m *marshaler) Marshal(t message.MsgID, v interface{}) ([]byte, error) {
	if bytes, ok := v.(ProvideMarshaler); ok {
		return bytes.Bytes()
	}

	// 判断
	marshal := m.specialMarshaler[t]
	if marshal == nil {
		marshal = m.marshal
	}

	return marshal.Marshal(t, v)
}
