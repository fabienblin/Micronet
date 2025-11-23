package clientServer

import (
	"net/rpc"
)

// MockClientServer implements I_ClientServer for unit testing.
type MockClientServer struct {
	// ----- CLIENT methods -----
	DialFunc  func() error
	CallFunc  func(serviceMethod string, args any, reply any) error
	GoFunc    func(serviceMethod string, args any, reply any, done chan *rpc.Call) *rpc.Call
	CloseFunc func() error
	PingFunc  func() error

	// ----- SERVER methods -----
	RegisterFunc func(rcvr any) error
	StartFunc    func() error
	StopFunc     func()
}

// Ensure interface compliance.
var _ I_ClientServer = (*MockClientServer)(nil)

// ----- CLIENT SIDE -----

func (m *MockClientServer) Dial() error {
	if m.DialFunc != nil {
		return m.DialFunc()
	}
	return nil
}

func (m *MockClientServer) Call(serviceMethod string, args any, reply any) error {
	if m.CallFunc != nil {
		return m.CallFunc(serviceMethod, args, reply)
	}
	return nil
}

func (m *MockClientServer) Go(serviceMethod string, args any, reply any, done chan *rpc.Call) *rpc.Call {
	if m.GoFunc != nil {
		return m.GoFunc(serviceMethod, args, reply, done)
	}
	return &rpc.Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}
}

func (m *MockClientServer) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *MockClientServer) Ping() error {
	if m.PingFunc != nil {
		return m.PingFunc()
	}
	return nil
}

// ----- SERVER SIDE -----

func (m *MockClientServer) Register(rcvr any) error {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(rcvr)
	}
	return nil
}

func (m *MockClientServer) Start() error {
	if m.StartFunc != nil {
		return m.StartFunc()
	}
	return nil
}

func (m *MockClientServer) Stop() {
	if m.StopFunc != nil {
		m.StopFunc()
	}
}
