package tcp

import (
	"net"
	"sync"
	"testing"
	"time"

	"go-common/app/infra/databus/conf"
)

var (
	_testGroup = "test-consumer-group"
	_testTopic = "test_topic"
	_testAddr  = "172.22.33.174:9092"

	_testConfig = &conf.Kafka{
		Cluster: _testGroup,
		Brokers: []string{_testAddr},
	}
)

func TestNewSub(t *testing.T) {
	var (
		mu  sync.Mutex
		err error
	)
	subs := []*Sub{}
	c, _ := net.Dial("tcp", _testAddr)
	sub, err := NewSub(newConn(c, time.Second, time.Second), _testGroup, _testTopic, "", _testConfig, 100)
	if err != nil {
		t.Fatal(err)
	}
	subs = append(subs, sub)
	time.Sleep(time.Second * 5)
	go func() {
		for i := 0; i < 200; i++ {
			sub, err := NewSub(newConn(c, time.Second, time.Second), _testGroup, _testTopic, "", _testConfig, 100)
			if err != nil {
				t.Errorf("NewSub error(%v)", err)
				continue
			}
			mu.Lock()
			subs = append(subs, sub)
			mu.Unlock()
		}
	}()
	time.Sleep(time.Second * 5)
	for i := 0; i < 20; i++ {
		sub, err := NewSub(newConn(c, time.Second, time.Second), _testGroup, _testTopic, "", _testConfig, 100)
		if err != nil {
			t.Fatal(err)
			continue
		}
		mu.Lock()
		subs = append(subs, sub)
		mu.Unlock()
	}
	time.Sleep(time.Second * 5)
}
