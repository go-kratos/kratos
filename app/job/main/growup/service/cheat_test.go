package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpdateCheatHTTP(t *testing.T) {
	Convey("growup-job UpdateCheatHTTP", t, WithService(func(s *Service) {
		date := time.Now().Add(-24 * time.Hour)
		err := s.UpdateCheatHTTP(context.Background(), date)
		So(err, ShouldBeNil)
	}))
}

func Test_CheatStatistics(t *testing.T) {
	Convey("growup-job CheatStatistics", t, WithService(func(s *Service) {
		date := time.Now().Add(-24 * time.Hour)
		err := s.CheatStatistics(context.Background(), date)
		So(err, ShouldBeNil)
	}))
}
