package service

import (
	"testing"

	pb "go-common/app/service/main/coin/api"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUserCoins(t *testing.T) {
	var mid int64 = 2
	Convey("coin", t, func() {
		arg1 := &pb.UserCoinsReq{
			Mid: mid,
		}
		reply, err := s.UserCoins(ctx, arg1)
		So(err, ShouldBeNil)
		arg := &pb.ModifyCoinsReq{
			Mid:       mid,
			Count:     2,
			Reason:    "test",
			IP:        "",
			Operator:  "",
			CheckZero: 0,
			Ts:        0,
		}
		s.ModifyCoins(ctx, arg)

		ec, err := s.UserCoins(ctx, arg1)
		So(err, ShouldBeNil)
		So(ec.Count, ShouldEqual, reply.Count+2)
	})
}

func TestLog(t *testing.T) {
	Convey("log", t, func() {
		arg := &pb.CoinsLogReq{
			Mid:       88888929,
			Recent:    false,
			Translate: true,
		}
		ls, err := s.CoinsLog(ctx, arg)
		So(err, ShouldBeNil)
		So(ls.List, ShouldNotBeEmpty)
	})
}

func TestTranslateLog(t *testing.T) {
	Convey("cv", t, func() {
		l := "cv Rating for 123 : 3565 from 1234"
		exp := "专栏 cv123 收到打赏"
		So(translateLog(l), ShouldEqual, exp)
		l = "cv Rating for 123"
		exp = "给专栏 cv123 打赏"
		So(translateLog(l), ShouldEqual, exp)
		l = "2015萌战活动"
		exp = "投票资格"
		So(translateLog(l), ShouldEqual, exp)
	})
}

func TestRound(t *testing.T) {
	Convey("work", t, func() {
		So(Round(0.7), ShouldEqual, 0.7)
		So(Round(1.16), ShouldEqual, 1.16)
		So(Round(33.3), ShouldEqual, 33.3)
		So(Round(0.2900000000000001), ShouldEqual, 0.29)
	})
}

func TestCheckBusiness(t *testing.T) {
	Convey("present", t, func() {
		tp, err := s.CheckBusiness("article")
		So(tp, ShouldEqual, 1)
		So(err, ShouldBeNil)
	})

	Convey("not present", t, func() {
		tp, err := s.CheckBusiness("xxx")
		So(tp, ShouldEqual, 0)
		So(err, ShouldNotBeNil)
	})
}
