package znet

import (
	"io"
	"net"
	"time"

	"github.com/ranxx/ztcp/pack"
)

// WriteFunc ...
type WriteFunc func() (pack.MsgType, []byte)

type connection struct {
	conn        net.Conn
	close       chan struct{}
	writeFnChan chan WriteFunc
}

func newConnection(conn net.Conn) *connection {
	return &connection{conn: conn, close: make(chan struct{}), writeFnChan: make(chan WriteFunc, 100)}
}

func (c *connection) Start() {
	// 开启读
	// 开启写
}

// 读写，开启重连
func (c *connection) goRead() {
	for {
		select {
		// 关闭
		case <-c.close:
			return
		default:
			// 解包
			c.realRead()
		}
	}
}

func (c *connection) realRead() {
	// 首先读 header
	headBytes := make([]byte, 0, pack.HeadLength())
	if _, err := io.ReadFull(c.conn, headBytes); err != nil {
		panic(err)
		return
	}

	// 解析头
	msg := pack.Empty()
	if err := msg.UnpackHeadBytes(headBytes); err != nil {
		panic(err)
		return
	}

	// 读取body
	msg.Msg = make([]byte, 0, msg.Length)
	if _, err := io.ReadFull(c.conn, msg.Msg); err != nil {
		panic(err)
		return
	}

	// dispatch
	go c.dispatch(msg)
}

func (c *connection) dispatch(msg *pack.Package) {
	// TODO: 做什么
}

func (c *connection) goWrite() {
	for {
		select {
		case <-c.close:
		case fn := <-c.writeFnChan:
			c.realWrite(fn)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (c *connection) realWrite(fn WriteFunc) {
	msgid, body := fn()
	msg := pack.NewPackage(msgid, body)
	bytes, err := msg.PackBytes()
	if err != nil {
		panic(err)
	}
	_, err = c.conn.Write(bytes)
	if err != nil {
		panic(err)
	}
}
