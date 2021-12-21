package write

import (
	"io"

	"github.com/ranxx/ztcp/pkg/message"
)

// Writer ...
type Writer interface {
	With(io.Writer) Writer

	Write(message.MsgID, []byte) (int, error)

	WriteMessager(message.Messager) (int, error)

	WriteValue(interface{}) (int, error)
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

func (w *writer) Write(id message.MsgID, data []byte) (int, error) {
	msg := w.opt.genMessage(id, data)

	return w.WriteMessager(msg)
}

func (w *writer) WriteMessager(msg message.Messager) (int, error) {
	data, err := w.opt.packer.Pack(msg)
	if err != nil {
		return 0, err
	}

	n, err := w.Writer.Write(data)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (w *writer) WriteValue(v interface{}) (int, error) {
	key := typeUniqueString(v)

	id := w.opt.typeMsgID[key]

	data, err := w.opt.marshal.Marshal(id, v)
	if err != nil {
		return 0, err
	}

	return w.Write(id, data)
}
