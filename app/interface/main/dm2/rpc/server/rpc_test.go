package server

import (
	"flag"
	"net/rpc"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/dm2/conf"
	"go-common/app/interface/main/dm2/model"
	"go-common/app/interface/main/dm2/service"
	rpcx "go-common/library/net/rpc"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	client *rpc.Client
	_noArg = &struct{}{}
)

const (
	_addr          = "127.0.0.1:6709"
	_subjectInfos  = "RPC.SubjectInfos"
	_buyAdvance    = "RPC.BuyAdvance"
	_advanceState  = "RPC.AdvanceState"
	_advances      = "RPC.Advances"
	_passAdvance   = "RPC.PassAdvance"
	_denyAdvance   = "RPC.DenyAdvance"
	_cancelAdvance = "RPC.CancelAdvance"
	_mask          = "RPC.Mask"
)

func TestMain(m *testing.M) {
	var err error
	dir, _ := filepath.Abs("../../cmd/dm2-test.toml")
	if err = flag.Set("conf", dir); err != nil {
		panic(err)
	}
	if err = conf.Init(); err != nil {
		panic(err)
	}
	svr := service.New(conf.Conf)
	r := &RPC{s: svr}
	server := rpcx.NewServer(conf.Conf.RPCServer)
	if err = server.Register(r); err != nil {
		panic(err)
	}
	if client, err = rpc.Dial("tcp", _addr); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestSubjectInfos(t *testing.T) {
	var (
		tp   int32 = 1
		oids       = []int64{1221, 1231}
		res        = make(map[int64]*model.SubjectInfo)
	)
	Convey("get dm subject info", t, func() {
		arg := model.ArgOids{Type: tp, Oids: oids}
		err := client.Call(_subjectInfos, arg, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		for cid, r := range res {
			t.Logf("=====cid:%d Infos:%+v", cid, r)
		}
	})
}

func TestBuyAdvance(t *testing.T) {
	var (
		mid  int64 = 27515260
		cid  int64 = 10107292
		mode       = "sp"
	)
	Convey("buy advance dm", t, func() {
		arg := &model.ArgAdvance{
			Mid:  mid,
			Cid:  cid,
			Mode: mode,
		}
		err := client.Call(_buyAdvance, arg, _noArg)
		So(err, ShouldBeNil)

	})
}

func TestAdvanceState(t *testing.T) {
	var (
		mid  int64 = 27515330
		cid  int64 = 10107292
		mode       = "sp"
		res        = &model.AdvState{}
	)
	Convey("get advance dm state", t, func() {
		arg := &model.ArgAdvance{
			Mid:  mid,
			Cid:  cid,
			Mode: mode,
		}
		err := client.Call(_advanceState, arg, res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestAdvances(t *testing.T) {
	var (
		mid int64 = 27515260
		res       = make([]*model.Advance, 10)
	)
	Convey("get advances dm", t, func() {
		arg := &model.ArgMid{
			Mid: mid,
		}
		err := client.Call(_advances, arg, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestPassAdvance(t *testing.T) {
	var (
		mid int64 = 7158471
		id  int64 = 2
	)
	Convey("pass advance dm ", t, func() {
		arg := &model.ArgUpAdvance{
			Mid: mid,
			ID:  id,
		}
		err := client.Call(_passAdvance, arg, _noArg)
		So(err, ShouldBeNil)
	})
}

func TestDenyAdvance(t *testing.T) {
	var (
		mid int64 = 27515615
		id  int64 = 107
	)
	Convey("deny advance dm", t, func() {
		arg := &model.ArgUpAdvance{
			Mid: mid,
			ID:  id,
		}
		err := client.Call(_denyAdvance, arg, _noArg)
		So(err, ShouldBeNil)
	})
}

func TestCancelAdvance(t *testing.T) {
	var (
		mid int64 = 27515615
		id  int64 = 122
	)
	Convey("cancel advance dm", t, func() {
		arg := &model.ArgUpAdvance{
			Mid: mid,
			ID:  id,
		}
		err := client.Call(_cancelAdvance, arg, _noArg)
		So(err, ShouldBeNil)
	})
}

func TestMask(t *testing.T) {
	var (
		cid int64 = 32
		res       = &model.Mask{}
	)
	Convey("test mask list", t, func() {
		arg := &model.ArgMask{
			Cid: cid,
		}
		err := client.Call(_mask, arg, res)
		t.Logf("=========%+v", res)
		So(err, ShouldBeNil)
	})
}
