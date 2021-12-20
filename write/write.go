package write

import (
	"io"
	"reflect"

	"github.com/ranxx/ztcp/message"
)

// Writer ...
type Writer interface {
	With(io.Writer) Writer

	Write(message.MsgID, []byte) error

	WriteMessager(message.Messager) error

	WriteValue(interface{}) error
}

// 最终 打包之后发给conn
type writer struct {
	io.Writer

	opt *Options
}

// DefaultWriter ...
func DefaultWriter(w io.Writer, opts ...Option) Writer {
	opt := DefaultOptions()

	for _, v := range opts {
		v(opt)
	}

	return &writer{
		Writer: w,
		opt:    opt,
	}
}

func (w *writer) With(iw io.Writer) Writer {
	w.Writer = iw
	return w
}

func (w *writer) Write(id message.MsgID, data []byte) error {
	msg := w.opt.genMessage(id, data)

	return w.WriteMessager(msg)
}

func (w *writer) WriteMessager(msg message.Messager) error {
	data, err := w.opt.packer.Pack(msg)
	if err != nil {
		return err
	}

	if _, err := w.Writer.Write(data); err != nil {
		return err
	}

	return nil
}

func (w *writer) WriteValue(v interface{}) error {
	id := w.opt.typeMsgID[reflect.TypeOf(v)]

	data, err := w.opt.marshal.Marshal(id, v)
	if err != nil {
		return err
	}

	return w.Write(id, data)
}
