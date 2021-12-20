/*Package conn 读写数据

读函数，写函数，消息分发，打包机制，解析函数

*/
package conn

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ranxx/ztcp/read"
	"github.com/ranxx/ztcp/write"
)

// Conner conn
type Conner interface {
	net.Conn
	ID() int64
	Start()
	Writer() write.Writer
	Reader() read.Reader
}

// conn ...
type conn struct {
	net.Conn       // net
	id       int64 // 唯一标识
	rlock    sync.Mutex
	wlock    sync.Mutex
	opt      *Options // 可选项
}

// NewConn ...
func NewConn(id int64, _conn net.Conn, opts ...Option) Conner {
	conn := &conn{
		id:    id,
		Conn:  _conn,
		rlock: sync.Mutex{},
		wlock: sync.Mutex{},
	}

	opt := DefaultOptions()

	for _, v := range opts {
		v(opt)
	}

	opt.writer.With(conn)
	opt.reader.With(conn)

	conn.opt = opt
	return conn
}

func (c *conn) ID() int64 {
	return c.id
}

func (c *conn) Start() {
	go c.gReading()
}

func (c *conn) Read(b []byte) (n int, err error) {
	c.rlock.Lock()
	defer c.rlock.Unlock()
	return c.Conn.Read(b)
}

func (c *conn) Write(b []byte) (n int, err error) {
	c.wlock.Lock()
	defer c.wlock.Unlock()
	return c.Conn.Write(b)
}

func (c *conn) gReading() {
	// 是否开启
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
	msg, err := c.opt.reader.Read()
	if err != nil {
		fmt.Println(err)
		return err
	}
	go c.opt.dispatcher.Dispatch(msg, c)
	return nil
}

func (c *conn) Writer() write.Writer {
	return c.opt.writer
}

func (c *conn) Reader() read.Reader {
	return c.opt.reader
}
