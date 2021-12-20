/*Package conn 读写数据

读函数，写函数，消息分发，打包机制，解析函数

*/
package conn

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Conner conn
type Conner interface {
	net.Conn
	ID() int64
	Start()
}

// conn ...
type conn struct {
	net.Conn          // net
	id       int64    // 唯一标识
	opt      *Options // 可选项
	rlock    sync.Mutex
}

// NewConn ...
func NewConn(id int64, _conn net.Conn, opts ...Option) Conner {
	opt := DefaultOptions()

	for _, v := range opts {
		v(opt)
	}

	return &conn{
		id:    id,
		Conn:  _conn,
		opt:   opt,
		rlock: sync.Mutex{},
	}
}

func (c *conn) ID() int64 {
	return c.id
}

func (c *conn) Start() {
	go c.gReading()

	// go c.gWriting()
}

func (c *conn) Read(b []byte) (n int, err error) {
	c.rlock.Lock()
	defer c.rlock.Unlock()
	return c.Conn.Read(b)
}

func (c *conn) gReading() {
	// 是否开启
	fmt.Println("开启 reading")
	for {
		if c.opt.closeConnRead {
			time.Sleep(time.Second)
			continue
		}
		select {
		case <-c.opt.close:
		default:
			c.reading()
		}
	}
}

func (c *conn) reading() error {
	c.rlock.Lock()
	defer c.rlock.Unlock()
	// 读 head
	headData := make([]byte, c.opt.packer.GetHeadLength())
	if _, err := io.ReadFull(c.Conn, headData); err != nil {
		return err
	}

	// 解包 head
	msg, err := c.opt.packer.UnpackHead(headData)
	if err != nil {
		return err
	}

	// 读取 body
	data := make([]byte, msg.GetDataLength())
	if _, err := io.ReadFull(c.Conn, data); err != nil {
		return err
	}

	// 设置 data
	msg.SetData(data)

	go c.opt.dispatcher.Dispatch(msg, c)
	return nil
}

// func (c *conn) gWriting() {
// 	// 是否开启
// 	for {
// 		if c.opt.closeConnWrite {
// 			time.Sleep(time.Second)
// 			continue
// 		}
// 		select {
// 		case <-c.opt.close:
// 		case fn := <-c.opt.writeChan:
// 			c.writing(fn)
// 		default:
// 			time.Sleep(time.Millisecond * 100)
// 		}
// 	}
// }

// func (c *conn) writing(fn writeFunc) error {
// 	msgid, value := fn()

// 	data, err := c.opt.marshaler.Marshal(msgid, value)
// 	if err != nil {
// 		return err
// 	}

// 	// 打包
// 	data, err = c.opt.packer.Pack(data)
// 	if err != nil {
// 		return err
// 	}

// 	if _, err := c.conn.Write(data); err != nil {
// 		return err
// 	}
// 	return nil
// }
