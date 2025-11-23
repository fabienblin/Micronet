package clientServer

import (
	"testing"
	
	"micronet/common"
)

func TestClientServerPing(t *testing.T) {
	// Create a sample common.NetConf for testing
	netConf1 := common.NetConf{
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "1234",
	}
	netConf2 := common.NetConf{
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "4321",
	}
	var err error

	srv1, err := InitClientServer(netConf1, netConf2)
	if err != nil {
		t.Error(err)
	}
	go func() {
		err = srv1.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	srv2, err := InitClientServer(netConf2, netConf1)
	if err != nil {
		t.Error(err)
	}
	go func() {
		err = srv2.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	err = srv1.Dial()
	if err != nil {
		t.Error(err)
	}

	err = srv2.Dial()
	if err != nil {
		t.Error(err)
	}

	err = srv1.Ping()
	if err != nil {
		t.Error(err)
	}

	err = srv2.Ping()
	if err != nil {
		t.Error(err)
	}

	err = srv1.Stop()
	if err != nil {
		t.Error(err)
	}

	err = srv2.Stop()
	if err != nil {
		t.Error(err)
	}
}
