package client

import (
	"testing"
)

func TestClientCall(t *testing.T) {
	// Create a sample NetConf for testing
	netConf := NetConf{
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "1234",
	}
	var err error

	server, err := InitServer(netConf)
	if err != nil {
		t.Error(err)
	}

	go func(){
		err = server.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	// Initialize the client
	client := InitClient(netConf)
	
	err = client.Dial()
	if err != nil {
		t.Error(err)
	}

	err = client.Ping()
	if err != nil {
		t.Error(err)
	}

	err = client.Close()
	if err != nil {
		t.Error(err)
	}

	server.Stop()
}

func TestClientReconnect(t *testing.T) {
	netConf := NetConf{
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "1234",
	}
	var err error

	server, err := InitServer(netConf)
	if err != nil {
		t.Error(err)
	}

	go func(){
		err = server.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	// Initialize the client
	client := InitClient(netConf)
	
	// err = client.Dial()
	// if err != nil {
	// 	t.Error(err)
	// }

	err = client.Ping()
	if _, ok := err.(MicronetReconnectTimeoutError); ok {
		t.Error(err)
	}

	err = client.Close()
	if err != nil {
		t.Error(err)
	}

	server.Stop()
}

func TestClientCallTimeout(t *testing.T) {
	netConf := NetConf{
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "1234",
	}
	var err error

	// Initialize the client
	client := InitClient(netConf)
	
	// err = client.Dial()
	// if err != nil {
	// 	t.Error(err)
	// }

	err = client.Ping()
	e, ok := err.(*MicronetReconnectTimeoutError)
	if !ok && e != nil{
		t.Error(err)
	}
}

func TestClientGo(t *testing.T) {
	// Create a sample NetConf for testing
	netConf := NetConf{
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "1234",
	}
	var err error

	server, err := InitServer(netConf)
	if err != nil {
		t.Error(err)
	}

	go func(){
		err = server.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	// Initialize the client
	client := InitClient(netConf)
	
	err = client.Dial()
	if err != nil {
		t.Error(err)
	}

	request := Ping{Data: "PING"}
	response := Pong{}
	call := client.Go("PingHandler.Ping", &request, &response, nil)
	if call.Error != nil {
		return
	}
	<-call.Done

	err = client.Close()
	if err != nil {
		t.Error(err)
	}

	server.Stop()
}



func TestClientGoTimeout(t *testing.T) {
	netConf := NetConf{
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "1234",
	}
	var err error

	// Initialize the client
	client := InitClient(netConf)
	
	// err = client.Dial()
	// if err != nil {
	// 	t.Error(err)
	// }

	request := Ping{Data: "PING"}
	response := Pong{}
	call := client.Go("PingHandler.Ping", &request, &response, nil)
	e, ok := call.Error.(*MicronetReconnectTimeoutError)
	if !ok && e != nil{
		t.Error(err)
	}
}