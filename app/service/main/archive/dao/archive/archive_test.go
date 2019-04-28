package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveMaxAID(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("MaxAID", t, func(ctx convey.C) {
		id, err := d.MaxAID(c)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivearchive3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10097272)
	)
	convey.Convey("archive3", t, func(ctx convey.C) {
		_, err := d.archive3(c, aid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchivearchives3(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{10097272}
	)
	convey.Convey("archives3", t, func(ctx convey.C) {
		res, err := d.archives3(c, aids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
