package service

import (
	"encoding/json"
	"testing"

	"go-common/app/service/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestServiceByMid
func TestServiceByMid(t *testing.T) {
	Convey("ByMid err == nil", t, func() {
		res, err := s.ByMid(c, 0)
		t.Logf("ByMid(%+v)", res)
		So(err, ShouldBeNil)
	})
}

func TestServiceVipInfos(t *testing.T) {
	Convey("ByMid err == nil", t, func() {
		mids := append(make([]int64, 1), int64(2089810))
		res, err := s.VipInfos(c, mids)
		t.Logf("VipInfos(%v)", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceVipInfo

func TestServiceVipInfo(t *testing.T) {
	Convey("ByMid err == nil", t, func() {
		res, err := s.VipInfo(c, int64(2))
		t.Logf("VipInfo(%v)", res)
		So(err, ShouldBeNil)
	})
}

func TestServiceH5History(t *testing.T) {
	Convey("ByMid err == nil", t, func() {
		arg := &model.ArgChangeHistory{Mid: 1001}
		res, err := s.H5History(c, arg)
		t.Logf("VipInfo(%v)", res)
		So(err, ShouldBeNil)
	})
}

func TestService_History(t *testing.T) {
	Convey("history err == nil", t, func() {
		arg := new(model.ArgChangeHistory)
		arg.Mid = 88889017
		arg.Pn = 1
		arg.Ps = 20
		vh, count, err := s.History(c, arg)
		bytes, _ := json.Marshal(vh)
		t.Logf("history(%+v) count(%v)", string(bytes), count)
		So(err, ShouldBeNil)
	})
}

func TestService_H5History(t *testing.T) {
	Convey("h5 history err == nil", t, func() {
		arg := new(model.ArgChangeHistory)
		arg.Mid = 1001
		vh, err := s.H5History(c, arg)
		t.Logf("h5history(%+v)", vh)
		So(err, ShouldBeNil)
	})
}

func TestService_OrderMng(t *testing.T) {
	Convey("order mng err == nil", t, func() {
		order, err := s.OrderMng(c, 10)
		t.Logf("%+v", order)
		So(err, ShouldBeNil)
	})
}

func TestService_Rescision(t *testing.T) {
	Convey("rescision err == nil", t, func() {
		err := s.Rescision(c, 10, 3)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceVipInfoBo
func TestServiceVipInfoBo(t *testing.T) {
	Convey("VipInfoBo err == nil", t, func() {
		res, err := s.VipInfoBo(c, int64(2089809))
		t.Logf("VipInfoBo(%v)", res)
		t.Logf("ios_overdue_time(%v)", res.IosOverdueTime)
		So(err, ShouldBeNil)
	})
}

func TestServiceSurplusFrozenTime(t *testing.T) {
	Convey("surplus frozen time", t, func() {
		stime, err := s.SurplusFrozenTime(c, 7593623)
		t.Logf("time:%v", stime)
		So(err, ShouldBeNil)
	})
}

func TestServiceUnfrozen(t *testing.T) {
	Convey("un frozen ", t, func() {
		err := s.Unfrozen(c, 1001)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceVipInfoAutoRenew
func TestServiceVipInfoAutoRenew(t *testing.T) {
	Convey("un frozen ", t, func() {
		var (
			v   *model.VipInfoResp
			err error
		)
		v, err = s.ByMid(c, iapAutoRenewOverdueMid)
		t.Logf("%+v", v)
		So(v, ShouldNotBeNil)
		So(v.PayType == 0, ShouldBeTrue)
		So(err, ShouldBeNil)

		v, err = s.ByMid(c, wechatAutoRenewOverdueMid)
		t.Logf("%+v", v)
		So(v, ShouldNotBeNil)
		So(v.PayType == 1, ShouldBeTrue)
		So(err, ShouldBeNil)

		v, err = s.ByMid(c, iapAutoRenewTodayOverdueMid)
		t.Logf("%+v", v)
		So(v, ShouldNotBeNil)
		So(v.PayType == 1, ShouldBeTrue)
		So(err, ShouldBeNil)

		v, err = s.ByMid(c, iapAutoRenewNotOverdueMid)
		t.Logf("%+v", v)
		So(v, ShouldNotBeNil)
		So(v.PayType == 1, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

var (
	iapAutoRenewOverdueMid       int64 = 4004193
	wechatAutoRenewOverdueMid    int64 = 27515398
	wechatAutoRenewNotOverdueMid int64 = 27515586
	iapAutoRenewTodayOverdueMid  int64 = 1
	iapAutoRenewNotOverdueMid    int64 = 27515406
)
