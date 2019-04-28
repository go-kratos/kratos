package client

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"go-common/app/service/main/passport/model"
)

var (
	once        sync.Once
	passportSvc *Client2
)

func startRPCServer() {
	passportSvc = New(nil)
	time.Sleep(time.Second * 2)
}

func TestService2_LoginLogs(t *testing.T) {
	once.Do(startRPCServer)
	arg := &model.ArgLoginLogs{
		Mid: 88888970,
	}
	if res, err := passportSvc.LoginLogs(context.TODO(), arg); err != nil {
		t.Errorf("failed to call rpc, passportSvc.LoginLogs(%v) error(%v)", arg, err)
		t.FailNow()
	} else if len(res) == 0 {
		t.Errorf("res is incorrect, expected res length > 0 but got 0")
		t.FailNow()
	} else {
		for i, v := range res {
			str, _ := json.Marshal(v)
			t.Logf("res[%d]: %s", i, str)
		}
	}
}
