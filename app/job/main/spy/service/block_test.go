package service

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/job/main/spy/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testBlockMid int64 = 4780461
	testLowScore int8  = 7
)

func Test_BlockReason(t *testing.T) {
	Convey("Test_BlockReason get block reason", t, WithService(func(s *Service) {
		reason, remake := s.blockReason(context.TODO(), testBlockMid)
		fmt.Println("reason remake", reason, remake)
		So(reason, ShouldNotBeEmpty)
		So(remake, ShouldNotBeEmpty)
	}))
}

func Test_CanBlock(t *testing.T) {
	Convey("Test_CanBlock can block ", t, WithService(func(s *Service) {
		tx, err := s.dao.BeginTran(c)
		So(err, ShouldBeNil)
		ui := &model.UserInfo{Mid: testBlockMid, State: model.StateNormal}
		err = s.dao.TxUpdateUserState(c, tx, ui)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)

		ui, ok := s.canBlock(context.TODO(), testBlockMid)
		fmt.Println("Test_CanBlock ui ", ui, testBlockMid)
		So(ui, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
	}))
}

func Test_Block(t *testing.T) {
	Convey("Test_Block block ", t, WithService(func(s *Service) {
		ui, err := s.dao.UserInfo(context.TODO(), testBlockMid)
		So(err, ShouldBeNil)

		tx, err := s.dao.BeginTran(context.TODO())
		So(err, ShouldBeNil)
		ui.State = model.StateNormal
		err = s.dao.TxUpdateUserState(c, tx, ui)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)

		ui, err = s.dao.UserInfo(context.TODO(), testBlockMid)
		So(err, ShouldBeNil)
		So(ui.State == model.StateNormal, ShouldBeTrue)

		reason, remake := s.blockReason(context.TODO(), testBlockMid)
		fmt.Println("reason remake", reason, remake)
		So(reason, ShouldNotBeEmpty)
		So(remake, ShouldNotBeEmpty)

		Convey("Test_CanBlock do block ", WithService(func(s *Service) {
			err := s.blockByMid(context.TODO(), testBlockMid)
			So(err, ShouldBeNil)
			Convey("Test_CanBlock get block user info ", WithService(func(s *Service) {
				ui, err := s.dao.UserInfo(context.TODO(), testBlockMid)
				So(err, ShouldBeNil)
				So(ui.State == model.StateBlock, ShouldBeTrue)
			}))
		}))
		fmt.Println("Test_Block end ")

	}))
}
