package common

import (
	"micronet/common"
	"net/rpc"
)

// =========================================================
//                    PUBLISHER MOCKS
// =========================================================

// MockPublisher implements I_Publisher
type MockPublisher struct {
	PublishFunc func(req string)
}

var _ I_Publisher = (*MockPublisher)(nil)

func (m *MockPublisher) Publish(req string) {
	if m.PublishFunc != nil {
		m.PublishFunc(req)
	}
}

// ---------------------------------------------------------
// MockPublisherHandler implements I_PublisherHandler
// ---------------------------------------------------------

type MockPublisherHandler struct {
	SubscribeFunc   func(req *common.SubscribeRequest, res *common.SubscribeResponse) error
	UnsubscribeFunc func(req *common.SubscribeRequest, res *common.SubscribeResponse) error
}

var _ I_PublisherHandler = (*MockPublisherHandler)(nil)

func (m *MockPublisherHandler) Subscribe(req *common.SubscribeRequest, res *common.SubscribeResponse) error {
	if m.SubscribeFunc != nil {
		return m.SubscribeFunc(req, res)
	}
	res.Ok = true
	return nil
}

func (m *MockPublisherHandler) Unsubscribe(req *common.SubscribeRequest, res *common.SubscribeResponse) error {
	if m.UnsubscribeFunc != nil {
		return m.UnsubscribeFunc(req, res)
	}
	res.Ok = true
	return nil
}

// =========================================================
//                   SUBSCRIBER MOCKS
// =========================================================

// MockSubscriber implements I_Subscriber
type MockSubscriber struct {
	SubscribeFunc   func(publisher Publisher) error
	UnsubscribeFunc func(publisher Publisher) error
}

var _ I_Subscriber = (*MockSubscriber)(nil)

func (m *MockSubscriber) Subscribe(publisher Publisher) error {
	if m.SubscribeFunc != nil {
		return m.SubscribeFunc(publisher)
	}
	return nil
}

func (m *MockSubscriber) Unsubscribe(publisher Publisher) error {
	if m.UnsubscribeFunc != nil {
		return m.UnsubscribeFunc(publisher)
	}
	return nil
}

// ---------------------------------------------------------
// MockSubscriberHandler implements I_SubscriberHandler
// ---------------------------------------------------------

type MockSubscriberHandler struct {
	UpdateFunc func(req any, res any) error
}

var _ I_SubscriberHandler = (*MockSubscriberHandler)(nil)

func (m *MockSubscriberHandler) Update(req any, res any) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(req, res)
	}
	return nil
}

// =========================================================
//            SubscriberClient MOCK (RPC Update)
// =========================================================

// MockSubscriberClient simulates the RPC client used in Publish()
type MockSubscriberClient struct {
	UpdateFunc func(req any, res any) error
}

func (m *MockSubscriberClient) Update(req any, res any) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(req, res)
	}
	return nil
}

// To satisfy how Publisher uses it:
func (m *MockSubscriberClient) Call(method string, args any, reply any) error {
	// Only SubscriberHandler.Update is ever used in the real code
	if m.UpdateFunc != nil {
		return m.UpdateFunc(args, reply)
	}
	return nil
}

// Prevents compile errors (fake field for compatibility only)
type _ = rpc.Call
