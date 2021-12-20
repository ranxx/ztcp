package read

import (
	"io"

	"github.com/ranxx/ztcp/message"
)

// Reader reader
type Reader interface {
	With(io.Reader) Reader

	Read() (message.Messager, error)
}

type reader struct {
	io.Reader

	opt *Options
}

// DefaultReader ...
func DefaultReader(r io.Reader, opts ...Option) Reader {
	opt := DefaultOptions()

	for _, v := range opts {
		v(opt)
	}

	return &reader{
		Reader: r,
		opt:    opt,
	}
}

func (r *reader) With(nio io.Reader) Reader {
	r.Reader = nio
	return r
}

func (r *reader) Read() (message.Messager, error) {
	// 读 head
	headData := make([]byte, r.opt.packer.GetHeadLength())
	if _, err := io.ReadFull(r.Reader, headData); err != nil {
		return nil, err
	}

	// 解包 head
	msg, err := r.opt.packer.UnpackHead(headData)
	if err != nil {
		return nil, err
	}

	// 读取 body
	data := make([]byte, msg.GetDataLength())
	if _, err := io.ReadFull(r.Reader, data); err != nil {
		return nil, err
	}

	// 设置 data
	msg.SetData(data)

	return msg, nil
}
