package common

import (
	"fmt"
	"log"

	"micronet/client"
	"micronet/clientServer"
	"micronet/common"
	"micronet/server"
)

/**
 * PUBLISHER :
 */

/**
 * The basic Publisher functions
 */
type I_Publisher interface {
	Publish(string)
}

/**
 * The Publisher is a Server that can subscribe, unsubscribe and publish to Subscribers
 */
type Publisher struct {
	I_Publisher
	*server.Server
	*PublisherHandler
}

/**
 * The basic PublisherHandler functions
 */
type I_PublisherHandler interface {
	Subscribe(*common.SubscribeRequest, *common.SubscribeResponse) error
	Unsubscribe(*common.SubscribeRequest, *common.SubscribeResponse) error
}

/**
 * The PublisherHandler can subscribe or unsubscribe Subscribers
 */
type PublisherHandler struct {
	I_PublisherHandler
	subscribers map[common.NetConf]*SubscriberClient
}

/**
 * The SubscriberClient is the client.Client used to communicate to the Subscriber
 */
type SubscriberClient struct {
	*client.Client
}

/**
 * Update the subscriber
 */
func (s *SubscriberClient) Update(req any, res any) error {
	return s.Call("Subscriber.Update", req, res)
}

/**
 * InitPublisher creates Server
 * @param network is the server's configuration
 * @return the initialized Server or error
 */
func InitPublisher(network common.NetConf) (*Publisher, error) {
	server, err := server.InitServer(network)
	if err != nil {
		return nil, err
	}

	handler := &PublisherHandler{subscribers: make(map[common.NetConf]*SubscriberClient)}

	pub := &Publisher{
		Server:           server,
		PublisherHandler: handler,
	}

	err = pub.Register(handler)
	if err != nil {
		return nil, err
	}

	return pub, nil
}

/**
 * Start the Server that was initialized with a netork config
 * You might consider starting the server in a goroutine
 * @return potential networking errors
 */
func (p *Publisher) Start() error {
	return p.Server.Start()
}

/**
 * Subscribe will add a SubscriberClient to the list of subscribers
 * @param req is the request containig networking config to initialize SubscriberClient
 * @param res is the response that will give Ok=true if subscription was effective
 * @return a potential network error
 */
func (p *PublisherHandler) Subscribe(req *common.SubscribeRequest, res *common.SubscribeResponse) error {
	_, exist := p.subscribers[req.Subscriber]
	if exist {
		return nil
	}

	newSubscriber := &SubscriberClient{Client: client.InitClient(req.Subscriber)}

	p.subscribers[req.Subscriber] = newSubscriber

	// TODO: remove dial and implement reconnection instead
	err := p.subscribers[req.Subscriber].Dial()
	if err != nil {
		return err
	}

	res.Ok = true

	return nil
}

/**
 * Unsubscribe will remove a SubscriberClient from the list of subscribers
 * @param req is the request containig networking config of SubscriberClient to remove
 * @param res is the response that will give Ok=true if unsubscription was effective
 * @return a potential network error
 */
func (p *PublisherHandler) Unsubscribe(req *common.SubscribeRequest, res *common.SubscribeResponse) error {
	delete(p.subscribers, req.Subscriber)
	res.Ok = true

	return nil
}

/**
 * Publish will cycle through all subscribers and send them the message
 * @param req is the request
 * @param res is the response
 * @return a potential network error
 */
func (p *Publisher) Publish(req any, res any) {
	for _, sub := range p.subscribers {
		err := sub.Call("SubscriberHandler.Update", &req, &res)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

/**
 * SUBSCRIBER :
 */

/**
 * The basic Subscriber functions
 */
type I_Subscriber interface {
	Subscribe(Publisher) error
	Unsubscribe(Publisher) error
}

/**
 * Subscriber is a ClientServer that can subscribe to Publishers and recieve updates
 */
type Subscriber struct {
	I_Subscriber
	*clientServer.ClientServer
	*SubscriberHandler
}

/**
 * The basic SubscriberHandler functions
 */
type I_SubscriberHandler interface {
	Update(any, any) error
}

/**
 * The SubscriberHandler has a channel to forward the update message
 */
type SubscriberHandler struct {
	I_SubscriberHandler
	msgChan chan any
}

/**
 * InitSubscriber creates a Subscriber, inheriting from ClientServer
 * @param selfNetwork is the server's network config
 * @param remoteNetwork is the remote server's network config
 * @return the initialized Subscriber or error
 */
func InitSubscriber(selfNetwork common.NetConf, remoteNetwork common.NetConf) (*Subscriber, error) {
	clientServer, err := clientServer.InitClientServer(selfNetwork, remoteNetwork)
	if err != nil {
		return nil, err
	}

	handler := &SubscriberHandler{msgChan: make(chan any)}

	subscriber := &Subscriber{
		ClientServer:      clientServer,
		SubscriberHandler: handler,
	}

	err = subscriber.Register(handler)
	if err != nil {
		return nil, err
	}

	return subscriber, nil
}

/**
 * Start the ClientServer that was initialized with a netork config
 * You might consider starting the server in a goroutine
 * @return potential networking errors
 */
func (s *Subscriber) Start() error {
	return s.Server.Start()
}

/**
 * Subscribe to the desired publisher
 * @param publisher is the target publisher
 * @return potential networking or subscription errors
 */
func (s *Subscriber) Subscribe(publisher common.NetConf) error {
	req := common.SubscribeRequest{Subscriber: s.Server.NetConf}
	res := common.SubscribeResponse{}
	err := s.Call("PublisherHandler.Subscribe", &req, &res)
	if err != nil {
		return err
	}

	if !res.Ok {
		return fmt.Errorf("publisher %+v could not subscribe %+v", publisher, s.Server.NetConf)
	}

	return nil
}

/**
 * Unsubscribe from the desired publisher
 * @param publisher is the target publisher
 * @return potential networking or unsubscription errors
 */
func (s *Subscriber) Unsubscribe(publisher common.NetConf) error {
	req := common.SubscribeRequest{Subscriber: s.Server.NetConf}
	res := common.SubscribeResponse{}
	err := s.Call("PublisherHandler.Unsubscribe", &req, &res)
	if err != nil {
		return err
	}

	if !res.Ok {
		return fmt.Errorf("publisher %+v could not unsubscribe %+v", publisher, s.Server.NetConf)
	}

	return nil
}

/**
 * Update will forward the incoming request data to the message channel
 */
func (s *SubscriberHandler) Update(req *any, res *any) error {
	s.msgChan <- *req

	return nil
}

/**
 * Message channel getter
 */
func (s *Subscriber) Chan() chan any {
	return s.msgChan
}

/**
 * Close the running server
 * Same as Stop()
 */
func (s *Subscriber) Close() error {
	close(s.msgChan)
	return s.ClientServer.Close()
}

/**
 * Stop the running server
 * Same as Close()
 */
func (s *Subscriber) Stop() error {
	return s.Close()
}
