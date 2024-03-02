package micronet

import (
	"context"
	"log"
	"net"
	"net/rpc"
)

/**
 * The basic Server functions
 */
type I_Server interface {
	Register(any) (error)
	Start() (error)
	Stop()
}

/**
 * The Server structure is an rpc server with it's network config and context
 */
type Server struct {
	*rpc.Server
	I_Server
	NetConf
	ctx            context.Context
	cancelFunction context.CancelFunc
}

/**
 * InitServer creates an rpc server with it's context and registers the default ping handler
 * @param network is the server's configuration
 * @return the initialized Server or error 
 */
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

/**
 * Register any additional handler
 * @param rcvr any structure that implements at leaste one handler prototyped function
 * @return an potential registration error
 */
func (s *Server) Register(rcvr any) error {
	err := s.Server.Register(rcvr)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Start the Server that was initialized with a netork config
 * You might consider starting the server in a goroutine
 * @return potential networking errors
 */
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

/**
 * Stop the running server
 */
func (s *Server) Stop() {
	log.Printf("Stoping sever %+v", s.NetConf)
	s.cancelFunction()
}
