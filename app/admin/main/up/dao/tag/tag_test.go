package tag

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestTagTagList(t *testing.T) {
	convey.Convey("TagList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{}
		)
		defer gock.OffAll()
		httpMock("GET", d.tagList).Reply(200).BodyString(`
{
	"code":0, 
	"data": []
}`)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tgs, err := d.TagList(c, ids)
			ctx.Convey("Then err should be nil.tgs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tgs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagTagCheck(t *testing.T) {
	convey.Convey("TagCheck", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(0)
			tagName = ""
		)
		defer gock.OffAll()
		httpMock("GET", d.tagCheck).Reply(200).BodyString(`
{
	"code":0, 
	"data": {"tags" : []}
}`)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			no, err := d.TagCheck(c, mid, tagName)
			ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(no, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagAppealTag(t *testing.T) {
	convey.Convey("AppealTag", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
			ip  = ""
		)
		defer gock.OffAll()
		httpMock("GET", d.appealTag).Reply(200).BodyString(`
{
	"code":0, 
	"data": {"tags" : []}
}`)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tid, err := d.AppealTag(c, aid, ip)
			ctx.Convey("Then err should be nil.tid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tid, convey.ShouldNotBeNil)
			})
		})
	})
}
