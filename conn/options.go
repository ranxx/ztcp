package conn

import (
	"github.com/ranxx/ztcp/dispatch"
	"github.com/ranxx/ztcp/message"
	"github.com/ranxx/ztcp/pack"
	"github.com/ranxx/ztcp/read"
	"github.com/ranxx/ztcp/router"
	"github.com/ranxx/ztcp/write"
)

// type writeFunc func() (pack.MsgType, interface{})

// Options ...
type Options struct {
	name          string
	close         chan struct{}
	packer        pack.Packer         // 消息打包解包
	dispatcher    dispatch.Dispatcher // 消息分发
	reader        read.Reader
	writer        write.Writer
	closeConnRead bool // 关闭写
}

// DefaultOptions 默认
func DefaultOptions() *Options {
	packer := pack.DefaultPacker(message.DefaultMessager)
	opt := &Options{
		name:          "conn",
		close:         make(chan struct{}),
		packer:        packer,
		dispatcher:    dispatch.DefaultDispatcher(router.NewRoot()),
		reader:        read.DefaultReader(nil, read.WithPacker(packer)),
		writer:        write.DefaultWriter(nil, write.WithPacker(packer)),
		closeConnRead: false,
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

// WithCloseChannel 关闭 channel
func WithCloseChannel(ch chan struct{}) Option {
	return func(o *Options) {
		o.close = ch
	}
}

// WithPacker 打包/解包 函数
func WithPacker(p pack.Packer) Option {
	return func(o *Options) {
		o.packer = p
	}
}

// WithDispatcher 分发消息
func WithDispatcher(d dispatch.Dispatcher) Option {
	return func(o *Options) {
		o.dispatcher = d
	}
}

// WithReader 读
func WithReader(reader read.Reader) Option {
	return func(o *Options) {
		o.reader = reader
	}
}

// WithWriter 写
func WithWriter(writer write.Writer) Option {
	return func(o *Options) {
		o.writer = writer
	}
}

// WithCloseConnRead 关闭读
func WithCloseConnRead(close bool) Option {
	return func(o *Options) {
		o.closeConnRead = close
	}
}
