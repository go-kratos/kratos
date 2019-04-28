package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DeleteAvRatio(t *testing.T) {
	Convey("growup-job", t, WithService(func(s *Service) {
		var (
			limit = int64(200)
		)
		_, err := s.DeleteAvRatio(context.Background(), limit)
		So(err, ShouldBeNil)
	}))
}

func Test_DeleteUpIncome(t *testing.T) {
	Convey("growup-job", t, WithService(func(s *Service) {
		var (
			limit = int64(200)
		)
		_, err := s.DeleteUpIncome(context.Background(), limit)
		So(err, ShouldBeNil)
	}))
}

func Test_UpdateAvRatio(t *testing.T) {
	Convey("growup-job", t, WithService(func(s *Service) {
		err := s.UpdateAvRatio()
		So(err, ShouldBeNil)
	}))
}

func Test_ExecIncomeForHTTP(t *testing.T) {
	Convey("growup-job", t, WithService(func(s *Service) {
		var (
			year  = 2018
			month = 1
			day   = 12
		)
		err := s.ExecIncomeForHTTP(context.Background(), year, month, day)
		So(err, ShouldBeNil)
	}))
}

func Test_ExecRatioForHTTP(t *testing.T) {
	Convey("growup-job", t, WithService(func(s *Service) {
		var (
			year  = 2018
			month = 1
			day   = 12
		)
		err := s.ExecRatioForHTTP(context.Background(), year, month, day)
		So(err, ShouldBeNil)
	}))
}
