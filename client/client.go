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
func NewClient(network common.NetConf) *Client {
	cli := &Client{remote: network}
	cli.SetReconnectionConf(3, 1)

	return cli
}

/**
 * Dial creates the client's connexion to the remote Server
 * @return a potential network error
 */
func (c *Client) Dial() error {
	if c == nil {
		return fmt.Errorf("ni client")
	}

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
func (c *Client) Call(serviceMethod string, args any, reply any) error {
	if c == nil || c.Client == nil {
		return fmt.Errorf("ni client")
	}

	var deferedError error
	defer func() {
		// normal call can panic if client is nil or if connection was lost
		if r := recover(); r != nil {
			deferedError := c.reconnect()
			if deferedError != nil {
				return
			}

			// retry to request
			deferedError = c.Client.Call(serviceMethod, args, reply)
			if deferedError != nil {
				return
			}
		}
	}()

	err := c.Client.Call(serviceMethod, args, reply)
	if err != nil {
		return err
	}

	return deferedError
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
func (c *Client) Go(serviceMethod string, args any, reply any, done chan *rpc.Call) *rpc.Call {
	if c == nil {
		return &rpc.Call{
			ServiceMethod: serviceMethod,
			Args:          args,
			Reply:         reply,
			Error:         fmt.Errorf("nil client"),
			Done:          nil,
		}
	}

	var deferedError *rpc.Call
	defer func() {
		// normal call can panic if client is nil or if connection was lost
		if r := recover(); r != nil {
			deferedError := c.reconnect()
			if deferedError != nil {
				return
			}

			// retry to request
			deferedError = c.Client.Call(serviceMethod, args, reply)
			if deferedError != nil {
				return
			}
		}
	}()

	call := c.Client.Go(serviceMethod, args, reply, done)
	if call != nil {
		return call
	}

	return deferedError
}

/**
 * Close calls the underlying codec's Close method. If the connection is already shutting down, ErrShutdown is returned
 */
func (c *Client) Close() error {
	if c == nil || c.Client == nil {
		return fmt.Errorf("ni client")
	}

	var deferedError error
	defer func() {
		if r := recover(); r != nil {
			deferedError = fmt.Errorf("%v", r)
		}
	}()

	err := c.Client.Close()
	if err != nil {
		return err
	}

	return deferedError
}

/**
 * Ping allows to test a client and server connection
 * It is registered by default by the Server
 */
func (c *Client) Ping() error {
	if c == nil || c.Client == nil {
		return fmt.Errorf("ni client")
	}

	request := common.Ping{Data: "PING"}
	response := common.Pong{}
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

/**
 * Set reconnection logic
 * @param iterationLimit is the number of times the reconnection should try
 * @param timeInterval is the time between each try
 */
func (c *Client) SetReconnectionConf(iterationLimit int, timeInterval time.Duration) {
	if c == nil {
		return
	}

	c.iterationLimit = iterationLimit
	c.timeInterval = timeInterval
}

func (c *Client) reconnect() error {
	if c == nil || c.Client == nil {
		return fmt.Errorf("ni client")
	}

	if c.isReconnecting {
		return nil
	}
	c.isReconnecting = true
	defer func() {
		c.isReconnecting = false
	}()

	// var err error

	for i := 0; i < c.iterationLimit; i++ {
		log.Printf("reconnexion attempt %d/%d to %+v\n", i, c.iterationLimit, c.remote)
		err := c.Dial()
		if err != nil {
			log.Println(err)
		} else {
			log.Println("reconnexion succeeded")
			return nil
		}

		time.Sleep(time.Second * c.timeInterval)
	}

	return common.MicronetReconnectTimeoutError{NetConf: c.remote}
}
