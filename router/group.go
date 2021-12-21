package router

import (
	"github.com/ranxx/ztcp/handle"
	"github.com/ranxx/ztcp/pkg/message"
)

// Group 分组
type Group struct {
	middlewares []handle.Handler // 只是删除middle

	root *Root
}

// NewGroup group
func NewGroup() *Group {
	return &Group{middlewares: make([]handle.Handler, 0, 10), root: &Root{
		routers:     make(map[message.MsgID]*Router),
		groups:      make(map[*Group]struct{}),
		middlewares: make([]handle.Handler, 0, 10),
	}}
}

// AddRouter 添加路由
func (r *Group) AddRouter(rts ...*Router) *Group {
	for _, rt := range rts {
		r.root.routers[rt.MsgID] = rt
		rt.group = r
	}
	return r
}

// Use 添加中间件，适用于当前分组
func (r *Group) Use(mid ...handle.Handler) *Group {
	r.middlewares = append(r.middlewares, mid...)
	return r
}
