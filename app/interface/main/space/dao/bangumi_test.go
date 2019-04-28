package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_BangumiList(t *testing.T) {
	convey.Convey("test bangumi list", t, func(ctx convey.C) {
		data, count, err := d.BangumiList(context.Background(), 908085, 1, 10)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%+v", data)
		convey.Printf("%d", count)
	})
}
