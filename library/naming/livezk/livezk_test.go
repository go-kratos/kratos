package livezk

import (
	"context"
	"fmt"
	"path"
	"testing"
	"time"

	"go-common/library/naming"
	xtime "go-common/library/time"
)

var appdID = "main.arch.test6"
var addr = "127.0.0.1:8080"

var ins1 = &naming.Instance{
	AppID:   appdID,
	Addrs:   []string{"grpc://" + addr},
	Version: "1",
	Metadata: map[string]string{
		"test":  "1",
		"color": "blue",
	},
}

var addrs = []string{"172.18.33.131:2181", "172.18.33.168:2181", "172.18.33.169:2181"}
var zkConfig = &Zookeeper{
	Addrs:   addrs,
	Timeout: xtime.Duration(time.Second),
}

func TestLiveZK(t *testing.T) {
	reg, err := New(zkConfig)
	if err != nil {
		t.Fatal(err)
	}
	cancel, err := reg.Register(context.TODO(), ins1)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()
	lzk := reg.(*livezk)
	nodePath := path.Join(basePath, appdID, addr)
	ok, _, err := lzk.zkConn.Exists(nodePath)
	if err != nil {
		if err != nil {
			t.Fatal(err)
		}
	}
	if !ok {
		t.Errorf("path not exists %s", nodePath)
	}
}

func TestLiveZKCancel(t *testing.T) {
	reg, err := New(zkConfig)
	if err != nil {
		t.Fatal(err)
	}
	cancel, err := reg.Register(context.TODO(), ins1)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	lzk := reg.(*livezk)
	nodePath := path.Join(basePath, fmt.Sprintf("b%s", ins1.AppID), addr)
	ok, _, err := lzk.zkConn.Exists(nodePath)
	if err != nil {
		if err != nil {
			t.Fatal(err)
		}
	}
	if ok {
		t.Errorf("path should not exists %s", nodePath)
	}
}
