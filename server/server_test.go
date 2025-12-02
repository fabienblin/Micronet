package server

import (
	"testing"

	"micronet/common"

	"github.com/stretchr/testify/assert"
)

func TestServerLifecycle(t *testing.T) {
	t.Run("Nominal case", func(t *testing.T) {
		netConf := common.NetConf{
			Protocol: "tcp",
			Port:     "12345",
			Ip:       "localhost",
		}
	
		server, errNew := NewServer(netConf)
		if !assert.NoError(t, errNew) {
			t.FailNow()
		}
	
		go func(){
			errStart := server.Start()
			if !assert.NoError(t, errStart) {
				return
			}
		}()
	
		server.Stop()
	})
}

// Add more test functions as needed
