package handle

import (
	"context"

	"github.com/ranxx/ztcp/request"
)

// Handler ...
type Handler interface {
	Serve(context.Context, *request.Request)
}

// WrapHandler ...
type WrapHandler func(context.Context, *request.Request)

// Serve ...
func (w WrapHandler) Serve(ctx context.Context, req *request.Request) {
	w(ctx, req)
}
