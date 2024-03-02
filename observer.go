package micronet

import (
	"fmt"
	"log"
)

/**
 * PUBLISHER :
 */

 type I_Publisher interface {
	Publish(string)
}

type Publisher struct {
	I_Publisher
	*Server
	*PublisherHandler
}

type I_PublisherHandler interface {
	Subscribe(*SubscribeRequest, *SubscribeResponse) (error)
	Unsubscribe(*SubscribeRequest, *SubscribeResponse) (error)
}

type PublisherHandler struct {
	I_PublisherHandler
	subscribers map[NetConf]*SubscriberClient
}

type SubscriberClient struct {
	*Client
}
func (s *SubscriberClient) Update(req *any, res *any) (error) {
	return s.Call("Subscriber.Update", req, res)
}

func InitPublisher(selfNetwork NetConf) (*Publisher, error) {
	server, err := InitServer(selfNetwork)
	if err != nil {
		return nil, err
	}
	
	handler := &PublisherHandler{subscribers: make(map[NetConf]*SubscriberClient)}
	
	pub := &Publisher{
		Server: server,
		PublisherHandler: handler,
	}

	err = pub.Register(handler)
	if err != nil {
		return nil, err
	}
	
	return pub, nil
}

func (p *PublisherHandler) Subscribe(req *SubscribeRequest, res *SubscribeResponse) (error) {
	_, exist := p.subscribers[req.Subscriber]
	if exist {
		return nil
	}

	newSubscriber := &SubscriberClient{Client: InitClient(req.Subscriber)}

	p.subscribers[req.Subscriber] = newSubscriber

	// TODO: remove dial and implement reconnection instead
	err := p.subscribers[req.Subscriber].Dial()
	if err != nil {
		return err
	}

	res.Ok = true

	return nil
}

func (p *PublisherHandler) Unsubscribe(req *SubscribeRequest, res *SubscribeResponse) (error) {
	delete(p.subscribers, req.Subscriber)
	res.Ok = true

	return nil
}

func (p *Publisher) Publish(msg any, res any) {
	for _, sub := range p.subscribers {
		err := sub.Call("SubscriberHandler.Update", &msg, &res)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

/**
 * SUBSCRIBER :
 */

 type I_Subscriber interface {
	Subscribe(Publisher) (error)
	Unsubscribe(Publisher) (error)
}

type Subscriber struct {
	I_Subscriber
	*ClientServer
	*SubscriberHandler
}

type I_SubscriberHandler interface {
	Update(any, any) (error)
}

type SubscriberHandler struct {
	I_SubscriberHandler
	msgChan chan any
}

func InitSubscriber(selfNetwork NetConf, remoteNetwork NetConf) (*Subscriber, error) {
	clientServer, err := InitClientServer(selfNetwork, remoteNetwork)
	if err != nil {
		return nil, err
	}

	handler := &SubscriberHandler{msgChan: make(chan any)}

	subscriber := &Subscriber{
		ClientServer: clientServer,
		SubscriberHandler: handler,
	}

	err = subscriber.Register(handler)
	if err != nil {
		return nil, err
	}

	return subscriber, nil
}

func (s *Subscriber) Subscribe(publisher NetConf) (error) {
	req := SubscribeRequest{Subscriber: s.Server.NetConf}
	res := SubscribeResponse{}
	err := s.Call("PublisherHandler.Subscribe", &req, &res)
	if err != nil {
		return err
	}

	if !res.Ok {
		return fmt.Errorf("publisher %+v could not subscribe %+v", publisher, s.Server.NetConf)
	}
	
	return nil
}

func (s *Subscriber) Unsubscribe(publisher NetConf) (error) {
	req := SubscribeRequest{Subscriber: s.Server.NetConf}
	res := SubscribeResponse{}
	err := s.Call("PublisherHandler.Unsubscribe", &req, &res)
	if err != nil {
		return err
	}

	if !res.Ok {
		return fmt.Errorf("publisher %+v could not unsubscribe %+v", publisher, s.Server.NetConf)
	}
	
	return nil
}

func (s *SubscriberHandler) Update(req *any, res *any) (error) {
	s.msgChan <- *req

	return nil
}

func (s *Subscriber) Chan() (chan any) {
	return s.msgChan
}

func (s *Subscriber) Close() (error) {
	close(s.msgChan)
	return s.ClientServer.Close()
}

func (s *Subscriber) Stop() (error) {
	return s.Close()
}