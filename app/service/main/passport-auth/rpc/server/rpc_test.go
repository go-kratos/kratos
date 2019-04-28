package server

import (
	"net/rpc"
	"sync"
	"testing"
	"time"

	"go-common/app/service/main/passport-auth/conf"
	"go-common/app/service/main/passport-auth/model"
	"go-common/app/service/main/passport-auth/service"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	addr      = "127.0.0.1:7579"
	_testPing = "RPC.Ping"
	// token
	_tokenInfo = "RPC.TokenInfo"
	_delToken  = "RPC.DelToken"

	// cookie
	_cookieInfo = "RPC.CookieInfo"
	_delCookie  = "RPC.DelCookie"
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

func TestRPC_TokenInfo(t *testing.T) {
	var (
		err error
		arg = "64294c76972aee8cf4af51566c76ed0d"
		res *model.Token
	)
	once.Do(startServer)
	Convey("Test RPC server get token by token", t, func() {
		err = client.Call(_tokenInfo, arg, res)
		So(err, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestRPC_DelToken(t *testing.T) {
	var (
		err error
		arg = "64294c76972aee8cf4af51566c76ed0d"
		res int64
	)
	once.Do(startServer)
	time.Sleep(3 * time.Second)
	Convey("Test RPC server del token by token", t, func() {
		err = client.Call(_delToken, arg, res)
		So(err, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestRPC_CookieInfo(t *testing.T) {
	var (
		err error
		arg = "c1300f65,1519273116,f05bd5ef"
		res *model.Cookie
	)
	once.Do(startServer)
	Convey("Test RPC server get coookie by ssda", t, func() {
		err = client.Call(_cookieInfo, arg, res)
		So(err, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestRPC_DelCookie(t *testing.T) {
	var (
		err error
		arg = "c1300f65,1519273116,f05bd5ef"
		res int64
	)
	once.Do(startServer)
	time.Sleep(3 * time.Second)
	Convey("Test RPC server del coookie by cookie", t, func() {
		err = client.Call(_delCookie, arg, res)
		So(err, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
