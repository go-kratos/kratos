package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/growup/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetWithdraw(t *testing.T) {
	var (
		dateVersion   = "2018-04"
		offset, limit = 0, 10
	)
	Convey("GetWithdraw", t, WithService(func(s *Service) {
		_, res, err := s.GetWithdraw(context.Background(), dateVersion, offset, limit)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func Test_UpWithdraw(t *testing.T) {
	var (
		dateVersion   = "2018-04"
		offset, limit = 0, 10
	)
	Convey("UpWithdraw", t, WithService(func(s *Service) {
		_, res, err := s.UpWithdraw(context.Background(), dateVersion, offset, limit)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func Test_InsertUpWithdrawRecord(t *testing.T) {
	var (
		a = &model.UpIncomeWithdraw{MID: int64(1001), WithdrawIncome: 10, DateVersion: "2018-05", State: 1}
	)
	Convey("Test_InsertUpWithdrawRecord", t, WithService(func(s *Service) {
		err := s.InsertUpWithdrawRecord(context.Background(), a)
		So(err, ShouldBeNil)
	}))
}

func Test_WithdrawSuccess(t *testing.T) {
	var (
		orderNo     int64 = 10
		tradeStatus       = 10
	)
	Convey("WithdrawSuccess", t, WithService(func(s *Service) {
		err := s.WithdrawSuccess(context.Background(), orderNo, tradeStatus)
		So(err, ShouldBeNil)
	}))
}

func Test_WithdrawDetail(t *testing.T) {
	var (
		mid int64 = 10
	)
	Convey("WithdrawDetail", t, WithService(func(s *Service) {
		res, err := s.WithdrawDetail(context.Background(), mid)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}
