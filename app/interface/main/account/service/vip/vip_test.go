package vip

import (
	"context"
	"flag"
	"testing"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/model"

	vipmod "go-common/app/service/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	flag.Set("conf", "../../cmd/account-interface-example.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}

func TestService_CodeOpen(t *testing.T) {
	Convey("code open", t, func() {

		codeInfo, err := s.CodeOpen(context.TODO(), 123, "7b6e2263b8355928", "fd09f95433ed4c579f03ca7112b843ab", "45tn")
		t.Logf("%v", codeInfo)
		So(err, ShouldBeNil)
	})
}

func TestService_CodeVerify(t *testing.T) {
	Convey("code verify", t, func() {
		_, err := s.CodeVerify(context.TODO())
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceTips
func TestServiceTips(t *testing.T) {
	Convey("TestServiceTips", t, func() {
		req := &model.TipsReq{
			Version:  int64(6000),
			Platform: "ios",
		}
		res, err := s.Tips(context.TODO(), req)
		t.Logf("data(+%v)", res)
		So(err, ShouldBeNil)
	})
}

func TestService_CodeOpeneds(t *testing.T) {
	Convey("test service code opened", t, func() {
		arg := new(model.CodeInfoReq)
		resp, err := s.CodeOpeneds(context.TODO(), arg, "127.0.0.1")
		t.Logf("data(%+v)", resp)
		So(err, ShouldBeNil)
	})
}

func TestService_Unfrozen(t *testing.T) {
	Convey("test unfrozen", t, func() {
		err := s.Unfrozen(context.TODO(), 10001)
		So(err, ShouldBeNil)
	})
}

func TestService_FrozenTime(t *testing.T) {
	Convey("test frozen time", t, func() {
		ctime, err := s.FrozenTime(context.TODO(), 10001)
		t.Logf("%+v", ctime)
		So(err, ShouldBeNil)
	})
}

func TestService_checkIp(t *testing.T) {
	Convey("test check ip", t, func() {
		err := s.checkIP("b2cf4e9dbe9fd2e3", "111.203.12.97")
		So(err, ShouldBeNil)
		err = s.checkIP("b2cf4e9dbe9fd2e31", "111.203.12.97")
		So(err, ShouldNotBeNil)
	})
}

func TestService_OrderStatus(t *testing.T) {
	Convey("TestService_OrderStatus", t, func() {
		arg := &vipmod.ArgDialog{OrderNo: "1"}
		res, err := s.OrderStatus(context.Background(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
