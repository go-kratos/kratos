package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Replys(t *testing.T) {
	Convey("Replys", t, WithService(func(s *Service) {
		_, _, err := s.Replys(context.Background(), "1466046600", "ios", "0", "phone", "abc", "123", "android-creative", 11111, 1, 10)
		So(err, ShouldBeNil)
	}))
}

func Test_Sessions(t *testing.T) {
	Convey("Sessions", t, WithService(func(s *Service) {
		_, _, err := s.Sessions(context.Background(), 1, "2", "3", "ios", time.Now(), time.Now(), 1, 10)
		So(err, ShouldBeNil)
	}))
}

func Test_Tags(t *testing.T) {
	Convey("Tags", t, WithService(func(s *Service) {
		_, err := s.Tags(context.Background(), 1, 1, "ios")
		So(err, ShouldBeNil)
	}))
}

func Test_WebReplys(t *testing.T) {
	Convey("WebReplys", t, WithService(func(s *Service) {
		_, err := s.WebReplys(context.Background(), 1, 1)
		So(err, ShouldBeNil)
	}))
}
