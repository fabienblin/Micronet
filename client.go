package micronet

import (
	"fmt"
	"log"
	"net/rpc"
)

type I_Client interface {
	Dial() (error)
	Call(string, any, any) (error)
	Go(string, any, any, chan *rpc.Call) (*rpc.Call)
	Close() (error)
	Ping() (error)
}

type Client struct {
	*rpc.Client
	I_Client
	remote NetConf
}

func InitClient(network NetConf) (*Client) {
	cli := &Client{remote: network}

	return cli
}

func (c *Client) Dial() (error) {
	var err error
	c.Client, err = rpc.Dial(c.remote.Protocol, c.remote.Ip+":"+c.remote.Port)
	if err != nil {
		return err
	}
	
	return nil
}

func (c *Client) Call(serviceMethod string, args any, reply any) (error) {
	return c.Client.Call(serviceMethod, args, reply)
}

func (c *Client) Go(serviceMethod string, args any, reply any, done chan *rpc.Call) (*rpc.Call) {
	return c.Client.Go(serviceMethod, args, reply, done)
}

func (c *Client) Close() (error) {
	return c.Client.Close()
}

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