package client

import (
	"micronet/common"
	"net"
	"net/rpc"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const ListenReadynessDuration time.Duration = time.Millisecond * 100
const MockMethodResponseValue int = 42

var netConf = common.NetConf{
	Name:     "mock",
	Protocol: "tcp",
	Ip:       "127.0.0.1",
	Port:     "12345",
}

type MockService struct{}

func (s *MockService) MockMethod(req bool, resp *int) error {
	*resp = MockMethodResponseValue
	return nil
}

func TestClient_NewClient(t *testing.T) {
	t.Run("Client connection success", func(t *testing.T) {
		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}
		defer listener.Close()
		go rpc.Accept(listener)
		time.Sleep(ListenReadynessDuration)

		_, err := NewClient(netConf)
		assert.NoError(t, err)

	})

	t.Run("Client connection fail", func(t *testing.T) {
		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}
		defer listener.Close()
		go rpc.Accept(listener)
		time.Sleep(ListenReadynessDuration)

		var badNetConf = common.NetConf{
			Name:     "wrongPort",
			Protocol: "tcp",
			Ip:       "127.0.0.1",
			Port:     "54321",
		}

		_, err := NewClient(badNetConf)
		assert.Error(t, err)

	})
}

func TestClient_Call(t *testing.T) {
	t.Run("Invalid service method", func(t *testing.T) {
		mockServer := rpc.NewServer()

		mockService := &MockService{}
		err := mockServer.Register(mockService)
		if err != nil {
			t.Fatal(err)
		}

		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}
		defer listener.Close()
		go mockServer.Accept(listener)
		time.Sleep(ListenReadynessDuration)

		client, errDial := NewClient(netConf)
		assert.NoError(t, errDial)

		errCall := client.Call("unRegisteredMethod", nil, nil)

		assert.Error(t, errCall)
	})

	t.Run("Invalid request type", func(t *testing.T) {
		mockServer := rpc.NewServer()

		mockService := &MockService{}
		err := mockServer.Register(mockService)
		if err != nil {
			t.Fatal(err)
		}

		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}
		defer listener.Close()
		go mockServer.Accept(listener)
		time.Sleep(ListenReadynessDuration)

		client, errDial := NewClient(netConf)
		assert.NoError(t, errDial)

		request := "true"
		var response int

		errCall := client.Call("MockService.MockMethod", request, &response)

		assert.Error(t, errCall)
	})

	t.Run("Invalid response type", func(t *testing.T) {
		mockServer := rpc.NewServer()

		mockService := &MockService{}
		err := mockServer.Register(mockService)
		if err != nil {
			t.Fatal(err)
		}

		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}
		defer listener.Close()
		go mockServer.Accept(listener)
		time.Sleep(ListenReadynessDuration)

		client, errDial := NewClient(netConf)
		assert.NoError(t, errDial)

		request := true
		var response string

		errCall := client.Call("MockService.MockMethod", request, &response)

		assert.Error(t, errCall)
	})

	t.Run("Connection lost", func(t *testing.T) {
		mockServer := rpc.NewServer()

		mockService := &MockService{}
		err := mockServer.Register(mockService)
		if err != nil {
			t.Fatal(err)
		}

		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}

		var serverConn net.Conn
		var serverConnMu sync.Mutex
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					return
				}
				serverConnMu.Lock()
				serverConn = conn
				serverConnMu.Unlock()
				go mockServer.ServeConn(conn)
			}
		}()
		time.Sleep(ListenReadynessDuration)

		client, errDial := NewClient(netConf)
		client.SetReconnectionConf(0, 0)
		assert.NoError(t, errDial)

		request := true
		var response int

		// ---- Force next Call() to fail ----
		listener.Close()
		serverConnMu.Lock()
		if serverConn != nil {
			serverConn.Close()
		}
		serverConnMu.Unlock()
		// -----------------------------------

		errCall := client.Call("MockService.MockMethod", request, &response)

		assert.Error(t, errCall)
	})

	t.Run("Reconnection success", func(t *testing.T) {
		mockServer := rpc.NewServer()

		mockService := &MockService{}
		err := mockServer.Register(mockService)
		if err != nil {
			t.Fatal(err)
		}

		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}

		var serverConn net.Conn
		var serverConnMu sync.Mutex
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					return
				}
				serverConnMu.Lock()
				serverConn = conn
				serverConnMu.Unlock()
				go mockServer.ServeConn(conn)
			}
		}()
		time.Sleep(ListenReadynessDuration)

		client, errDial := NewClient(netConf)
		assert.NoError(t, errDial)

		request := true
		var response int

		// ---- Force next Call() to fail ----
		listener.Close()
		serverConnMu.Lock()
		if serverConn != nil {
			serverConn.Close()
		}
		serverConnMu.Unlock()
		// ------- Create a new server -------
		secondMockServer := rpc.NewServer()

		secondMockService := &MockService{}
		errSecondRegister := secondMockServer.Register(secondMockService)
		if errSecondRegister != nil {
			t.Fatal(errSecondRegister)
		}

		secondListener, errSecondListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errSecondListen != nil {
			t.Fatal(errSecondListen)
		}
		defer secondListener.Close()
		go secondMockServer.Accept(secondListener)
		time.Sleep(ListenReadynessDuration)
		// -----------------------------------

		errCall := client.Call("MockService.MockMethod", request, &response)

		assert.NoError(t, errCall)
	})

	t.Run("Nominal case", func(t *testing.T) {
		mockServer := rpc.NewServer()

		mockService := &MockService{}
		err := mockServer.Register(mockService)
		if err != nil {
			t.Fatal(err)
		}

		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}
		defer listener.Close()
		go mockServer.Accept(listener)
		time.Sleep(ListenReadynessDuration)

		client, errDial := NewClient(netConf)
		assert.NoError(t, errDial)

		request := true
		var response int

		errCall := client.Call("MockService.MockMethod", request, &response)

		assert.NoError(t, errCall)
		assert.Equal(t, MockMethodResponseValue, response)
	})
}

