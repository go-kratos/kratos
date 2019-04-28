package server

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/passport/conf"
	"go-common/app/service/main/passport/model"
	"go-common/app/service/main/passport/service"
	"go-common/library/log"
	"go-common/library/net/rpc"
)

const (
	_rpcLoginLogs = "RPC.LoginLogs"
)

func TestRPC_LoginLogs(t *testing.T) {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	s := service.New(conf.Conf)
	r := New(conf.Conf, s)
	defer r.Close()
	c2 := rpc.NewDiscoveryCli("passport.service", nil)
	time.Sleep(time.Second * 2)

	ms := make([]*model.LoginLog, 0)
	if err := c2.Call(context.TODO(), _rpcLoginLogs, &model.ArgLoginLogs{
		Mid: 88888970,
	}, &ms); err != nil {
		t.Errorf("failed to call %s, error(%v)", _rpcLoginLogs, err)
		t.FailNow()
	}
}
