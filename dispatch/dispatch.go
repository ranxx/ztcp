package dispatch

import (
	"github.com/ranxx/ztcp/conner"
	"github.com/ranxx/ztcp/pkg/message"
	"github.com/ranxx/ztcp/router"
)

// Dispatcher ...
type Dispatcher interface {
	Dispatch(message.Messager, conner.Conner)
}

type dispatcher struct {
	root *router.Root
	opt  *Options
}

// DefaultDispatcher default dispatch
func DefaultDispatcher(r *router.Root, opts ...Option) Dispatcher {
	opt := DefaultOptions()

	for _, v := range opts {
		v(opt)
	}
	if r == nil {
		r = router.NewRoot()
	}
	return &dispatcher{
		root: r,
		opt:  opt,
	}
}

func (d *dispatcher) Dispatch(msg message.Messager, conn conner.Conner) {
	// 分发
	d.root.Dispatch(conn, msg)
}
