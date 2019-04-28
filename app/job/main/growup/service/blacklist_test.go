package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_InitBlacklistMID(t *testing.T) {
	Convey("growup-job InitBlacklistMID", t, WithService(func(s *Service) {
		err := s.InitBlacklistMID(context.Background())
		So(err, ShouldBeNil)
	}))
}

func Test_UpdateBlacklist(t *testing.T) {
	Convey("growup-job UpdateBlacklist", t, WithService(func(s *Service) {
		err := s.UpdateBlacklist(context.Background())
		So(err, ShouldBeNil)
	}))
}

func Test_GetAvsMID(t *testing.T) {
	Convey("growup-job GetAvsMID", t, WithService(func(s *Service) {
		var (
			avs = []int64{int64(1), int64(2)}
		)
		_, err := s.GetAvsMID(context.Background(), avs)
		So(err, ShouldBeNil)
	}))
}
