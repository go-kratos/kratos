package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_InsertTagIncome(t *testing.T) {
	Convey("growup-job InsertTagIncome", t, WithService(func(s *Service) {
		date := time.Now().Add(-24 * time.Hour)
		t := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		err := s.InsertTagIncome(context.Background(), t)
		So(err, ShouldBeNil)
	}))
}

func Test_GetAvIncomeStatis(t *testing.T) {
	Convey("growup-job GetAvIncomeStatis", t, WithService(func(s *Service) {
		date := time.Now().Add(-24 * time.Hour)
		t := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		err := s.GetAvIncomeStatis(context.Background(), t.Format(_layout))
		So(err, ShouldBeNil)
	}))
}

func Test_GetAvIncomes(t *testing.T) {
	Convey("growup-job GetAvIncome", t, WithService(func(s *Service) {
		res, err := s.GetAvIncome(context.Background())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func Test_GetUpIncomeStatis(t *testing.T) {
	Convey("growup-job GetUpIncomeStatis", t, WithService(func(s *Service) {
		date := time.Now().Add(-24 * time.Hour)
		t := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		hasWithdraw := 1
		err := s.GetUpIncomeStatis(context.Background(), t.Format(_layout), hasWithdraw)
		So(err, ShouldBeNil)
	}))
}

func Test_GetUpAccount(t *testing.T) {
	Convey("growup-job GetUpAccount", t, WithService(func(s *Service) {
		date := time.Now().Add(-24 * time.Hour)
		t := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		res, err := s.GetUpAccount(context.Background(), t.Format(_layout), date.Format("2006-01-02 15:04:05"))
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func Test_GetUpIncome(t *testing.T) {
	Convey("growup-job GetUpIncome", t, WithService(func(s *Service) {
		date := time.Now().Add(-24 * time.Hour)
		t := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		table := "up_account"
		res, err := s.GetUpIncome(context.Background(), table, t.Format(_layout))
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func Test_GetUpWithdraw(t *testing.T) {
	Convey("growup-job GetUpWithdraw", t, WithService(func(s *Service) {
		date := time.Now().Add(-24 * time.Hour)
		t := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		res, err := s.GetUpWithdraw(context.Background(), t.Format(_layout))
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func Test_GetUpNickname(t *testing.T) {
	Convey("growup-job GetUpNickname", t, WithService(func(s *Service) {
		mids := []int64{int64(1), int64(2)}
		res, err := s.GetUpNickname(context.Background(), mids)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}
