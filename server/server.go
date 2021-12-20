package server

import (
	"fmt"
	"net"

	"github.com/ranxx/ztcp/index"
)

// Server ...
type Server struct {
	opt     *Options
	network string
	ip      string
	port    int

	listener net.Listener
}

// NewServer new
func NewServer(network, ip string, port int, opts ...Option) *Server {
	opt := DefaultOptions()

	for _, v := range opts {
		v(opt)
	}

	return &Server{
		indexMgr: index.NewIndexI64(),
		opt:      opt,
		network:  network,
		ip:       ip,
		port:     port,
		manager:  &Manager{},
		close:    make(chan struct{}),
		listener: opt.listener,
	}
}

// Start 开始
func (s *Server) Start() error {
	if s.listener == nil {
		listener, err := net.Listen(s.network, fmt.Sprintf("%s:%d", s.ip, s.port))
		if err != nil {
			return err
		}
		s.listener = listener

		s.opt.listenAfter(s.listener)
	}

	go s.doListener()
	return nil
}

func (s *Server) doListener() {
	defer func() {
		s.manager.Close()
	}()
	for {
		select {
		case <-s.close:
			return
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			// 是否退出
			panic(err)
		}

		c := NewConn(conn, s.indexMgr.NewIndex(), Options2Option(s.opt))

		s.opt.newConnMiddle(c)

		// 具体什么时候开启读写
		s.manager.AddConn(c)
	}
}
