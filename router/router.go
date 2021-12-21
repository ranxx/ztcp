package router

import (
	"net"

	"github.com/ranxx/ztcp/handle"
	"github.com/ranxx/ztcp/options"
	"github.com/ranxx/ztcp/pkg/message"
)

// Router 路由
type Router struct {
	opt         *options.Options
	middlewares []handle.Handler
	conn        net.Conn
	group       *Group
	MsgID       message.MsgID
	Headler     handle.Handler
}

// NewRouter 新加路由
func NewRouter(msgid message.MsgID, handler handle.Handler, opts ...options.Option) *Router {
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
func (r *Router) Use(mid ...handle.Handler) *Router {
	r.middlewares = append(r.middlewares, mid...)
	return r
}

// SetGroup 设置分组
func (r *Router) SetGroup(g *Group) *Router {
	r.group = g
	return r
}

// Bind 绑定消息
func (r *Router) Bind() *Router {
	return r
}
