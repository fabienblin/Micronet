package micronet

import (
	"testing"
)

func TestClientPing(t *testing.T) {
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