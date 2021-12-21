package server

import (
	"log"
	"net"

	"github.com/ranxx/ztcp/conn"
	"github.com/ranxx/ztcp/conner"
	"github.com/ranxx/ztcp/pkg/index"
)

// Options ...
type Options struct {
	name        string
	indexMgr    *index.Index64
	manager     *conner.Manager
	listener    net.Listener
	listenAfter func(net.Listener) // 监听之后
	close       chan struct{}
	genConner   func(int64, net.Conn) conner.Conner
	genOptions  []conn.Option
}

// DefaultOptions 默认
func DefaultOptions() *Options {
	opt := &Options{
		name:        "ztcp",
		indexMgr:    index.NewIndexI64(),
		manager:     conner.NewManager(),
		listenAfter: nil,
		listener:    nil,
		close:       make(chan struct{}),
		genOptions:  make([]conn.Option, 0, 10),
		genConner:   nil,
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

// WithIndexer 下标管理
func WithIndexer(i *index.Index64) Option {
	return func(o *Options) {
		o.indexMgr = i
	}
}

// WithConnManager conn的管理器
func WithConnManager(m *conner.Manager) Option {
	return func(o *Options) {
		o.manager = m
	}
}

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

// WithClose ...
func WithClose(close chan struct{}) Option {
	return func(o *Options) {
		o.close = close
	}
}

// WithGenConner gen conner
// 如果此选项被构造，会忽略 WithConnOptions 选项
func WithGenConner(gen func(int64, net.Conn) conner.Conner) Option {
	return func(o *Options) {
		o.genConner = gen
	}
}

// WithConnOptions gen conner 的options
// 如果已经构造了 genConner，将会忽略此选项，否则将会被使用
func WithConnOptions(opts ...conn.Option) Option {
	return func(o *Options) {
		o.genOptions = opts
	}
}
