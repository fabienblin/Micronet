package client

import (
	"net/rpc"
)

type MockClient struct {
	DialFunc  func() error
	CallFunc  func(serviceMethod string, args any, reply any) error
	GoFunc    func(serviceMethod string, args any, reply any, done chan *rpc.Call) *rpc.Call
	CloseFunc func() error
	PingFunc  func() error
}

// Ensure MockClient implements I_Client
var _ I_Client = (*MockClient)(nil)

func (m *MockClient) Dial() error {
	if m.DialFunc != nil {
		return m.DialFunc()
	}
	return nil
}

func (m *MockClient) Call(serviceMethod string, args any, reply any) error {
	if m.CallFunc != nil {
		return m.CallFunc(serviceMethod, args, reply)
	}
	return nil
}

func (m *MockClient) Go(serviceMethod string, args any, reply any, done chan *rpc.Call) *rpc.Call {
	if m.GoFunc != nil {
		return m.GoFunc(serviceMethod, args, reply, done)
	}
	return &rpc.Call{ServiceMethod: serviceMethod, Args: args, Reply: reply, Done: done}
}

func (m *MockClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *MockClient) Ping() error {
	if m.PingFunc != nil {
		return m.PingFunc()
	}
	return nil
}
