package dispatch

import (
	"net"

	"github.com/ranxx/ztcp/encoding"
	"github.com/ranxx/ztcp/message"
	"github.com/ranxx/ztcp/router"
)

// Dispatcher ...
type Dispatcher interface {
	Dispatch(message.Messager, net.Conn)
}

type dispatcher struct {
	root               *router.Root
	unmarshaler        encoding.Unmarshaler                   // 消息反序列化
	specialUnmarshaler map[message.MsgID]encoding.Unmarshaler // 特效消息体 反序列化
}

// DefaultDispatcher default dispatch
func DefaultDispatcher(r *router.Root) Dispatcher {
	return &dispatcher{
		root:               r,
		specialUnmarshaler: make(map[message.MsgID]encoding.Unmarshaler),
		unmarshaler: encoding.Unmarshal(func(mi message.MsgID, b []byte) (interface{}, error) {
			return b, nil
		}),
	}
}

func (d *dispatcher) Dispatch(msg message.Messager, conn net.Conn) {
	// 先 返序列化
	su := d.specialUnmarshaler[msg.GetMsgID()]
	if su == nil {
		su = d.unmarshaler
	}

	v, err := su.Unmarshal(msg.GetMsgID(), msg.GetData())
	if err != nil {
		return
	}

	// 分发
	d.root.Dispatch(msg.GetMsgID(), conn, v)
}
