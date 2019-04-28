package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_CheatUps(t *testing.T) {
	mid := int64(1011)
	nickname := "helloworld"
	from, limit := 0, 1000
	Convey("admins", t, WithService(func(s *Service) {
		_, _, err := s.CheatUps(context.Background(), mid, nickname, from, limit)
		So(err, ShouldBeNil)
	}))
}

func Test_CheatArchives(t *testing.T) {
	mid, avID := int64(1011), int64(5011)
	nickname := "helloworld"
	from, limit := 0, 1000
	Convey("admins", t, WithService(func(s *Service) {
		_, _, err := s.CheatArchives(context.Background(), mid, avID, nickname, from, limit)
		So(err, ShouldBeNil)
	}))
}

func Test_ExportCheatUps(t *testing.T) {
	mid := int64(1011)
	nickname := "helloworld"
	from, limit := 0, 1000
	Convey("admins", t, WithService(func(s *Service) {
		res, err := s.ExportCheatUps(context.Background(), mid, nickname, from, limit)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func Test_ExportCheatAvs(t *testing.T) {
	mid, avID := int64(1011), int64(5011)
	nickname := "helloworld"
	from, limit := 0, 1000
	Convey("admins", t, WithService(func(s *Service) {
		res, err := s.ExportCheatAvs(context.Background(), mid, avID, nickname, from, limit)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}
