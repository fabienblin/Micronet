package server

import (
	"context"
	"log"
	"net"
	"net/rpc"

	"micronet/common"
)

/**
 * The basic Server functions
 */
type I_Server interface {
	Register(any) error
	Start() error
	Stop()
}

/**
 * The Server structure is an rpc server with it's network config and context
 */
type Server struct {
	*rpc.Server
	I_Server
	common.NetConf
	ctx            context.Context
	cancelFunction context.CancelFunc
}

/**
 * NewServer creates an rpc server with it's context and registers the default ping handler
 * @param network is the server's configuration
 * @return the initialized Server or error
 */
func NewServer(network common.NetConf) (*Server, error) {
	srv := &Server{NetConf: network}
	srv.Server = rpc.NewServer()
	srv.ctx, srv.cancelFunction = context.WithCancel(context.Background())

	errRegister := srv.Register(new(common.PingHandler))
	if errRegister != nil {
		return nil, errRegister
	}

	return srv, nil
}

/**
 * Register any additional handler
 * @param rcvr any structure that implements at leaste one handler prototyped function
 * @return an potential registration error
 */
func (s *Server) Register(rcvr any) error {
	errRegister := s.Server.Register(rcvr)
	if errRegister != nil {
		return errRegister
	}

	return nil
}

/**
 * Start the Server that was initialized with a netork config
 * You might consider starting the server in a goroutine
 * @return potential networking errors
 */
func (s *Server) Start() (error) {
	listener, errListen := net.Listen(s.Protocol, ":"+s.Port)
	if errListen != nil {
		return errListen
	}
	defer listener.Close()

	log.Printf("Server is running %+v", s.NetConf)

	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			conn, errAccept := listener.Accept()
			if errAccept != nil {
				log.Printf("Error accepting connection: %s", errAccept)
				continue
			}
			go s.ServeConn(conn)
		}
	}
}

/**
 * Stop the running server
 */
func (s *Server) Stop() {
	log.Printf("Stoping server %+v", s.NetConf)
	s.cancelFunction()
}
