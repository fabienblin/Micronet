package client

import (
	"fmt"
	"log"
	"net/rpc"
	"time"

	"micronet/common"
)

/**
 * The basic Client functions
 */
type I_Client interface {
	Dial() error
	Call(string, any, any) error
	Go(string, any, any, chan *rpc.Call) *rpc.Call
	Close() error
	Ping() error
}

/**
 * The Client structure is an rpc client with the remote server's config
 */
type Client struct {
	*rpc.Client
	I_Client
	remote         common.NetConf
	isReconnecting bool
	iterationLimit int
	timeInterval   time.Duration
}

/**
 * NewClient creates an rpc client and saves the remote server's network config
 * @param network is the remote server to call
 * @return the initialized Client
 */
func NewClient(network common.NetConf) (*Client, error) {
	cli := &Client{
		remote: network,
	}

	errDial := cli.Dial()
	if errDial != nil {
		return nil, errDial
	}

	cli.SetReconnectionConf(3, 1)

	return cli, nil
}

/**
 * Dial creates the client's connexion to the remote Server
 * You should use NewClient instead, it will Dial for you.
 * @return a potential network error
 */
func (c *Client) Dial() error {
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
 * @param request is the derefenced request of any type
 * @param response is the derefenced response of any type
 * @return a potential network error
 */
func (c *Client) Call(serviceMethod string, request any, response any) error {
	if c.Client == nil {
		return fmt.Errorf("nil client")
	}

	errCall := c.Client.Call(serviceMethod, request, response)
	if errCall == nil {
		return nil
	}

	if errReconnect := c.reconnect(); errReconnect != nil {
		return errReconnect
	}

	return c.Client.Call(serviceMethod, request, response)
}

/**
 * Go sends a asynchronous request request to to the remote Server
 * Use Call() for a synchronous request
 * @param serviceMethod is the remote's "handler.function" to call
 * @param request is the derefenced request of any type
 * @param response is the derefenced response of any type
 * @param done channel will signal when the call is complete by returning the same Call object.
 * 		If done is nil, Go will allocate a new channel.
 * 		If non-nil, done must be buffered or Go will deliberately crash
 * @return the done channel
 */
func (c *Client) Go(serviceMethod string, request any, response any, done chan *rpc.Call) *rpc.Call {
	if c.Client == nil {
		call := &rpc.Call{
			ServiceMethod: serviceMethod,
			Args:          request,
			Reply:         response,
			Error:         fmt.Errorf("nil client"),
			Done:          make(chan *rpc.Call, 1),
		}
		call.Done <- call
		return call
	}

	// first async attempt
	origCall := c.Client.Go(serviceMethod, request, response, done)

	// The call returned to the user
	wrappedCall := &rpc.Call{
		ServiceMethod: serviceMethod,
		Args:          request,
		Reply:         response,
		Done:          make(chan *rpc.Call, 1),
	}

	go func() {
		result := <-origCall.Done

		// If first attempt succeeded
		if result.Error == nil {
			wrappedCall.Error = nil
			wrappedCall.Done <- wrappedCall
			return
		}

		// Try reconnect
		if err := c.reconnect(); err != nil {
			wrappedCall.Error = result.Error // propagate error
			wrappedCall.Done <- wrappedCall
			return
		}

		// Retry
		retryCall := c.Client.Go(serviceMethod, request, response, nil)
		retryResult := <-retryCall.Done

		wrappedCall.Error = retryResult.Error // propagate retry error
		wrappedCall.Done <- wrappedCall
	}()

	return wrappedCall
}

/**
 * Close calls the underlying codec's Close method. If the connection is already shutting down, ErrShutdown is returned
 */
func (c *Client) Close() error {
	if c.Client == nil {
		return fmt.Errorf("nil client")
	}

	err := c.Client.Close()
	if err != nil {
		return err
	}

	return nil
}

/**
 * Ping allows to test a client and server connection
 * It is registered by default by the Server
 */
func (c *Client) Ping() error {
	if c.Client == nil {
		return fmt.Errorf("nil client")
	}

	request := common.Ping{Data: common.PING}
	response := common.Pong{}
	err := c.Call("PingHandler.Ping", &request, &response)
	if err != nil {
		return err
	}

	if response.Data != common.PONG {
		return fmt.Errorf("response is not PONG")
	}

	log.Printf("PING %+v responded with %s", c.remote, response.Data)

	return nil
}

/**
 * Set reconnection logic
 * @param iterationLimit is the number of times the reconnection should try
 * @param timeInterval is the time (in seconds) between each try
 */
func (c *Client) SetReconnectionConf(iterationLimit int, timeInterval time.Duration) {
	if c == nil {
		return
	}

	c.iterationLimit = iterationLimit
	c.timeInterval = timeInterval
}

func (c *Client) reconnect() error {
	if c.Client == nil {
		return fmt.Errorf("nil client")
	}

	if c.isReconnecting {
		return nil
	}
	c.isReconnecting = true

	for i := 0; i < c.iterationLimit; i++ {
		log.Printf("reconnexion attempt %d/%d to %+v\n", i+1, c.iterationLimit, c.remote)

		if err := c.Dial(); err != nil {
			log.Println("reconnect failed:", err)
		} else {
			log.Println("reconnexion succeeded")
			c.isReconnecting = false
			return nil
		}

		time.Sleep(time.Second * c.timeInterval)
	}

	c.isReconnecting = false
	return common.MicronetReconnectTimeoutError{NetConf: c.remote}
}
