package article

import (
	"context"

	"testing"

	artMdl "go-common/app/interface/main/creative/model/article"
	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestArticleAddDraft(t *testing.T) {
	var (
		c   = context.TODO()
		art = &artMdl.ArtParam{}
	)
	convey.Convey("AddDraft", t, func(ctx convey.C) {
		id, err := d.AddDraft(c, art)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotEqual, ecode.ArtCreationNoPrivilege)
			ctx.So(id, convey.ShouldEqual, 0)
		})
	})
}

func TestArticleDelDraft(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
		mid = int64(0)
		ip  = ""
	)
	convey.Convey("DelDraft", t, func(ctx convey.C) {
		err := d.DelDraft(c, aid, mid, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotEqual, ecode.ArtCreationNoPrivilege)
		})
	})
}

func TestArticleDraft(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(0)
		mid = int64(0)
		ip  = ""
	)
	convey.Convey("Draft", t, func(ctx convey.C) {
		res, err := d.Draft(c, aid, mid, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotEqual, ecode.ArtCreationNoPrivilege)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestArticleDrafts(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		pn  = int(0)
		ps  = int(0)
		ip  = ""
	)
	convey.Convey("Drafts", t, func(ctx convey.C) {
		res, err := d.Drafts(c, mid, pn, ps, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotEqual, ecode.ArtCreationNoPrivilege)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}
