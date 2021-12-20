package handler

import (
	"context"
	"net"
)

// Handler 处理
type Handler interface {
	Handle(context.Context, net.Conn, interface{})
}

type handler struct {
	typeHandler   map[interface{}]func(interface{})
	typeOfHandler map[interface{}]string
}

func (h *handler) AddHandler() {

}
