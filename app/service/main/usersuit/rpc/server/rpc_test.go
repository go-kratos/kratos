package server

import (
	"net/rpc"
	"sync"
	"testing"
	"time"

	"go-common/app/service/main/usersuit/conf"
	"go-common/app/service/main/usersuit/service"
	"go-common/library/log"
)

const (
	addr      = "127.0.0.1:7269"
	_testPing = "RPC.Ping"
)

var (
	_noArg = &struct{}{}
	client *rpc.Client
	once   sync.Once
)

func startServer() {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	svr := service.New(conf.Conf)
	New(conf.Conf, svr)
	time.Sleep(time.Second * 3)
	var err error
	client, err = rpc.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
}

func TestRPC_Ping(t *testing.T) {
	once.Do(startServer)
	if err := client.Call(_testPing, &_noArg, &_noArg); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _testPing, err)
		t.FailNow()
	}
}
