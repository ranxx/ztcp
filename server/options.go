package server

import (
	"encoding/json"
	"log"
	"net"

	"github.com/ranxx/ztcp/conn"
	"github.com/ranxx/ztcp/index"
)

type (
	// Encode 编码
	Encode func(interface{}) ([]byte, error)

	// Decode 解码
	Decode func(data []byte, v interface{}) error
)

// Options ...
type Options struct {
	name     string
	indexMgr *index.Index64
	manager  *conn.Manager

	encode Encode // 消息编码
	decode Decode // 消息解码
	// manager            conn.IManager             // conn 的管理
	newConnWithNetConn func(net.Conn) conn.IConn // 通过 net.Conn 生成 Conn
	listener           net.Listener
	listenAfter        func(net.Listener) // 监听之后
	close              chan struct{}
}

// DefaultOptions 默认
func DefaultOptions() *Options {
	opt := &Options{
		name:    "ztcp",
		encode:  json.Marshal,
		decode:  json.Unmarshal,
		manager: conn.NewManager(),
		// closeConnWrite: true,
		// closeConnRead:  true,
		listenAfter: nil,
		listener:    nil,
	}

	opt.newConnWithNetConn = func(c net.Conn) Conn {
		return conn.NewConn(c, opt.indexMgr.NewIndex())
	}

	opt.listenAfter = func(l net.Listener) {
		log.Printf("name:%s moduel:server listener success addr:%s\n", opt.name, l.Addr().String())
	}

	return opt
}

// Option ...
type Option func(*Options)

// Options2Option ...
func Options2Option(opt *Options) Option {
	return func(o *Options) {
		o = opt
	}
}

// WithName name
func WithName(name string) Option {
	return func(o *Options) {
		o.name = name
	}
}

// WithEncode 编码
func WithEncode(ec Encode) Option {
	return func(o *Options) {
		o.encode = ec
	}
}

// WithDecode 解码
func WithDecode(de Decode) Option {
	return func(o *Options) {
		o.decode = de
	}
}

// WithConnManager conn的管理器
func WithConnManager(m ConnManager) Option {
	return func(o *Options) {
		o.manager = m
	}
}

// WithNewConnWithNetConn 新建 Conn
func WithNewConnWithNetConn(fn func(c net.Conn) Conn) Option {
	return func(o *Options) {
		o.newConnWithNetConn = fn
	}
}

// // WithCloseConnRead 设置读
// func WithCloseConnRead(close bool) Option {
// 	return func(o *Options) {
// 		o.closeConnRead = close
// 	}
// }

// // WithCloseConnWrite 设置写
// func WithCloseConnWrite(close bool) Option {
// 	return func(o *Options) {
// 		o.closeConnWrite = close
// 	}
// }

// WithListenAfter ...
func WithListenAfter(fn func(net.Listener)) Option {
	return func(o *Options) {
		o.listenAfter = fn
	}
}

// WithListener ...
func WithListener(l net.Listener) Option {
	return func(o *Options) {
		o.listener = l
	}
}
