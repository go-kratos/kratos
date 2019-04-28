package v2

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestV2Tags(t *testing.T) {
	var (
		appID = int64(24)
	)
	convey.Convey("Tags", t, func(ctx convey.C) {
		tags, err := d.Tags(appID)
		ctx.Convey("Then err should be nil.tags should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tags, convey.ShouldNotBeNil)
		})
	})
}

func TestV2ConfIDs(t *testing.T) {
	var (
		ID = int64(460)
	)
	convey.Convey("ConfIDs", t, func(ctx convey.C) {
		ids, err := d.ConfIDs(ID)
		ctx.Convey("Then err should be nil.ids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ids, convey.ShouldNotBeNil)
		})
	})
}
