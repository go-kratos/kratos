package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_MaxLikeCache(t *testing.T) {
	var (
		aid   = int64(100)
		value = int64(200)
		err   error
	)
	Convey("add cache", t, WithCleanCache(func() {
		err = d.SetMaxLikeCache(context.TODO(), aid, value)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			res, err := d.MaxLikeCache(context.TODO(), aid)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, value)
		})
		Convey("expire cache", func() {
			res, err := d.ExpireMaxLikeCache(context.TODO(), aid)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, true)
		})
	}))
}
