package read

import (
	"io"

	"github.com/ranxx/ztcp/pkg/buffer"
	"github.com/ranxx/ztcp/pkg/message"
	"github.com/ranxx/ztcp/pkg/pack"
)

// Reader reader
type Reader interface {
	buffer.Buffer

	// io.Reader
	Read(p []byte) (n int, err error)

	With(io.Reader)

	// Read 如果 stop为true
	//
	// 则返回的 Messager 为 nil
	ReadMessage() (message.Messager, error)

	// ReadHeader 包括以下方法
	//
	//  GetDataLength() uint32
	//  GetMsgID() MsgID
	ReadHeader() (message.Messager, []byte, error)

	// ReadBody 包括以下方法
	//
	//  SetData([]byte)
	//  GetData() []byte
	ReadBody(message.Messager) (message.Messager, error)

	WithStop(stop bool)

	Packer() pack.Packer
}

type reader struct {
	buffer.Buffer
	// io.Reader

	opt *Options
}

// DefaultReader ...
func DefaultReader(r io.Reader, opts ...Option) Reader {
	opt := DefaultOptions()

	for _, v := range opts {
		v(opt)
	}

	read := &reader{
		Buffer: buffer.NewBuffer(r),
		opt:    opt,
	}

	return read
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.opt.stop {
		return 0, nil
	}
	return r.Buffer.Read(p)
}

func (r *reader) With(ir io.Reader) {
	r.Buffer.With(ir)
}

func (r *reader) ReadMessage() (message.Messager, error) {
	if r.opt.stop {
		return nil, nil
	}

	// 读 head
	headData := make([]byte, r.opt.packer.GetHeadLength())
	if _, err := io.ReadFull(r, headData); err != nil {
		return nil, err
	}

	// 解包 head
	msg, err := r.opt.packer.UnpackHead(headData)
	if err != nil {
		return nil, err
	}

	// 读取 body
	data := make([]byte, msg.GetDataLength())
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	// 设置 data
	msg.SetData(data)

	return msg, nil
}

func (r *reader) ReadHeader() (message.Messager, []byte, error) {
	if r.opt.stop {
		return nil, nil, nil
	}

	// 读 head
	headData := make([]byte, r.opt.packer.GetHeadLength())
	if _, err := io.ReadFull(r, headData); err != nil {
		return nil, nil, err
	}

	// 解包 head
	msg, err := r.opt.packer.UnpackHead(headData)
	if err != nil {
		return nil, nil, err
	}

	return msg, headData, nil
}

func (r *reader) ReadBody(m message.Messager) (message.Messager, error) {
	if r.opt.stop {
		return nil, nil
	}

	data := make([]byte, m.GetDataLength())
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	m.SetData(data)

	return m, nil
}

func (r *reader) WithStop(stop bool) {
	r.opt.stop = stop
}

func (r *reader) Packer() pack.Packer {
	return r.opt.packer
}
