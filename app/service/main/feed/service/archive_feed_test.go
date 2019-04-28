package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ArchiveFeed(t *testing.T) {
	Convey("archive feed", t, WithService(t, func(svf *Service) {
		res, err := svf.ArchiveFeed(context.TODO(), _mid, 1, 2, _ip)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		Convey("return feed for page 2", func() {
			time.Sleep(time.Millisecond * 300) // wait cache ready
			res, err := svf.ArchiveFeed(context.TODO(), _mid, 2, 2, _ip)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		})
	}))
}
