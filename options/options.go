package options

import "encoding/json"

type (
	// Encode 编码
	Encode func(interface{}) ([]byte, error)

	// Decode 解码
	Decode func(data []byte, v interface{}) error
)

// Options ...
type Options struct {
	Encode Encode // 消息编码
	Decode Decode // 消息解码
}

// DefaultOptions 默认
func DefaultOptions() *Options {
	return &Options{
		Encode: json.Marshal,
		Decode: json.Unmarshal,
	}
}

// Option ...
type Option func(*Options)

// WithEncode 编码
func WithEncode(ec Encode) Option {
	return func(o *Options) {
		o.Encode = ec
	}
}

// WithDecode 解码
func WithDecode(de Decode) Option {
	return func(o *Options) {
		o.Decode = de
	}
}
