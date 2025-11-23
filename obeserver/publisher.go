package common

import (
	"log"

	"micronet/client"
	"micronet/common"
	"micronet/server"
)

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
