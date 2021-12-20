package write

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ranxx/ztcp/encoding"
	"github.com/ranxx/ztcp/message"
	"github.com/ranxx/ztcp/pack"
)

// Options ...
type Options struct {
	genMessage message.GenMessage
	packer     pack.Packer
	typeMsgID  map[reflect.Type]message.MsgID
	marshal    encoding.Marshaler
}

// Option ...
type Option func(*Options)

// DefaultOptions ...
func DefaultOptions() *Options {
	return &Options{
		genMessage: message.DefaultMessager,
		packer:     pack.DefaultPacker(message.DefaultMessager),
		marshal: encoding.NewMarshaler(encoding.Marshal(func(mi message.MsgID, i interface{}) ([]byte, error) {
			if b, ok := i.([]byte); ok {
				return b, nil
			}
			return json.Marshal(i)
		})),
	}
}

// WithGenMessage new 消息
func WithGenMessage(gen message.GenMessage) Option {
	return func(o *Options) {
		o.genMessage = gen
	}
}

// WithPacker 打包
func WithPacker(packer pack.Packer) Option {
	return func(o *Options) {
		o.packer = packer
	}
}

// WithTypeMsgID 消息类型的id
func WithTypeMsgID(tfunc ...func() (interface{}, message.MsgID)) Option {
	return func(o *Options) {
		o.typeMsgID = make(map[reflect.Type]message.MsgID)
		for _, f := range tfunc {
			t, v := f()
			tt := reflect.TypeOf(t)
			fmt.Println(tt.Kind())
			o.typeMsgID[tt] = v
		}
	}
}

// WithMarshal 打包
func WithMarshal(marshal encoding.Marshaler) Option {
	return func(o *Options) {
		o.marshal = marshal
	}
}
