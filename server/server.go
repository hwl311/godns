package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"ruotian.vip/godns/message"
)

type Backend interface {
	Query(msg *message.Message, q *message.Question, mb *message.MsgBuilder) error //
	STATUS() int                                                                   //服务器状态
	RecursionAvailable() int                                                       //服务器是否支持递归查询
}

type Server struct {
	Backend Backend
	ctx     context.Context
	cancel  context.CancelFunc
	conn    net.Conn
}

func (s *Server) Close() {
	s.cancel()
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *Server) ListenAndServe(ip string, port int) error {
	defer s.Close()

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer conn.Close()
	s.conn = conn

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
			log.Println("WriteTo:", caddr, res)
			resmsg := message.Parse(res)
			j, _ := json.MarshalIndent(resmsg, "", "    ")
			log.Println(string(j))
			len, err = conn.WriteTo(res, caddr)
			log.Println(len, err)
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
	mb.SetRecursion(s.Backend.RecursionAvailable())
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
