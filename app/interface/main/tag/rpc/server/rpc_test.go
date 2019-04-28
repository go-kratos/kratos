package rpc

import (
	"flag"
	"net/rpc"
	"testing"

	"go-common/app/interface/main/tag/conf"
	"go-common/app/interface/main/tag/model"
	"go-common/app/interface/main/tag/service"
	"go-common/library/log"
)

const (
	_testAddr = "127.0.0.1:6099"

	// test rpc
	_testPing        = "RPC.Ping"
	_testInfoByID    = "RPC.InfoByID"
	_testInfoByIDs   = "RPC.InfoByIDs"
	_testInfoByName  = "RPC.InfoByName"
	_testInfoByNames = "RPC.InfoByNames"
	_testArcTags     = "RPC.ArcTags"
	_testSubTags     = "RPC.SubTags"
	// test rpc2
	_testUpBind    = "RPC.UpBind"
	_testAdminBind = "RPC.AdminBind"
	_testResTags   = "RPC.ResTags"
)

func TestRPC(t *testing.T) {
	var (
		err error
	)
	flag.Set("conf", "../tag-example.toml")
	flag.Parse()
	if err = conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	// service
	svr := service.New(conf.Conf)
	// rpc
	Init(conf.Conf, svr)
	client, err := rpc.Dial("tcp", _testAddr)
	if err != nil {
		t.Errorf("rpc.Dial(tcp:%s) error(%v)", _testAddr, err)
		t.FailNow()
	}
	defer client.Close()
	// test rpc
	testPing(client, t)
	testInfoByID(client, t)
	testInfoByIDs(client, t)
	testInfoByName(client, t)
	testInfoByNames(client, t)
	// test rpc2
	testUpBind(client, t)
	testAdminBind(client, t)
	testResTags(client, t)
	testArcTags(client, t)
	testSubTags(client, t)
}

// test rpc

func testPing(client *rpc.Client, t *testing.T) {
	var (
		err error
		arg = &struct{}{}
		res = &struct{}{}
	)
	if err = client.Call(_testPing, arg, res); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _testPing, err)
	}
}

func testInfoByID(client *rpc.Client, t *testing.T) {
	var (
		err error
		arg = &model.ArgID{ID: 1, Mid: 27515274}
		res = &model.Tag{}
	)
	if err = client.Call(_testInfoByID, arg, res); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _testInfoByID, err)
	}
	t.Logf("InfoByID arg:%+v res:%+v", arg, res)
}

func testInfoByIDs(client *rpc.Client, t *testing.T) {
	var (
		err error
		arg = &model.ArgIDs{IDs: []int64{1, 2}, Mid: 27515274}
		res = &[]model.Tag{}
	)
	if err = client.Call(_testInfoByIDs, arg, res); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _testInfoByIDs, err)
	}
	t.Logf("InfoByIDs arg:%+v res:%+v", arg, res)
}

func testInfoByName(client *rpc.Client, t *testing.T) {
	var (
		err error
		arg = &model.ArgName{Name: "朱杰测试8", Mid: 27515274}
		res = &model.Tag{}
	)
	if err = client.Call(_testInfoByName, arg, res); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _testInfoByName, err)
	}
	t.Logf("InfoByName arg:%+v res:%+v", arg, res)
}

func testInfoByNames(client *rpc.Client, t *testing.T) {
	var (
		err error
		arg = &model.ArgNames{Names: []string{"朱杰测试8", "2012"}, Mid: 27515274}
		res = &[]model.Tag{}
	)
	if err = client.Call(_testInfoByNames, arg, res); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _testInfoByNames, err)
	}
	t.Logf("InfoByNames arg:%+v res:%+v", arg, res)
}

func testArcTags(client *rpc.Client, t *testing.T) {
	var (
		err error
		arg = &model.ArgAid{Aid: 1, Mid: 27515274}
		res = &[]model.Tag{}
	)
	if err = client.Call(_testArcTags, arg, res); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _testArcTags, err)
	}
	t.Logf("ArcTags arg:%+v res:%+v", arg, res)
}

func testSubTags(client *rpc.Client, t *testing.T) {
	var (
		err error
		arg = &model.ArgSub{Mid: 15555180, Pn: 1, Ps: 20, Order: -1}
		res = &[]model.Tag{}
	)
	if err = client.Call(_testSubTags, arg, res); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _testArcTags, err)
	}
	t.Logf("SubTags arg:%+v res:%+v", arg, res)
}

// test rpc2

func testUpBind(client *rpc.Client, t *testing.T) {
	var (
		err error
		arg = &model.ArgBind{Oid: 1, Mid: 123, Type: model.PicResType, Names: []string{"platform1", "platform2"}}
		res = &struct{}{}
	)
	if err = client.Call(_testUpBind, arg, res); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _testUpBind, err)
	}
	t.Logf("UpBind arg:%+v res:%+v", arg, res)
}

func testAdminBind(client *rpc.Client, t *testing.T) {
	var (
		err error
		arg = &model.ArgBind{Oid: 1, Mid: 1234, Type: model.PicResType, Names: []string{"platform1", "platform2", "platform3"}}
		res = &model.ArgBind{}
	)
	if err = client.Call(_testAdminBind, arg, res); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _testAdminBind, err)
	}
	t.Logf("AdminBind arg:%+v res:%+v", arg, res)
}

func testResTags(client *rpc.Client, t *testing.T) {
	var (
		err error
		arg = &model.ArgResTags{Oids: []int64{1}, Type: model.PicResType, Mid: 123}
		res map[int64][]*model.Tag
	)
	if err = client.Call(_testResTags, arg, &res); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _testResTags, err)
	}
	t.Logf("ResTags arg:%+v res:%+v", arg, res)
}
