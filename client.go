package micronet

import (
	"fmt"
	"log"
	"net/rpc"
)
/**
 * The basic Client functions
 */
type I_Client interface {
	Dial() (error)
	Call(string, any, any) (error)
	Go(string, any, any, chan *rpc.Call) (*rpc.Call)
	Close() (error)
	Ping() (error)
}

/**
 * The Client structure is an rpc client with the remote server's config
 */
type Client struct {
	*rpc.Client
	I_Client
	remote NetConf
}

/**
 * InitClient creates an rpc client and saves the remote server's network config
 * @param network is the remote server to call
 * @return the initialized Client
 */
func InitClient(network NetConf) (*Client) {
	cli := &Client{remote: network}

	return cli
}

/**
 * Dial creates the client's connexion to the remote Server
 * @return a potential network error
 */
func (c *Client) Dial() (error) {
	var err error
	c.Client, err = rpc.Dial(c.remote.Protocol, c.remote.Ip+":"+c.remote.Port)
	if err != nil {
		return err
	}
	
	return nil
}

/**
 * Call sends a synchronous request to the remote Server
 * Use Go() for async request
 * @param serviceMethod is the remote's "handler.function" to call
 * @param args is the derefenced request of any type
 * @param reply is the derefenced response of any type 
 * @return a potential network error
 */
func (c *Client) Call(serviceMethod string, args any, reply any) (error) {
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
func (c *Client) Go(serviceMethod string, args any, reply any, done chan *rpc.Call) (*rpc.Call) {
	return c.Client.Go(serviceMethod, args, reply, done)
}

/**
 * Close calls the underlying codec's Close method. If the connection is already shutting down, ErrShutdown is returned
 */
func (c *Client) Close() (error) {
	return c.Client.Close()
}

/**
 * Ping allows to test a client and server connection
 * It is registered by default by the Server
 */
func (c *Client) Ping() (error) {
	request := Ping{Data: "PING"}
	response := Pong{}
	err := c.Call("PingHandler.Ping", &request, &response)
	if err != nil {
		return err
	}
	
	if response.Data != "PONG" {
		return fmt.Errorf("response is not PONG")
	}
	
	log.Printf("PING %+v responded with %s", c.remote, response.Data)
	
	return nil
}