func TestClient_Go(t *testing.T) {
	t.Run("Invalid service method", func(t *testing.T) {
		mockServer := rpc.NewServer()

		mockService := &MockService{}
		err := mockServer.Register(mockService)
		if err != nil {
			t.Fatal(err)
		}

		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}
		defer listener.Close()
		go mockServer.Accept(listener)
		time.Sleep(ListenReadynessDuration)

		client, errDial := NewClient(netConf)
		assert.NoError(t, errDial)

		doneChan := client.Go("unRegisteredMethod", nil, nil, nil)
		<-doneChan.Done

		assert.Error(t, doneChan.Error)
	})

	t.Run("Invalid request type", func(t *testing.T) {
		mockServer := rpc.NewServer()

		mockService := &MockService{}
		err := mockServer.Register(mockService)
		if err != nil {
			t.Fatal(err)
		}

		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}
		defer listener.Close()
		go mockServer.Accept(listener)
		time.Sleep(ListenReadynessDuration)

		client, errDial := NewClient(netConf)
		assert.NoError(t, errDial)

		request := "true"
		var response int

		doneChan := client.Go("MockService.MockMethod", request, &response, nil)
		<-doneChan.Done

		assert.Error(t, doneChan.Error)
	})

	t.Run("Invalid response type", func(t *testing.T) {
		mockServer := rpc.NewServer()

		mockService := &MockService{}
		err := mockServer.Register(mockService)
		if err != nil {
			t.Fatal(err)
		}

		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}
		defer listener.Close()
		go mockServer.Accept(listener)
		time.Sleep(ListenReadynessDuration)

		client, errDial := NewClient(netConf)
		assert.NoError(t, errDial)

		request := true
		var response string

		doneChan := client.Go("MockService.MockMethod", request, &response, nil)
		<-doneChan.Done

		assert.Error(t, doneChan.Error)
	})

	t.Run("Connection lost", func(t *testing.T) {
		mockServer := rpc.NewServer()

		mockService := &MockService{}
		err := mockServer.Register(mockService)
		if err != nil {
			t.Fatal(err)
		}

		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}

		var serverConn net.Conn
		var serverConnMu sync.Mutex
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					return
				}
				serverConnMu.Lock()
				serverConn = conn
				serverConnMu.Unlock()
				go mockServer.ServeConn(conn)
			}
		}()
		time.Sleep(ListenReadynessDuration)

		client, errDial := NewClient(netConf)
		client.SetReconnectionConf(0, 0)
		assert.NoError(t, errDial)

		request := true
		var response int

		// ---- Force next Call() to fail ----
		listener.Close()
		serverConnMu.Lock()
		if serverConn != nil {
			serverConn.Close()
		}
		serverConnMu.Unlock()
		// -----------------------------------

		doneChan := client.Go("MockService.MockMethod", request, &response, nil)
		<-doneChan.Done

		assert.Error(t, doneChan.Error)
	})

	t.Run("Reconnection success", func(t *testing.T) {
		mockServer := rpc.NewServer()

		mockService := &MockService{}
		err := mockServer.Register(mockService)
		if err != nil {
			t.Fatal(err)
		}

		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}

		var serverConn net.Conn
		var serverConnMu sync.Mutex
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					return
				}
				serverConnMu.Lock()
				serverConn = conn
				serverConnMu.Unlock()
				go mockServer.ServeConn(conn)
			}
		}()
		time.Sleep(ListenReadynessDuration)

		client, errDial := NewClient(netConf)
		assert.NoError(t, errDial)

		request := true
		var response int

		// ---- Force next Call() to fail ----
		listener.Close()
		serverConnMu.Lock()
		if serverConn != nil {
			serverConn.Close()
		}
		serverConnMu.Unlock()
		// ------- Create a new server -------
		secondMockServer := rpc.NewServer()

		secondMockService := &MockService{}
		errSecondRegister := secondMockServer.Register(secondMockService)
		if errSecondRegister != nil {
			t.Fatal(errSecondRegister)
		}

		secondListener, errSecondListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errSecondListen != nil {
			t.Fatal(errSecondListen)
		}
		defer secondListener.Close()
		go secondMockServer.Accept(secondListener)
		time.Sleep(ListenReadynessDuration)
		// -----------------------------------

		doneChan := client.Go("MockService.MockMethod", request, &response, nil)
		<-doneChan.Done

		assert.NoError(t, doneChan.Error)
	})

	t.Run("Nominal case", func(t *testing.T) {
		mockServer := rpc.NewServer()

		mockService := &MockService{}
		err := mockServer.Register(mockService)
		if err != nil {
			t.Fatal(err)
		}

		listener, errListen := net.Listen(netConf.Protocol, netConf.Ip+":"+netConf.Port)
		if errListen != nil {
			t.Fatal(errListen)
		}
		defer listener.Close()
		go mockServer.Accept(listener)
		time.Sleep(ListenReadynessDuration)

		client, errDial := NewClient(netConf)
		assert.NoError(t, errDial)

		request := true
		var response int

		doneChan := client.Go("MockService.MockMethod", request, &response, nil)
		<-doneChan.Done

		assert.NoError(t, doneChan.Error)
		assert.Equal(t, MockMethodResponseValue, response)
	})
}
