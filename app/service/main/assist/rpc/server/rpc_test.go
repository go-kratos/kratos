package server

import (
	"go-common/app/service/main/assist/conf"
	"go-common/app/service/main/assist/model/assist"
	model "go-common/app/service/main/assist/model/assist"
	"go-common/app/service/main/assist/service"
	"net/rpc"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const (
	addr             = "127.0.0.1:6729"
	_assists         = "RPC.Assists"
	_assist          = "RPC.Assist"
	_addAssist       = "RPC.AddAssist"
	_delAssist       = "RPC.DelAssist"
	_assistExit      = "RPC.AssistExit"
	_assistLogInfo   = "RPC.AssistLogInfo"
	_assistLogs      = "RPC.AssistLogs"
	_assistLogAdd    = "RPC.AssistLogAdd"
	_assistLogCancel = "RPC.AssistLogCancel"
	_assistUps       = "RPC.AssistUps"
)

func TestAssistExit(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	svr := service.New(conf.Conf)
	New(conf.Conf, svr)
	client, err := rpc.Dial("tcp", addr)
	defer client.Close()
	if err != nil {
		t.Errorf("rpc.Dial(tcp, \"%s\") error(%v)", addr, err)
		t.FailNow()
	}
	assistExit(client, t)
}

func assistExit(client *rpc.Client, t *testing.T) {
	assistInfo := new(struct{})
	arg := &model.ArgAssist{
		Mid:       254386,
		AssistMid: 2,
		RealIP:    "127.0.0.1",
	}
	err := client.Call(_assistExit, arg, assistInfo)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	spew.Dump(assistInfo)
}

func TestAssistUps(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	svr := service.New(conf.Conf)
	New(conf.Conf, svr)
	client, err := rpc.Dial("tcp", addr)
	defer client.Close()
	if err != nil {
		t.Errorf("rpc.Dial(tcp, \"%s\") error(%v)", addr, err)
		t.FailNow()
	}
	assistUps(client, t)
}

func assistUps(client *rpc.Client, t *testing.T) {
	var res = &assist.AssistUpsPager{}
	arg := &model.ArgAssistUps{
		AssistMid: 88889017,
		Pn:        1,
		Ps:        20,
		RealIP:    "127.0.0.1",
	}
	err := client.Call(_assistUps, arg, &res)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	t.Logf("%+v", res)
	spew.Dump(res)
}

func TestAssists(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	svr := service.New(conf.Conf)
	New(conf.Conf, svr)
	client, err := rpc.Dial("tcp", addr)
	defer client.Close()
	if err != nil {
		t.Errorf("rpc.Dial(tcp, \"%s\") error(%v)", addr, err)
		t.FailNow()
	}
	assists(client, t)
}

func TestAssistInfo(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	svr := service.New(conf.Conf)
	New(conf.Conf, svr)
	client, err := rpc.Dial("tcp", addr)
	defer client.Close()
	if err != nil {
		t.Errorf("rpc.Dial(tcp, \"%s\") error(%v)", addr, err)
		t.FailNow()
	}
	assistInfo(client, t)
}

func TestAddAssist(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	svr := service.New(conf.Conf)
	New(conf.Conf, svr)
	client, err := rpc.Dial("tcp", addr)
	defer client.Close()
	if err != nil {
		t.Errorf("rpc.Dial(tcp, \"%s\") error(%v)", addr, err)
		t.FailNow()
	}
	addAssist(client, t)
}

func TestDelAssist(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	svr := service.New(conf.Conf)
	New(conf.Conf, svr)
	client, err := rpc.Dial("tcp", addr)
	defer client.Close()
	if err != nil {
		t.Errorf("rpc.Dial(tcp, \"%s\") error(%v)", addr, err)
		t.FailNow()
	}
	delAssist(client, t)
}

func TestAssistLogInfo(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	svr := service.New(conf.Conf)
	New(conf.Conf, svr)
	client, err := rpc.Dial("tcp", addr)
	defer client.Close()
	if err != nil {
		t.Errorf("rpc.Dial(tcp, \"%s\") error(%v)", addr, err)
		t.FailNow()
	}
	assistLogInfo(client, t)
}

func TestAssistLogAdd(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	svr := service.New(conf.Conf)
	New(conf.Conf, svr)
	client, err := rpc.Dial("tcp", addr)
	defer client.Close()
	if err != nil {
		t.Errorf("rpc.Dial(tcp, \"%s\") error(%v)", addr, err)
		t.FailNow()
	}
	assistLogAdd(client, t)
}

func TestAssistLogCancel(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	svr := service.New(conf.Conf)
	New(conf.Conf, svr)
	client, err := rpc.Dial("tcp", addr)
	defer client.Close()
	if err != nil {
		t.Errorf("rpc.Dial(tcp, \"%s\") error(%v)", addr, err)
		t.FailNow()
	}
	assistLogCancel(client, t)
}

func TestAssistLogs(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	svr := service.New(conf.Conf)
	New(conf.Conf, svr)
	client, err := rpc.Dial("tcp", addr)
	defer client.Close()
	if err != nil {
		t.Errorf("rpc.Dial(tcp, \"%s\") error(%v)", addr, err)
		t.FailNow()
	}
	assistLogs(client, t)
}

func assistLogCancel(client *rpc.Client, t *testing.T) {
	res := new(struct{})
	arg := &model.ArgAssistLog{
		Mid:       254386,
		AssistMid: 2089809,
		LogID:     670,
		RealIP:    "127.0.0.1",
	}
	err := client.Call(_assistLogCancel, arg, &res)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	t.Logf("%+v", res)
}

func assistLogAdd(client *rpc.Client, t *testing.T) {
	res := new(struct{})
	arg := &model.ArgAssistLogAdd{
		Mid:       254386,
		AssistMid: 2089809,
		Type:      model.TypeComment,
		Action:    model.ActDelete,
		SubjectID: 111,
		ObjectID:  "444",
		Detail:    "testing",
		RealIP:    "127.0.0.1",
	}
	err := client.Call(_assistLogAdd, arg, &res)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	t.Logf("%+v", res)
}

func assistInfo(client *rpc.Client, t *testing.T) {
	var res = &assist.AssistRes{}
	arg := &model.ArgAssist{
		Mid:       27515256,
		AssistMid: 27515255,
		RealIP:    "127.0.0.1",
	}
	err := client.Call(_assist, arg, &res)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	t.Logf("%+v", res)
}

func addAssist(client *rpc.Client, t *testing.T) {
	assistInfo := new(struct{})
	arg := &model.ArgAssist{
		Mid:       27515256,
		AssistMid: 27515255,
		RealIP:    "127.0.0.1",
	}
	err := client.Call(_addAssist, arg, assistInfo)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	spew.Dump(assistInfo)
}

func delAssist(client *rpc.Client, t *testing.T) {
	assistInfo := new(struct{})
	arg := &model.ArgAssist{
		Mid:       27515256,
		AssistMid: 27515255,
		RealIP:    "127.0.0.1",
	}
	err := client.Call(_delAssist, arg, assistInfo)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	spew.Dump(assistInfo)
}

func assists(client *rpc.Client, t *testing.T) {
	var res = make([]*assist.Assist, 0)
	arg := &model.ArgAssists{
		Mid:    254386,
		RealIP: "127.0.0.1",
	}
	err := client.Call(_assists, arg, &res)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	t.Logf("%+v", res)
	spew.Dump(res)
}

func assistLogs(client *rpc.Client, t *testing.T) {
	var res = make([]*assist.Log, 0)
	arg := &model.ArgAssistLogs{
		Mid:       254386,
		AssistMid: 2089809,
		Stime:     time.Unix(1496205563, 0),
		Etime:     time.Unix(1496291963, 0),
		Pn:        1,
		Ps:        30,
		RealIP:    "127.0.0.1",
	}
	err := client.Call(_assistLogs, arg, &res)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	t.Logf("%+v", res)
	spew.Dump(res)
}

func assistLogInfo(client *rpc.Client, t *testing.T) {
	var res = &assist.Log{}
	arg := &model.ArgAssistLog{
		Mid:       254386,
		AssistMid: 2089809,
		LogID:     15,
		RealIP:    "127.0.0.1",
	}
	err := client.Call(_assistLogInfo, arg, &res)
	if err != nil {
		t.Logf("err:%v.", err)
	}
	t.Logf("%+v", res)
	spew.Dump(res)
}
