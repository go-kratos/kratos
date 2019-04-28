package academy

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAcademyTagList(t *testing.T) {
	convey.Convey("TagList", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, tagMap, parentChildMap, err := d.TagList(c)
			ctx.Convey("Then err should be nil.res,tagMap,parentChildMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(parentChildMap, convey.ShouldNotBeNil)
				ctx.So(tagMap, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAcademyTag(t *testing.T) {
	convey.Convey("Tag", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			no, err := d.Tag(c, id)
			ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(no, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAcademyTags(t *testing.T) {
	convey.Convey("Tags", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Tags(c, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAcademyLinkTags(t *testing.T) {
	convey.Convey("LinkTags", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LinkTags(c, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
