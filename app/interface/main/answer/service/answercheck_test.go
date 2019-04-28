package service

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceProCheck(t *testing.T) {
	convey.Convey("ProCheck", t, func() {
		hid, err := s.ProCheck(context.Background(), 14771787, []int64{}, nil, "")
		convey.So(err, convey.ShouldBeNil)
		convey.So(hid, convey.ShouldNotBeNil)
	})
}

func TestServiceCheckBase(t *testing.T) {
	convey.Convey("CheckBase", t, func() {
		res, err := s.CheckBase(context.Background(), 14771787, []int64{}, nil, "")
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServiceCaptcha(t *testing.T) {
	convey.Convey("Captcha", t, func() {
		res, err := s.Captcha(context.Background(), 14771787, "", 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServiceValidate(t *testing.T) {
	convey.Convey("Validate", t, func() {
		res, err := s.Validate(context.Background(), "", "", "", "", 0, 0, "", "", nil)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServicecheckQsIDs(t *testing.T) {
	convey.Convey("checkQsIDs", t, func() {
		ok, err := s.checkQsIDs(context.Background(), []int64{}, 0, []int64{}, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ok, convey.ShouldNotBeNil)
	})
}

func TestServicecheckAns(t *testing.T) {
	convey.Convey("checkAns", t, func() {
		errIds, rc, err := s.checkAns(context.Background(), 14771787, []int64{}, nil, "", 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rc, convey.ShouldNotBeNil)
		convey.So(errIds, convey.ShouldNotBeNil)
	})
}

func TestServicebasePass(t *testing.T) {
	convey.Convey("basePass", t, func() {
		s.basePass(context.Background(), 0, nil, time.Now())
	})
}

func TestServicependant(t *testing.T) {
	convey.Convey("pendant", t, func() {
		err := s.pendant(context.Background(), nil, 0, "", nil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestServicecheckAnswerBlock(t *testing.T) {
	convey.Convey("checkAnswerBlock", t, func() {
		block := s.checkAnswerBlock(context.Background(), 0)
		convey.So(block, convey.ShouldNotBeNil)
	})
}

func TestServicesendData(t *testing.T) {
	convey.Convey("sendData", t, func() {
		s.sendData(context.Background(), nil, nil, "")
	})
}

func TestServiceExtraCheck(t *testing.T) {
	convey.Convey("ExtraCheck", t, func() {
		err := s.ExtraCheck(context.Background(), 14771787, []int64{}, nil, "", "", "", "")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestServicecheckExtraPassAns(t *testing.T) {
	convey.Convey("checkExtraPassAns", t, func() {
		ret, qs, err := s.checkExtraPassAns(context.Background(), 14771787, []int64{}, nil, "", 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(qs, convey.ShouldNotBeNil)
		convey.So(ret, convey.ShouldNotBeNil)
	})
}

func TestServicesendExtraRetMsg(t *testing.T) {
	convey.Convey("sendExtraRetMsg", t, func() {
		rs, err := s.sendExtraRetMsg(context.Background(), 0, nil, []int64{}, nil, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rs, convey.ShouldNotBeNil)
	})
}

func TestServicePendantRec(t *testing.T) {
	convey.Convey("PendantRec", t, func() {
		err := s.PendantRec(context.Background(), nil)
		convey.So(err, convey.ShouldBeNil)
	})
}
