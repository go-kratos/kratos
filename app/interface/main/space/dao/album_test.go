package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_AlbumCount(t *testing.T) {
	convey.Convey("test album count", t, func(ctx convey.C) {
		mid := int64(28272030)
		data, err := d.AlbumCount(context.Background(), mid)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%d", data)
	})
}

func TestDao_AlbumList(t *testing.T) {
	convey.Convey("test album list", t, func(ctx convey.C) {
		mid := int64(28272030)
		pn := 0
		ps := 1
		data, err := d.AlbumList(context.Background(), mid, pn, ps)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%+v", data)
	})
}
