package music

import (
	"context"
	"go-common/app/interface/main/creative/model/search"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"reflect"
	"testing"

	"github.com/bouk/monkey"

	"github.com/smartystreets/goconvey/convey"
)

func TestMusicSearchBgmSIDs(t *testing.T) {
	var (
		c     = context.TODO()
		ret   []int64
		kw    = "kw"
		pn    = 1
		ps    = 10
		page  *search.Pager
		err   error
		esReq *elastic.Request
	)
	convey.Convey("SearchBgmSIDs", t, func(ctx convey.C) {
		monkey.PatchInstanceMethod(reflect.TypeOf(esReq), "Scan", func(_ *elastic.Request, _ context.Context, _ interface{}) (err error) {
			return ecode.CreativeSearchErr
		})
		ret, page, err = d.SearchBgmSIDs(c, kw, pn, ps)
		ctx.Convey("SearchBgmSIDs", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeSearchErr)
			ctx.So(ret, convey.ShouldBeNil)
			ctx.So(page, convey.ShouldBeNil)
		})
	})
}

func TestExtAidsWithSameBgm(t *testing.T) {
	var (
		c     = context.TODO()
		sid   = int64(87600)
		total int
		aids  []int64
		err   error
		esReq *elastic.Request
	)
	convey.Convey("ExtAidsWithSameBgm", t, func(ctx convey.C) {
		monkey.PatchInstanceMethod(reflect.TypeOf(esReq), "Scan", func(_ *elastic.Request, _ context.Context, _ interface{}) (err error) {
			return ecode.CreativeSearchErr
		})
		aids, total, err = d.ExtAidsWithSameBgm(c, sid, 1)
		ctx.Convey("ExtAidsWithSameBgm err should be nil.res,resMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeSearchErr)
			ctx.So(aids, convey.ShouldHaveLength, 0)
			ctx.So(total, convey.ShouldNotBeNil)
		})
	})
}
