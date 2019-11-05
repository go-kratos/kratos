package zookeeper

import (
	"context"
	"testing"

	"github.com/bilibili/kratos/pkg/naming"
)

var (
	_testAppid = "test_appid"

	_testConf = &Config{
		Root:      "/test",
		Endpoints: []string{"127.0.0.1:2181"},
	}
	_testIns = &naming.Instance{
		AppID: _testAppid,
		Addrs: []string{"grpc://127.0.0.1:9000"},
		Metadata: map[string]string{
			"test_key": "test_value",
		},
	}
)

func TestZookeeper(t *testing.T) {
	zk, err := New(_testConf)
	if err != nil {
		t.Fatal(err)
	}
	_, err = zk.Register(context.TODO(), _testIns)
	if err != nil {
		t.Fatal(err)
	}
	// fetch&watch
	res := zk.Build(_testAppid)
	event := res.Watch()
	<-event
	in, ok := res.Fetch(context.TODO())
	if !ok {
		t.Fatal("failed to resolver fetch")
	}
	if len(in.Instances) != 1 {
		t.Fatalf("Instances not match, got:%d want:1", len(in.Instances))
	}
}
