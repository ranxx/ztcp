package router

import (
	"context"

	"github.com/ranxx/ztcp/conner"
	"github.com/ranxx/ztcp/handle"
	"github.com/ranxx/ztcp/options"
	"github.com/ranxx/ztcp/pkg/message"
	"github.com/ranxx/ztcp/request"
)

// Root ...
type Root struct {
	opt          *options.Options          // 选项
	routers      map[message.MsgID]*Router // 真正的 router
	groups       map[*Group]struct{}       // 分组 middle
	defaultGroup *Group                    // 默认分组
	middlewares  []handle.Handler          // 全局 middle
	unknown      handle.Handler            // 404 处理
}

// NewRoot root
func NewRoot(opt ...*options.Options) *Root {
	r := Root{
		routers:     make(map[message.MsgID]*Router),
		groups:      make(map[*Group]struct{}),
		middlewares: make([]handle.Handler, 0, 10),
	}
	g := NewGroup()
	g.root = &r
	r.groups[g] = struct{}{}
	return &r
}

// NewGroup 新建 group
func (r *Root) NewGroup() *Group {
	group := Group{middlewares: make([]handle.Handler, 0, 10), root: r}
	r.groups[&group] = struct{}{}
	return &group
}

// Merge 合并
func (r *Root) Merge(r2 *Root) *Root {
	for k, v := range r2.routers {
		r.routers[k] = v
	}
	for k, v := range r2.groups {
		r.groups[k] = v
	}
	r.middlewares = append(r.middlewares, r2.middlewares...)
	return r
}

// AddGroup 添加 group
func (r *Root) AddGroup(groups ...*Group) *Root {
	for _, g := range groups {
		r.groups[g] = struct{}{}
		if g.root == r {
			continue
		}
		g.root = r.Merge(g.root)
	}
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
func (r *Root) Use(mid ...handle.Handler) *Root {
	r.middlewares = append(r.middlewares, mid...)
	return r
}

// NotFound 未知消息处理
func (r *Root) NotFound(h handle.Handler) *Root {
	r.unknown = h
	return r
}

// Dispatch 分发
func (r *Root) Dispatch(conn conner.Conner, msg message.Messager) {
	r.dispatch(conn, msg)
}

func (r *Root) dispatch(conn conner.Conner, msg message.Messager) {
	router := r.routers[msg.GetMsgID()]
	if router == nil {
		ctx := context.Background()
		req := request.NewRequest(msg, conn)
		r.unknown.Serve(ctx, req)
		return
	}

	// TODO: worker
	go func() {
		ctx := context.Background()
		req := request.NewRequest(msg, conn)

		handlers := make([]handle.Handler, 0, 30)
		handlers = append(handlers, r.middlewares...)

		// 如果group被删了，不应该触发 group的中间件
		if router.group != nil {
			handlers = append(handlers, router.group.middlewares...)
		}

		handlers = append(handlers, router.middlewares...)
		handlers = append(handlers, router.Headler)

		for _, handler := range handlers {
			if req.GetAbort() {
				break
			}
			handler.Serve(ctx, req)
		}
	}()
}
