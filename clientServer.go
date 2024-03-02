package micronet

import "net/rpc"

/**
 * The basic ClientServer functions
 */
type I_ClientServer interface {
	I_Client
	I_Server
}

/**
 * The ClientServer structure is an rpc client and server that can recieve and send requests to and from other rpc servers and clients
 */
type ClientServer struct {
	I_ClientServer
	*Server
	*Client
}

/**
 * InitClientServer creates a ClientServer, inheriting from Client and Server
 * @param selfNetwork is the server's network config
 * @param remoteNetwork is the remote server's network config
 * @return the initialized ClientServer or error 
 */
func InitClientServer(selfNetwork NetConf, remoteNetwork NetConf) (*ClientServer, error) {
	cli := InitClient(remoteNetwork)
	srv, err := InitServer(selfNetwork)
	if err != nil {
		return nil, err
	}

	return &ClientServer{Client: cli, Server: srv}, nil
}

/**
 * Register any additional handler
 * @param rcvr any structure that implements at leaste one handler prototyped function
 * @return an potential registration error
 */
func (s *ClientServer) Register(rcvr any) error {
	return s.Server.Register(rcvr)
}

/**
 * Start the ClientServer that was initialized with a netork config
 * You might consider starting the server in a goroutine
 * @return potential networking errors
 */
func (s *ClientServer) Start() error {
	return s.Server.Start()
}

/**
 * Dial creates the client's connexion to the remote Server
 * @return a potential network error
 */
func (c *ClientServer) Dial() (error) {
	return c.Client.Dial()
}

/**
 * Call sends a synchronous request to the remote Server
 * Use Go() for async request
 * @param serviceMethod is the remote's "handler.function" to call
 * @param args is the derefenced request of any type
 * @param reply is the derefenced response of any type 
 * @return a potential network error
 */
func (c *ClientServer) Call(serviceMethod string, args any, reply any) (error) {
	return c.Client.Call(serviceMethod, args, reply)
}

/**
 * Go sends a asynchronous request request to to the remote Server
 * Use Call() for a synchronous request
 * @param serviceMethod is the remote's "handler.function" to call
 * @param args is the derefenced request of any type
 * @param reply is the derefenced response of any type 
 * @param done channel will signal when the call is complete by returning the same Call object. If done is nil, Go will allocate a new channel. If non-nil, done must be buffered or Go will deliberately crash
 * @return the done channel
 */
func (c *ClientServer) Go(serviceMethod string, args any, reply any, done chan *rpc.Call) (*rpc.Call) {
	return c.Client.Go(serviceMethod, args, reply, done)
}

/**
 * Stop the running server
 * Same as Close()
 */
func (c *ClientServer) Stop() (error) {
	c.Server.Stop()
	err := c.Client.Close()
	return err
}

/**
 * Close the running server
 * Same as Stop()
 */
func (c *ClientServer) Close() (error) {
	return c.Stop()
}

/**
 * Ping allows to test a client and server connection
 * It is registered by default by the ClientServer
 */
func (c *ClientServer) Ping() (error) {
	return c.Client.Ping()
}