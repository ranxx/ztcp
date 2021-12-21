package dispatch

import (
	"github.com/ranxx/ztcp/pkg/encoding"
)

// Options ...
type Options struct {
	unmarshaler encoding.Unmarshaler // 消息反序列化
}

// Option ...
type Option func(*Options)

// DefaultOptions ...
func DefaultOptions() *Options {
	return &Options{
		unmarshaler: nil,
	}
}

// WithUnmarshaler unmarshaler
func WithUnmarshaler(unmarshaler encoding.Unmarshaler) Option {
	return func(o *Options) {
		o.unmarshaler = unmarshaler
	}
}
