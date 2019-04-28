package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Stat(t *testing.T) {

	var (
		err error
		c   = context.TODO()
	)

	Convey("Stat", t, WithService(func(s *Service) {
		_, err = s.Stat(c, "music_songs")
		So(err, ShouldBeNil)
	}))
}
