package request

import (
	"net"

	"github.com/ranxx/ztcp/pkg/message"
)

// Request request
type Request struct {
	M     message.Messager
	C     net.Conn
	abort bool
}

// NewRequest new requests
func NewRequest(m message.Messager, c net.Conn) *Request {
	return &Request{
		M: m,
		C: c,
	}
}

// Abort ...
func (r *Request) Abort() {
	r.abort = true
}

// GetAbort ...
func (r *Request) GetAbort() bool {
	return r.abort
}
