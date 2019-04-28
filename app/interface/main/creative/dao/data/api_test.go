package data

import (
	"context"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"net/url"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDatastat(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		ip  = ""
	)
	convey.Convey("stat1", t, func(ctx convey.C) {
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.client), "RESTfulGet",
			func(_ *bm.Client, _ context.Context, _ string, _ string, _ url.Values, _ interface{}, _ ...interface{}) (err error) {
				return ecode.CreativeDataErr
			})
		defer mock.Unpatch()
		_, err := d.stat(c, mid, ip)
		ctx.Convey("Then err should be nil.st should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})

}

func TestDataTagsWithChecked(t *testing.T) {
	var (
		c        = context.TODO()
		mid      = int64(1)
		tid      = uint16(1)
		title    = ""
		filename = ""
		desc     = ""
		cover    = ""
		tagFrom  = int8(0)
		err      error
	)
	convey.Convey("TagsWithChecked1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.tagV2URI).Reply(200).JSON(`{"code":20003}`)
		_, err = d.TagsWithChecked(c, mid, tid, title, filename, desc, cover, tagFrom)
		ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("TagsWithChecked2", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.tagV2URI).Reply(200).JSON(`{"code":0,"data":{"tags":[{"tag":"1","checked":1}]}}`)
		_, err = d.TagsWithChecked(c, mid, tid, title, filename, desc, cover, tagFrom)
		ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataRecommendCovers(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		fns = []string{}
		err error
	)
	convey.Convey("RecommendCovers1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.coverBFSURI).Reply(200).JSON(`{"code":20003}`)
		_, err = d.RecommendCovers(c, mid, fns)
		ctx.Convey("Then err should be nil.cvs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("RecommendCovers2", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.coverBFSURI).Reply(200).JSON(`{"code":0,"data":["bfs1"]}`)
		_, err = d.RecommendCovers(c, mid, fns)
		ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
