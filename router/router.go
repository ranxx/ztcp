package router

import (
	"net"

	"github.com/ranxx/ztcp/message"
	"github.com/ranxx/ztcp/options"
)

// Router 路由
type Router struct {
	opt         *options.Options
	middlewares []Handler
	conn        net.Conn
	group       *Group
	MsgID       message.MsgID
	Headler     Handler
}

// NewRouter 新加路由
func NewRouter(msgid message.MsgID, handler Handler, opts ...options.Option) *Router {
	opt := options.DefaultOptions()

	for _, v := range opts {
		v(opt)
	}

	return &Router{
		opt:     opt,
		MsgID:   msgid,
		Headler: handler,
	}
}

// Use 添加中间件，适用于该路由
func (r *Router) Use(mid ...Handler) *Router {
	r.middlewares = append(r.middlewares, mid...)
	return r
}

// SetGroup 设置分组
func (r *Router) SetGroup(g *Group) *Router {
	r.group = g
	return r
}
