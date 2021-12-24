package buffer

import (
	"io"
)

// Buffer buffer
type Buffer interface {
	WithReader(r io.Reader) Buffer

	Reset()

	Read(p []byte) (n int, err error)

	Write(data []byte) (n int, err error)
}

type buffer struct {
	io.Reader
	buf    []byte
	offset int
}

// NewBuffer new Buffer
func NewBuffer(r io.Reader) Buffer {
	return &buffer{
		Reader: r,
		buf:    make([]byte, 0, 1024),
	}
}

func (b *buffer) WithReader(r io.Reader) Buffer {
	b.Reader = r
	return b
}

func (b *buffer) Reset() {
	b.buf = b.buf[:0]
	b.offset = 0
}

func (b *buffer) Write(data []byte) (n int, err error) {
	// 判断 buf是否被读完, 如果是，则重新copy
	length := len(b.buf)

	// 如果 buf 已经被读完，并且 data 与 buf 一样长
	if b.offset >= length && len(data) == length {
		return copy(b.buf, data), nil
	}

	if b.offset >= length {
		b.Reset()
	}
	b.buf = append(b.buf, data...)
	return len(data), nil
}

func (b *buffer) Read(p []byte) (n int, err error) {
	if len(p) <= 0 {
		return 0, nil
	}

	length := len(b.buf)

	// buf 为空
	if b.offset > length {
		return b.Reader.Read(p)
	}

	// buf 不为空，并且比 p 大
	if len(p)+int(b.offset) <= length {
		copy(p, b.buf[b.offset:len(p)+int(b.offset)])
		b.offset += (len(p))
		return len(p), nil
	}

	// 已有长度
	has := length - b.offset
	// 先读取 buf 中的数
	n1 := copy(p[:has], b.buf[b.offset:length])
	if len(p) == n1 {
		b.offset = length
		return n1, nil
	}

	// 在从 Reader 中读取剩余的字段
	n2, err := b.Reader.Read(p[has:])
	if err != nil {
		return 0, err
	}
	b.offset = length
	return n1 + n2, nil
}
