package micronet

import (
	"testing"
	"time"
)

func TestObserver(t *testing.T) {
	// Create a sample NetConf for testing
	pubNetConf := NetConf{
		Name:     "Publisher",
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "1234",
	}
	
	subNetConf := NetConf{
		Name:     "Subscriber",
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "4321",
	}
	var err error

	pub, err := InitPublisher(pubNetConf)
	if err != nil {
		t.Error(err)
	}
	go func() {
		err = pub.Server.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	sub, err := InitSubscriber(subNetConf, pubNetConf)
	if err != nil {
		t.Error(err)
	}
	go func() {
		err = sub.Server.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	err = sub.Dial()
	if err != nil {
		t.Error(err)
	}

	err = sub.Subscribe(pubNetConf)
	if err != nil {
		t.Error(err)
	}

	msg := "message"
	res := ""
	go pub.Publish(&msg, &res)

	time.Sleep(time.Second)

	go func(){
		msg := (<-sub.Chan()).(string)
		if msg != "message" {
			t.Errorf("incorrect message : '%s'", msg)
		}
	}()

	time.Sleep(time.Second)

	err = sub.Unsubscribe(pubNetConf)
	if err != nil {
		t.Error(err)
	}


	err = sub.Stop()
	if err != nil {
		t.Error(err)
	}
	
	pub.Stop()

}
