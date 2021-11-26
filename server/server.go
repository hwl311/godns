package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"ruotian.vip/godns/message"
)

type Backend interface {
	Query(msg *message.Message, q *message.Question, mb *message.MsgBuilder) error //
	STATUS() int                                                                   //服务器状态
	RecursionAvailable() bool                                                      //服务器是否支持递归查询
}

type Server struct {
	Backend Backend
	ctx     context.Context
	cancel  context.CancelFunc
	conn    net.Conn
}

func (s *Server) Close() {
	s.cancel()
	s.conn.Close()
}

func (s *Server) ListenAndServe(ip string, port int) error {
	defer s.Close()

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	s.conn = conn
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := make([]byte, 4096)
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			len, caddr, err := conn.ReadFrom(buf)
			if err != nil {
				log.Println("ReadFrom error:", err.Error())
				continue
			}
			res := s.handle(buf[:len])
			conn.WriteTo(res, caddr)
		}
	}
}

func (s *Server) handle(request []byte) (response []byte) {
	req := message.Parse(request)
	mb := message.NewResMsgBuilder(req)
	if req.QR != 0 {
		//TODO: error
		return mb.ToError(message.CODE_FORMAT_ERROR)
	}
	for _, q := range req.Question {
		err := s.Backend.Query(req, q, mb)
		if err != nil {
			return mb.ToError(message.CODE_NOT_IMPLEMENTED)
		}
	}
	return mb.ToBytes()
}

func New(be Backend) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{be, ctx, cancel, nil}
}
