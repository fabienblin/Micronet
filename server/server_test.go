package server

import (
	"testing"
	
	"micronet/common"
)

func TestServerLifecycle(t *testing.T) {
	// Create a sample NetConf for testing
	netConf := common.NetConf{
		Protocol: "tcp",
		Port:     "12345",
		Ip:       "localhost",
	}

	// Initialize the server
	server, err := InitServer(netConf)
	if err != nil {
		t.Error(err)
	}

	// Start the server in a goroutine
	go server.Start()

	// Perform any necessary assertions or tests here

	// Stop the server
	server.Stop()

	// Perform any necessary assertions or tests here

	// Additional assertions can be added based on your specific use case
}

// Add more test functions as needed
