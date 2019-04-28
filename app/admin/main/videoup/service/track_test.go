package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_trackArchive(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("TrackArchive", t, WithService(func(s *Service) {
		_, err := svr.TrackArchive(c, 1)
		So(err, ShouldBeNil)
	}))
}

func TestService_TrackVideo(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("TrackVideo", t, WithService(func(s *Service) {
		_, err := svr.TrackVideo(c, "11111", 111)
		So(err, ShouldBeNil)
	}))
}
