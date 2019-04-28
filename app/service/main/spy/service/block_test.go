package service

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/service/main/spy/model"
	spy "go-common/app/service/main/spy/model"
	"go-common/library/net/ip"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testLowMid    int64 = 59999
	testLowScore  int8  = 5
	testHighScore int8  = 90
	testLv5Mid    int64 = 15555180
	testNoLv5Mid  int64 = 200
)

func Test_Verfiy(t *testing.T) {
	Convey("Test_Verfiy block ", t, WithService(func(s *Service) {
		ui, err := s.UserInfo(c, testLowMid, "")
		fmt.Println(ui)
		So(ui, ShouldNotBeEmpty)
		So(err, ShouldBeNil)

		hs := []HandlerFunc{s.scoreLessHandler}
		Convey("Test_Verfiy scoreLessHandler low score ", WithService(func(s *Service) {
			ui.Score = testLowScore
			args := &Args{ui: ui}
			b := s.Verify(c, args, hs...)
			fmt.Println("b", b, ui.Score)
			So(b, ShouldBeTrue)
		}))
		Convey("Test_Verfiy scoreLessHandler high score ", WithService(func(s *Service) {
			ui.Score = testHighScore
			args := &Args{ui: ui}
			b := s.Verify(c, args, hs...)
			fmt.Println("b", b, ui.Score)
			So(b, ShouldBeFalse)
		}))
	}))
}

func Test_BlockFilter(t *testing.T) {
	Convey("Test_BlockFilter filter ", t, WithService(func(s *Service) {
		var (
			c = context.TODO()
		)
		ui, err := s.UserInfo(c, testLowMid, "")
		fmt.Println(ui)
		So(ui, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
		tx, err := s.dao.BeginTran(c)
		So(err, ShouldBeNil)
		Convey("Test_BlockFilter block ", WithService(func(s *Service) {
			ret, err := s.dao.BlockMidCache(c, s.BlockNo(s.c.Property.Block.CycleTimes), 10)
			So(err, ShouldBeNil)
			So(ret, ShouldBeNil)

			ui.Score = testLowScore
			state, err := s.BlockFilter(context.TODO(), ui)
			So(err, ShouldBeNil)
			So(state == model.StateBlock, ShouldBeTrue)
			err = tx.Commit()
			So(err, ShouldBeNil)
			Convey("Test_BlockFilter cache had mid ", WithService(func(s *Service) {
				ret, err := s.dao.BlockMidCache(c, s.BlockNo(s.c.Property.Block.CycleTimes), 10)
				fmt.Println("ret", ret)
				So(err, ShouldBeNil)
				So(ret, ShouldContain, ui.Mid)
			}))
		}))
		Convey("Test_BlockFilter bo block ", WithService(func(s *Service) {
			ui.Score = testHighScore
			state, err := s.BlockFilter(context.TODO(), ui)
			So(err, ShouldBeNil)
			So(state == model.StateBlock, ShouldBeFalse)
			err = tx.Commit()
			So(err, ShouldBeNil)
		}))

	}))
}

func Test_ReVerify(t *testing.T) {
	Convey("Test_ReVerify user lv check ", t, WithService(func(s *Service) {
		var (
			c = context.TODO()
		)
		Convey("test vip user ", WithService(func(s *Service) {
			err := s.ClearReliveTimes(c, &spy.ArgReset{Mid: testLv5Mid, Operator: "yubaihai"})
			So(err, ShouldBeNil)

			ui, err := s.UserInfo(c, testLv5Mid, ip.InternalIP())
			So(err, ShouldBeNil)
			So(ui, ShouldNotBeEmpty)
			fmt.Println("ui", ui)

			b, err := s.reVerifyHandler(c, ui)
			So(err, ShouldBeNil)
			So(b, ShouldBeFalse)

			// two
			b, err = s.reVerifyHandler(c, ui)
			So(err, ShouldBeNil)
			So(b, ShouldBeFalse)

			// three
			b, err = s.reVerifyHandler(c, ui)
			So(err, ShouldBeNil)
			So(b, ShouldBeTrue)
		}))
		Convey("test no lv vip user ", WithService(func(s *Service) {

			ui, err := s.UserInfo(c, testNoLv5Mid, ip.InternalIP())
			So(err, ShouldBeNil)
			So(ui, ShouldNotBeEmpty)
			fmt.Println("ui", ui)

			b, err := s.reVerifyHandler(c, ui)
			So(err, ShouldBeNil)
			So(b, ShouldBeTrue)
		}))
	}))
}
