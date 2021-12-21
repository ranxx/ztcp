package request

import (
	// "github.com/go-kit/kit/util/conn"
	"github.com/ranxx/ztcp/conner"
	"github.com/ranxx/ztcp/pkg/message"
)

// Request request
type Request struct {
	M     message.Messager
	C     conner.Conner
	abort bool
}

// NewRequest new requests
func NewRequest(m message.Messager, c conner.Conner) *Request {
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
