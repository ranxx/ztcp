package conn

import (
	"github.com/ranxx/ztcp/dispatch"
	"github.com/ranxx/ztcp/message"
	"github.com/ranxx/ztcp/pack"
	"github.com/ranxx/ztcp/router"
)

// type writeFunc func() (pack.MsgType, interface{})

// Options ...
type Options struct {
	name       string
	close      chan struct{}
	packer     pack.Packer         // 消息打包解包
	dispatcher dispatch.Dispatcher // 消息分发,.;'
	// marshaler        encoding.Marshaler                  // 消息序列化
	// specialMarshaler map[pack.MsgType]encoding.Marshaler // 特效消息体 序列化
	// writeChan        chan writeFunc                      // 写消息 chanel
	// closeConnWrite   bool                                // 关闭读
	closeConnRead bool // 关闭写
}

// DefaultOptions 默认
func DefaultOptions() *Options {
	opt := &Options{
		name:          "conn",
		close:         make(chan struct{}),
		packer:        pack.DefaultPacker(message.Empty()),
		dispatcher:    dispatch.DefaultDispatcher(&router.Root{}),
		closeConnRead: false,
		// marshaler:     encoding.NewMarshaler(json.Marshal),
		// writeChan:     make(chan writeFunc, 100),
		// writeCh:  make(<-chan WriteFunc, 10),
		// dispatch: func(msgid pack.MsgType, body []byte) {},
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

// // WithMarshaler 序列化函数
// func WithMarshaler(m encoding.Marshaler) Option {
// 	return func(o *Options) {
// 		o.marshaler = m
// 	}
// }

// // WithUnmarshaler 反序列化函数
// func WithUnmarshaler(u encoding.Unmarshaler) Option {
// 	return func(o *Options) {
// 		o.unmarshaler = u
// 	}
// }

// // WithSpecialMarshaler 特殊的反序列函数
// func WithSpecialMarshaler(m map[pack.MsgType]encoding.Marshaler) Option {
// 	return func(o *Options) {
// 		o.specialMarshaler = m
// 	}
// }

// // AppendSpecialMarshaler 特殊的反序列函数
// func AppendSpecialMarshaler(sm map[pack.MsgType]encoding.Marshaler) Option {
// 	return func(o *Options) {
// 		for k, v := range sm {
// 			o.specialMarshaler[k] = v
// 		}
// 	}
// }

// // WithSpecialUnmarshaler 特殊的反序列函数
// func WithSpecialUnmarshaler(su map[pack.MsgType]encoding.Unmarshaler) Option {
// 	return func(o *Options) {
// 		o.specialUnmarshaler = su
// 	}
// }

// // AppendSpecialUnmarshaler 特殊的反序列函数
// func AppendSpecialUnmarshaler(su map[pack.MsgType]encoding.Unmarshaler) Option {
// 	return func(o *Options) {
// 		for k, v := range su {
// 			o.specialUnmarshaler[k] = v
// 		}
// 	}
// }

// // WithNewHeader header
// func WithNewHeader(f func(pack.MsgType, []byte) pack.Header) Option {
// 	return func(o *Options) {
// 		o.newHeader = f
// 	}
// }
