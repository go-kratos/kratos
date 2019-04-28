package databus_test

import (
	"context"
	"testing"
	"time"

	"go-common/library/naming/discovery"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"
)

var (
	pCfg = &databus.Config{
		// Key:          "0PvKGhAqDvsK7zitmS8t",
		// Secret:       "0PvKGhAqDvsK7zitmS8u",
		// Group:        "databus_test_group",
		// Topic:        "databus_test_topic",
		Key:    "dbe67e6a4c36f877",
		Secret: "8c775ea242caa367ba5c876c04576571",
		Group:  "Test1-MainCommonArch-P",
		Topic:  "test1",
		Action: "pub",
		Name:   "databus",
		Proto:  "tcp",
		// Addr:         "172.16.33.158:6205",
		Addr:         "172.18.33.50:6205",
		Active:       10,
		Idle:         5,
		DialTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		IdleTimeout:  xtime.Duration(time.Minute),
	}
	sCfg = &databus.Config{
		// Key:          "0PvKGhAqDvsK7zitmS8t",
		// Secret:       "0PvKGhAqDvsK7zitmS8u",
		// Group:        "databus_test_group",
		// Topic:        "databus_test_topic",
		Key:    "dbe67e6a4c36f877",
		Secret: "8c775ea242caa367ba5c876c04576571",
		Group:  "Test1-MainCommonArch-S",
		Topic:  "test1",
		Action: "sub",
		Name:   "databus",
		Proto:  "tcp",
		// Addr:         "172.16.33.158:6205",
		Addr:         "172.18.33.50:6205",
		Active:       10,
		Idle:         5,
		DialTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second * 35),
		IdleTimeout:  xtime.Duration(time.Minute),
	}
	dCfg = &discovery.Config{
		Nodes:  []string{"172.18.33.50:7171"},
		Key:    "0c4b8fe3ff35a4b6",
		Secret: "b370880d1aca7d3a289b9b9a7f4d6812",
		Zone:   "sh001",
		Env:    "uat",
	}
)

type TestMsg struct {
	Now int64 `json:"now"`
}

func testSub(t *testing.T, d *databus.Databus) {
	for {
		m, ok := <-d.Messages()
		if !ok {
			return
		}
		t.Logf("sub message: %s", string(m.Value))
		if err := m.Commit(); err != nil {
			t.Errorf("sub commit error(%v)\n", err)
		}
	}
}

func testPub(t *testing.T, d *databus.Databus) {
	// pub
	m := &TestMsg{Now: time.Now().UnixNano()}
	if err := d.Send(context.TODO(), "test", m); err != nil {
		t.Errorf("d.Send(test) error(%v)", err)
	} else {
		t.Logf("pub message %v", m)
	}
}

func TestDatabus(t *testing.T) {
	d := databus.New(pCfg)
	// pub
	testPub(t, d)
	testPub(t, d)
	testPub(t, d)
	d.Close()

	// sub
	d = databus.New(sCfg)
	go testSub(t, d)

	time.Sleep(time.Second * 15)
	d.Close()
}

func TestDiscoveryDatabus(t *testing.T) {
	d := databus.New(pCfg)
	// pub
	testPub(t, d)
	testPub(t, d)
	testPub(t, d)
	d.Close()

	// sub
	d = databus.New(sCfg)
	go testSub(t, d)

	time.Sleep(time.Second * 15)
	d.Close()
}

func BenchmarkPub(b *testing.B) {
	d := databus.New(pCfg)
	defer d.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m := &TestMsg{Now: time.Now().UnixNano()}
			if err := d.Send(context.TODO(), "test", m); err != nil {
				b.Errorf("d.Send(test) error(%v)", err)
				continue
			}
		}
	})
}

func BenchmarkDiscoveryPub(b *testing.B) {
	d := databus.New(pCfg)
	defer d.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m := &TestMsg{Now: time.Now().UnixNano()}
			if err := d.Send(context.TODO(), "test", m); err != nil {
				b.Errorf("d.Send(test) error(%v)", err)
				continue
			}
		}
	})
}
