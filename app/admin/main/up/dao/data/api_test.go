package data

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDataTags(t *testing.T) {
	convey.Convey("Tags", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			tid      = uint16(0)
			title    = ""
			filename = ""
			desc     = ""
			cover    = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.tagURI).Reply(200).BodyString(`
{
	"code":0, 
	"data": {"tags" : ["1"]}
}`)
			no, err := d.Tags(c, mid, tid, title, filename, desc, cover)
			ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(no, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataTagsWithChecked(t *testing.T) {
	convey.Convey("TagsWithChecked", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			tid      = uint16(0)
			title    = ""
			filename = ""
			desc     = ""
			cover    = ""
			tagFrom  = int8(0)
		)
		defer gock.OffAll()
		httpMock("GET", d.tagV2URI).Reply(200).BodyString(`
{
	"code":0, 
	"data": {"tags" : []}
}`)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			no, err := d.TagsWithChecked(c, mid, tid, title, filename, desc, cover, tagFrom)
			ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(no, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataRecommendCovers(t *testing.T) {
	convey.Convey("RecommendCovers", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			fns = []string{}
		)
		defer gock.OffAll()
		httpMock("GET", d.coverBFSURI).Reply(200).BodyString(`
{
	"code":0, 
	"data": ["1"]
}`)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cvs, err := d.RecommendCovers(c, mid, fns)
			ctx.Convey("Then err should be nil.cvs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cvs, convey.ShouldNotBeNil)
			})
		})
	})
}
