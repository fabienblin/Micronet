package common

import (
	"testing"
	"time"
	
	"micronet/common"
)

func TestObserver(t *testing.T) {
	// Create a sample common.NetConf for testing
	pubNetConf := common.NetConf{
		Name:     "Publisher",
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "1234",
	}

	subNetConf := common.NetConf{
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
		err = pub.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	sub, err := InitSubscriber(subNetConf, pubNetConf)
	if err != nil {
		t.Error(err)
	}
	go func() {
		err = sub.Start()
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

	go func() {
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

func TestObserver2(t *testing.T) {
	// Create a sample common.NetConf for testing
	pubNetConf := common.NetConf{
		Name:     "Publisher",
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "1234",
	}

	subNetConf1 := common.NetConf{
		Name:     "Subscriber",
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "4321",
	}

	subNetConf2 := common.NetConf{
		Name:     "Subscriber",
		Protocol: "tcp",
		Ip:       "localhost",
		Port:     "4322",
	}
	var err error

	pub, err := InitPublisher(pubNetConf)
	if err != nil {
		t.Error(err)
	}
	go func() {
		err = pub.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	sub1, err := InitSubscriber(subNetConf1, pubNetConf)
	if err != nil {
		t.Error(err)
	}
	go func() {
		err = sub1.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	sub2, err := InitSubscriber(subNetConf2, pubNetConf)
	if err != nil {
		t.Error(err)
	}
	go func() {
		err = sub2.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	err = sub1.Dial()
	if err != nil {
		t.Error(err)
	}

	err = sub2.Dial()
	if err != nil {
		t.Error(err)
	}

	err = sub1.Subscribe(pubNetConf)
	if err != nil {
		t.Error(err)
	}

	err = sub2.Subscribe(pubNetConf)
	if err != nil {
		t.Error(err)
	}

	msg := "message"
	res := ""
	go pub.Publish(&msg, &res)

	time.Sleep(time.Second)

	go func() {
		msg := (<-sub1.Chan()).(string)
		if msg != "message" {
			t.Errorf("incorrect message : '%s'", msg)
		}
	}()

	time.Sleep(time.Second)

	go func() {
		msg := (<-sub2.Chan()).(string)
		if msg != "message" {
			t.Errorf("incorrect message : '%s'", msg)
		}
	}()

	time.Sleep(time.Second)

	err = sub1.Unsubscribe(pubNetConf)
	if err != nil {
		t.Error(err)
	}

	err = sub2.Unsubscribe(pubNetConf)
	if err != nil {
		t.Error(err)
	}

	err = sub1.Stop()
	if err != nil {
		t.Error(err)
	}

	err = sub2.Stop()
	if err != nil {
		t.Error(err)
	}

	pub.Stop()

}
