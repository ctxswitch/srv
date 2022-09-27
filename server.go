package srv

import (
	"bufio"
	"context"
	"net"
	"time"
)

type Server struct {
	MaxConnections int32
	BaseContext    context.Context
	ReadTimeout    time.Duration
	handler        Handler
	pool           *Pool
}

func (s *Server) serve(network string, addr string) error {
	l, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	defer l.Close()

	ctx := s.BaseContext
	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return nil
		default:
			s.pool.Add()
			go s.handleConnection(c)
		}
	}
}

func (s *Server) handleConnection(c net.Conn) {
	ctx, cancel := context.WithTimeout(s.BaseContext, s.ReadTimeout)
	defer func() {
		cancel()
		c.Close()
		s.pool.Remove()
	}()

	r := &Request{
		Reader:  *bufio.NewReader(c),
		Context: ctx,
	}

	w := Writer{
		Connection: c,
	}

	s.handler.ServeTCP(w, r)
}

func (s *Server) defaulted() {
	if s.BaseContext == nil {
		s.BaseContext = context.Background()
	}

	if s.ReadTimeout == 0 {
		s.ReadTimeout = 10 * time.Minute
	}

	if s.MaxConnections == 0 {
		s.MaxConnections = 1
	}

	s.pool = NewPool(s.MaxConnections)
}

func ListenAndServe(network string, addr string, handler Handler) error {
	s := Server{}
	s.defaulted()
	s.handler = handler
	return s.serve(network, addr)
}
