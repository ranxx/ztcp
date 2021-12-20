package router

import (
	"context"
	"net"

	"github.com/ranxx/ztcp/message"
	"github.com/ranxx/ztcp/options"
)

// Handler ...
type Handler interface {
	ServeNET(context.Context, net.Conn, interface{})
}

// WrapHandler ...
type WrapHandler func(context.Context, net.Conn, interface{})

// ServeNET ...
func (w WrapHandler) ServeNET(ctx context.Context, c net.Conn, v interface{}) {
	w(ctx, c, v)
}

// Root ...
type Root struct {
	opt          *options.Options          // 选项
	routers      map[message.MsgID]*Router // 真正的 router
	groups       map[*Group]struct{}       // 分组 middle
	defaultGroup *Group                    // 默认分组
	middlewares  []Handler                 // 全局 middle
}

// NewRoot root
func NewRoot(opt ...*options.Options) *Root {
	r := Root{
		routers:     make(map[message.MsgID]*Router),
		groups:      make(map[*Group]struct{}),
		middlewares: make([]Handler, 0, 10),
	}
	g := NewGroup()
	g.root = &r
	r.groups[g] = struct{}{}
	return &r
}

// NewGroup 新建 group
func (r *Root) NewGroup() *Group {
	group := Group{root: r}
	r.groups[&group] = struct{}{}
	return &group
}

// AddGroup 添加 group
func (r *Root) AddGroup(g *Group) *Root {
	r.groups[g] = struct{}{}
	return r
}

// AddRouter 添加路由
func (r *Root) AddRouter(rts ...*Router) *Root {
	for _, rt := range rts {
		r.routers[rt.MsgID] = rt
		rt.group = r.defaultGroup
	}
	return r
}

// Use 添加中间件，适用于所有路由
func (r *Root) Use(mid Handler) *Root {
	r.middlewares = append(r.middlewares, mid)
	return r
}

// Dispatch 分发
func (r *Root) Dispatch(msgid message.MsgID, conn net.Conn, v interface{}) {
	r.dispatch(msgid, conn, v)
}

func (r *Root) dispatch(msgid message.MsgID, conn net.Conn, v interface{}) {
	router := r.routers[msgid]
	if router == nil {
		return
	}
	go func() {

		ctx := context.Background()

		handlers := make([]Handler, 0, 30)
		handlers = append(handlers, r.middlewares...)
		// 如果group被删了，不应该触发 group的中间件
		if router.group != nil {
			handlers = append(handlers, router.group.middlewares...)
		}
		handlers = append(handlers, router.middlewares...)
		handlers = append(handlers, router.Headler)

		// v := 0
		for _, handler := range handlers {
			handler.ServeNET(ctx, conn, v)
		}
	}()
}
