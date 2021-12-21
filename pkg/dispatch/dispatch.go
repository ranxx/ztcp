package dispatch

import (
	"fmt"
	"net"

	"github.com/ranxx/ztcp/pkg/encoding"
	"github.com/ranxx/ztcp/pkg/message"
	"github.com/ranxx/ztcp/router"
)

// Dispatcher ...
type Dispatcher interface {
	Dispatch(message.Messager, net.Conn)
}

type dispatcher struct {
	root        *router.Root
	unmarshaler encoding.Unmarshaler // 消息反序列化
}

// DefaultDispatcher default dispatch
func DefaultDispatcher(r *router.Root, unmarshaler encoding.Unmarshaler) Dispatcher {
	if unmarshaler == nil {
		unmarshaler = encoding.Unmarshal(func(mi message.MsgID, b []byte) (interface{}, error) {
			return b, nil
		})
	}
	return &dispatcher{
		root:        r,
		unmarshaler: unmarshaler,
	}
}

func (d *dispatcher) Dispatch(msg message.Messager, conn net.Conn) {
	// 先 返序列化
	su := d.unmarshaler

	v, err := su.Unmarshal(msg.GetMsgID(), msg.GetData())
	if err != nil {
		fmt.Println(err)
		return
	}

	// 分发
	d.root.Dispatch(msg.GetMsgID(), conn, v)
}
