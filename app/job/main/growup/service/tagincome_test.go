package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SendTagIncomeByHTTP(t *testing.T) {
	Convey("growup-job", t, WithService(func(s *Service) {
		var (
			year  = 2018
			month = 1
			day   = 26
		)
		err := s.SendTagIncomeByHTTP(context.Background(), year, month, day)
		So(err, ShouldBeNil)
	}))
}
