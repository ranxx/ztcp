package write

import (
	"fmt"
	"reflect"

	"github.com/ranxx/ztcp/pkg/encoding"
	"github.com/ranxx/ztcp/pkg/message"
	"github.com/ranxx/ztcp/pkg/pack"
)

// Options ...
type Options struct {
	genMessage message.GenMessage
	packer     pack.Packer
	typeMsgID  map[string]message.MsgID
	marshal    encoding.Marshaler
	stop       bool
}

// Option ...
type Option func(*Options)

// DefaultOptions ...
func DefaultOptions() *Options {
	return &Options{
		genMessage: message.DefaultMessager,
		packer:     pack.DefaultPacker(message.DefaultMessager),
		typeMsgID:  make(map[string]message.MsgID),
		marshal:    encoding.NewMarshaler(nil),
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
		for _, f := range tfunc {
			t, v := f()
			key := typeUniqueString(t)
			o.typeMsgID[key] = v
		}
	}
}

// WithMarshal 打包
func WithMarshal(marshal encoding.Marshaler) Option {
	return func(o *Options) {
		o.marshal = marshal
	}
}

// WithStop ...
func WithStop(stop bool) Option {
	return func(o *Options) {
		o.stop = stop
	}
}

func typeUniqueString(i interface{}) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 如果是基础类型
	if len(t.Name()) <= 0 && len(t.PkgPath()) <= 0 {
		return t.Kind().String()
	}

	return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
}
