package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_BangumiFeed(t *testing.T) {
	Convey("bangumi feed", t, WithService(t, func(svf *Service) {
		res, err := svf.BangumiFeed(context.TODO(), _bangumiMid, 1, 2, _ip)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		Convey("return feed for page 2", func() {
			time.Sleep(time.Millisecond * 300) // wait cache ready
			res, err := svf.BangumiFeed(context.TODO(), _bangumiMid, 2, 2, _ip)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		})
		Convey("bangumi feed cache", func() {
			time.Sleep(time.Millisecond * 300) // wait cache ready
			res, err := svf.bangumiFeedCache(context.TODO(), _bangumiMid, 3, 2, _ip)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		})
	}))
}
