package server

// MockServer is a test double for I_Server.
type MockServer struct {
	// --- Optional behavior to define in your tests ---

	RegisterFunc func(rcvr any) error
	StartFunc    func() error
	StopFunc     func()
}

// Ensure MockServer implements I_Server.
var _ I_Server = (*MockServer)(nil)

func (m *MockServer) Register(rcvr any) error {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(rcvr)
	}
	return nil
}

func (m *MockServer) Start() error {
	if m.StartFunc != nil {
		return m.StartFunc()
	}
	return nil
}

func (m *MockServer) Stop() {
	if m.StopFunc != nil {
		m.StopFunc()
	}
}
