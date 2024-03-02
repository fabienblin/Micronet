package micronet

import "net/rpc"

type I_ClientServer interface {
	I_Client
	I_Server
}

type ClientServer struct {
	I_ClientServer
	*Server
	*Client
}

func InitClientServer(selfNetwork NetConf, remoteNetwork NetConf) (*ClientServer, error) {
	cli := InitClient(remoteNetwork)
	srv, err := InitServer(selfNetwork)
	if err != nil {
		return nil, err
	}

	return &ClientServer{Client: cli, Server: srv}, nil
}

func (s *ClientServer) Register(rcvr any) error {
	return s.Server.Register(rcvr)
}

func (s *ClientServer) Start() error {
	return s.Server.Start()
}

func (c *ClientServer) Dial() (error) {
	return c.Client.Dial()
}

func (c *ClientServer) Call(serviceMethod string, args any, reply any) (error) {
	return c.Client.Call(serviceMethod, args, reply)
}

func (c *ClientServer) Go(serviceMethod string, args any, reply any, done chan *rpc.Call) (*rpc.Call) {
	return c.Client.Go(serviceMethod, args, reply, done)
}

func (c *ClientServer) Stop() (error) {
	c.Server.Stop()
	err := c.Client.Close()
	return err
}

func (c *ClientServer) Close() (error) {
	return c.Stop()
}

func (c *ClientServer) Ping() (error) {
	return c.Client.Ping()
}