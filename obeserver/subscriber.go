package common

import (
	"fmt"

	"micronet/clientServer"
	"micronet/common"
)

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
	clientServer, err := clientServer.NewClientServer(selfNetwork, remoteNetwork)
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
