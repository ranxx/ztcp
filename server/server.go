package server

import (
	"net"

	"github.com/ranxx/ztcp/conn"
	"github.com/ranxx/ztcp/conner"
)

// Server ...
type Server struct {
	network  string
	address  string
	listener net.Listener
	opt      *Options
}

// NewServer new
func NewServer(network, address string, opts ...Option) *Server {
	opt := DefaultOptions()

	for _, v := range opts {
		v(opt)
	}

	if opt.genConner == nil {
		opt.genConner = func(i int64, c net.Conn) (conner.Conner, error) {
			return conn.NewConn(i, c, opt.genOptions...), nil
		}
	}

	return &Server{
		network:  network,
		address:  address,
		listener: opt.listener,
		opt:      opt,
	}
}

// Start 开始
func (s *Server) Start(success ...func(l net.Listener) error) error {
	if s.listener == nil {
		listener, err := net.Listen(s.network, s.address)
		if err != nil {
			return err
		}
		s.listener = listener
	}
	for _, v := range success {
		if err := v(s.listener); err != nil {
			return err
		}
	}
	s.doListener()
	return nil
}

// Del 删除
func (s *Server) Del(id int64) {
	s.opt.manager.Del(id)
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
		conner, err := s.opt.genConner(s.opt.indexMgr.NewIndex(), cn)
		if err != nil {
			continue
		}
		s.opt.manager.AddConn(conner)
		// 开启
		conner.Start()
	}
}

// GetManager 获取 manager
func (s *Server) GetManager() *conner.Manager {
	return s.opt.manager
}

// Close 关闭
func (s *Server) Close() {
	close(s.opt.close)
	s.opt.manager.Close()
}
