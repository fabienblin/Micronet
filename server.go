package micronet

import (
	"context"
	"log"
	"net"
	"net/rpc"
)

type I_Server interface {
	Register(any) (error)
	Start() (error)
	Stop()
}

type Server struct {
	*rpc.Server
	I_Server
	NetConf
	ctx            context.Context
	cancelFunction context.CancelFunc
}

func InitServer(network NetConf) (*Server, error) {
	srv := &Server{NetConf: network}
	srv.Server = rpc.NewServer()
	srv.ctx, srv.cancelFunction = context.WithCancel(context.Background())

	// Register ping request
	err := srv.Server.Register(new(PingHandler))
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func (s *Server) Register(rcvr any) error {
	err := s.Server.Register(rcvr)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Start() error {
	// Create a listener on port 1234
	listener, err := net.Listen(s.Protocol, ":"+s.Port)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Server is running %+v", s.NetConf)

	// Accept and handle incoming connections
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %s", err)
				continue
			}
			go s.Server.ServeConn(conn)
		}
	}
}

func (s *Server) Stop() {
	log.Printf("Stoping sever %+v", s.NetConf)
	s.cancelFunction()
}
