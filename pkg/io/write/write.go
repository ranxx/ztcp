package write

import (
	"io"

	"github.com/ranxx/ztcp/pkg/message"
	"github.com/ranxx/ztcp/pkg/pack"
)

// Writer ...
type Writer interface {
	With(io.Writer) Writer

	// io.Writer
	Write(p []byte) (n int, err error)

	WriteBytes(message.MsgID, []byte) (int, error)

	WriteMessager(message.Messager) (int, error)

	// Marshal(v) -> Write
	WriteValue(v interface{}) (int, error)

	// Marshal(v) -> Write
	WriteValueWithID(id message.MsgID, v interface{}) (int, error)

	Packer() pack.Packer
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

func (w *writer) Write(p []byte) (n int, err error) {
	return w.Writer.Write(p)
}

func (w *writer) WriteBytes(id message.MsgID, data []byte) (int, error) {
	msg := w.opt.genMessage(id, data)

	return w.WriteMessager(msg)
}

func (w *writer) WriteMessager(msg message.Messager) (int, error) {
	if w.opt.stop {
		return 0, nil
	}

	data, err := w.opt.packer.Pack(msg)
	if err != nil {
		return 0, err
	}

	n, err := w.Write(data)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (w *writer) WriteValue(v interface{}) (int, error) {
	if v == nil {
		return 0, nil
	}

	key := typeUniqueString(v)

	id := w.opt.typeMsgID[key]

	data, err := w.opt.marshal.Marshal(id, v)
	if err != nil {
		return 0, err
	}

	return w.WriteBytes(id, data)
}

func (w *writer) WriteValueWithID(id message.MsgID, v interface{}) (int, error) {
	if v == nil {
		return 0, nil
	}
	data, err := w.opt.marshal.Marshal(id, v)
	if err != nil {
		return 0, err
	}

	return w.WriteBytes(id, data)
}

func (w *writer) WithStop(stop bool) {
	w.opt.stop = stop
}

func (w *writer) Packer() pack.Packer {
	return w.opt.packer
}
