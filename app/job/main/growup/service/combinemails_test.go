package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_CombineMailsByHTTP(t *testing.T) {
	Convey("growup-job", t, WithService(func(s *Service) {
		var (
			year  = 2018
			month = 5
			day   = 16
		)
		err := s.CombineMailsByHTTP(context.Background(), year, month, day)
		So(err, ShouldBeNil)
	}))
}

func Test_CombineMails(t *testing.T) {
	Convey("growup-job", t, WithService(func(s *Service) {
		err := s.CombineMails()
		So(err, ShouldBeNil)
	}))
}
