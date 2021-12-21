package server

import (
	"fmt"
	"net"

	"github.com/ranxx/ztcp/conn"
	"github.com/ranxx/ztcp/conner"
)

// Server ...
type Server struct {
	network  string
	ip       string
	port     int
	listener net.Listener
	opt      *Options
}

// NewServer new
func NewServer(network, ip string, port int, opts ...Option) *Server {
	opt := DefaultOptions()

	for _, v := range opts {
		v(opt)
	}

	if opt.genConner == nil {
		opt.genConner = func(i int64, c net.Conn) conner.Conner {
			return conn.NewConn(i, c, opt.genOptions...)
		}
	}

	return &Server{
		network:  network,
		ip:       ip,
		port:     port,
		listener: opt.listener,
		opt:      opt,
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
	s.doListener()
	return nil
}

func (s *Server) doListener() {
	defer func() {
		s.opt.manager.Close()
	}()
	for {
		select {
		case <-s.opt.close:
			return
		default:
		}

		cn, err := s.listener.Accept()
		if err != nil {
			// 是否退出
			panic(err)
		}
		conner := s.opt.genConner(s.opt.indexMgr.NewIndex(), cn)
		s.opt.manager.AddConn(conner)
		// 开启
		conner.Start()
	}
}
