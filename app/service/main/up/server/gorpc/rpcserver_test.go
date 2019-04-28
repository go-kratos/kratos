package gorpc

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/service/main/up/model"
	"go-common/library/net/rpc"
	xtime "go-common/library/time"

	"github.com/davecgh/go-spew/spew"
)

const (
	/*clientConfigStr = `
		 proto = "tcp"
	     timeout = "1s"
	     timer = 1000
	     token = "123456"
	     addr = "127.0.0.1:6079"
	     [breaker]
	     window  = "3s"
	     sleep   = "100ms"
	     bucket  = 10
	     ratio   = 0.5
	     request = 100`*/

	_Special     = "RPC.Special"
	_Info        = "RPC.Info"
	_SetUpSwitch = "RPC.SetUpSwitch"
	_UpSwitch    = "RPC.UpSwitch"
	_UpCards     = "RPC.UpCards"
)

func init() {
	dir, _ := filepath.Abs("../../cmd/up-service.toml")
	flag.Set("conf", dir)
}

func initSvrAndClient(t *testing.T) (client *rpc.Client, err error) {
	client = rpc.Dial("127.0.0.1:6079", xtime.Duration(time.Second), nil)
	return
}

func TestInfo(t *testing.T) {
	client, err := initSvrAndClient(t)
	if err != nil {
		t.Errorf("rpc.Dial error(%v)", err)
		t.FailNow()
	}
	defer client.Close()
	//time.Sleep(1 * time.Second)
	info(client, t)
}

func info(client *rpc.Client, t *testing.T) {
	var res *model.UpInfo
	arg := &model.ArgInfo{
		Mid:  2089809,
		From: 1,
	}
	err := client.Call(context.TODO(), _Info, arg, &res)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	spew.Dump(res)
}

func TestSpecial(t *testing.T) {
	client, err := initSvrAndClient(t)
	if err != nil {
		t.Errorf("rpc.Dial error(%v)", err)
		t.FailNow()
	}
	defer client.Close()
	//time.Sleep(1 * time.Second)
	special(client, t)
}

func special(client *rpc.Client, t *testing.T) {
	var res []model.UpSpecial
	arg := &model.ArgSpecial{
		GroupID: 2,
	}
	err := client.Call(context.TODO(), _Special, arg, &res)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	spew.Dump(res)
}

func Test_UpSwitch(t *testing.T) {
	client, err := initSvrAndClient(t)
	if err != nil {
		t.Errorf("rpc.Dial error(%v)", err)
		t.FailNow()
	}
	defer client.Close()
	var res *model.PBUpSwitch
	arg := &model.ArgUpSwitch{
		Mid:  1,
		From: 0,
	}
	err = client.Call(context.TODO(), _UpSwitch, arg, &res)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	spew.Dump(11111, res)
}

func Test_SetUpSwitch(t *testing.T) {
	client, err := initSvrAndClient(t)
	if err != nil {
		t.Errorf("rpc.Dial error(%v)", err)
		t.FailNow()
	}
	defer client.Close()
	var res *model.PBSetUpSwitchRes
	arg := &model.ArgUpSwitch{
		Mid:   1,
		From:  0,
		State: 1,
	}
	err = client.Call(context.TODO(), _SetUpSwitch, arg, &res)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	spew.Dump(11111, res)
}

func Test_UpCards(t *testing.T) {
	client, err := initSvrAndClient(t)
	if err != nil {
		t.Errorf("rpc.Dial error(%v)", err)
		t.FailNow()
	}
	defer client.Close()

	arg := &model.ListUpCardInfoArg{
		Pn: 1,
		Ps: 1,
	}
	var res *model.UpCardInfoPage
	err = client.Call(context.TODO(), _UpCards, arg, &res)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	spew.Dump(11111, res)
}
