package service

import (
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/aegis/model/monitor"
	accApi "go-common/app/service/main/account/api"
	"testing"
)

func WithMock(t *testing.T, f func(mock *gomock.Controller)) func() {
	return func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		f(mockCtrl)
	}
}

func TestService_monitorArchive(t *testing.T) {
	var (
		na = &monitor.BinlogArchive{
			ID:     10111555,
			State:  -100,
			Round:  10,
			MID:    666,
			TypeID: 2422,
		}
	)
	Convey("monitorUpDelArc", t, func(ctx C) {
		errs := s.monitorArchive("update", nil, na)
		So(errs, ShouldNotBeEmpty)
	})
}

func TestService_monitorUpDelArc(t *testing.T) {
	var (
		na = &monitor.BinlogArchive{
			ID:     10111555,
			State:  -100,
			Round:  10,
			MID:    666,
			TypeID: 24,
		}
		logs []string
	)
	Convey("monitorUpDelArc", t, func(ctx C) {
		_, logs, _ = s.monitorUpDelArc(1, na)
		So(logs, ShouldNotBeEmpty)
	})
}

func TestService_monitorVideo(t *testing.T) {
	var (
		na = &monitor.BinlogVideo{
			ID:     10134809,
			Status: 0,
		}
	)
	Convey("monitorVideo", t, func(ctx C) {
		errs := s.monitorVideo("update", nil, na)
		So(errs, ShouldBeEmpty)
	})
}

func TestService_reflectIntVal(t *testing.T) {
	var (
		a = &monitor.BinlogArchive{
			ID:     123,
			State:  0,
			Round:  10,
			MID:    666,
			TypeID: 22,
			Addit: &monitor.ArchiveAddit{
				MissionID: 999,
			},
		}
	)
	Convey("reflectIntVal", t, func(ctx C) {
		_, err := s.reflectIntVal(a, "Addit.MissionID", 0)
		So(err, ShouldBeNil)

		_, err = s.reflectIntVal(a, "Addit111", 0)
		So(err, ShouldNotBeNil)

		_, err = s.reflectIntVal(a, "ID", 0)
		So(err, ShouldBeNil)
	})
}

func TestService_monitorCompSatisfy(t *testing.T) {
	Convey("monitorCompSatisfy >=", t, func(ctx C) {
		is, err := s.monitorCompSatisfy(">=10", 11)
		So(err, ShouldBeNil)
		So(is, ShouldBeTrue)

		is, err = s.monitorCompSatisfy(">10", 10)
		So(err, ShouldBeNil)
		So(is, ShouldBeFalse)

		is, err = s.monitorCompSatisfy("=10", 10)
		So(err, ShouldBeNil)
		So(is, ShouldBeTrue)

		is, err = s.monitorCompSatisfy("in(10,20,30)", 10)
		So(err, ShouldBeNil)
		So(is, ShouldBeTrue)

		is, err = s.monitorCompSatisfy("in(10,20,30)", 40)
		So(err, ShouldBeNil)
		So(is, ShouldBeFalse)
	})
}

func TestService_monitorSave(t *testing.T) {
	Convey("monitorSave", t, func(ctx C) {
		_, errs := s.monitorSave([]string{"monitor_test_1"}, []string{"monitor_test_2"}, 123)
		So(errs, ShouldBeEmpty)
	})
}

func TestService_multiAccounts(t *testing.T) {
	var c = context.Background()
	Convey("multiAccounts", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		mock := accApi.NewMockAccountClient(mockCtrl)
		s.acc = mock
		mockReq := &accApi.MidReq{
			Mid: 123,
		}
		mock.EXPECT().ProfileWithStat3(gomock.Any(), mockReq).Return(&accApi.ProfileStatReply{}, nil)
		_, err := s.multiAccounts(c, []int64{123})
		So(err, ShouldBeNil)
	}))
}
