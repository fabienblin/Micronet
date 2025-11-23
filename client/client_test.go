package client

import (
	"testing"

	"micronet/common"
	"micronet/server"
)

func TestClientCall(t *testing.T) {
	// Create a sample common.NetConf for testing
	netConf := common.NetConf{
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "1234",
	}
	var err error

	server, err := server.InitServer(netConf)
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
	netConf := common.NetConf{
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "1234",
	}
	var err error

	server, err := server.InitServer(netConf)
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
	if _, ok := err.(common.MicronetReconnectTimeoutError); ok {
		t.Error(err)
	}

	err = client.Close()
	if err != nil {
		t.Error(err)
	}

	server.Stop()
}

func TestClientCallTimeout(t *testing.T) {
	netConf := common.NetConf{
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
	e, ok := err.(*common.MicronetReconnectTimeoutError)
	if !ok && e != nil{
		t.Error(err)
	}
}

func TestClientGo(t *testing.T) {
	// Create a sample common.NetConf for testing
	netConf := common.NetConf{
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "1234",
	}
	var err error

	server, err := server.InitServer(netConf)
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

	request := common.Ping{Data: "PING"}
	response := common.Pong{}
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
	netConf := common.NetConf{
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

	request := common.Ping{Data: "PING"}
	response := common.Pong{}
	call := client.Go("PingHandler.Ping", &request, &response, nil)
	e, ok := call.Error.(*common.MicronetReconnectTimeoutError)
	if !ok && e != nil{
		t.Error(err)
	}
}