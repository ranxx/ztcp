/*Package conn 读写数据

读函数，写函数，消息分发，打包机制，解析函数

*/
package conn

import (
	"net"
	"sync"
	"time"

	"github.com/ranxx/ztcp/conner"
	"github.com/ranxx/ztcp/pkg/io/read"
	"github.com/ranxx/ztcp/pkg/io/write"
)

// conn ...
type conn struct {
	net.Conn            // net
	id       int64      // 唯一标识
	rlock    sync.Mutex // 读锁
	wlock    sync.Mutex // 写锁
	opt      *Options   // 可选项
	closed   bool       // 关闭
}

// NewConn ...
func NewConn(id int64, _conn net.Conn, opts ...Option) conner.Conner {
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
		if c.opt.stop {
			time.Sleep(time.Millisecond * 100)
		}
		select {
		case <-c.opt.close:
		default:
			c.reading()
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func (c *conn) reading() error {
	msg, err := c.opt.reader.ReadMessage()
	if err != nil {
		return err
	}
	if msg == nil {
		return nil
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

func (c *conn) Close() error {
	c.closed = true
	close(c.opt.close)
	return c.Conn.Close()
}

func (c *conn) Closed() bool {
	return c.closed
}

func (c *conn) Extra() interface{} {
	return c.opt.extra
}

func (c *conn) Dispatcher() conner.Dispatcher {
	return c.opt.dispatcher
}

func (c *conn) WithStop(stop bool) {
	c.opt.stop = stop
}
