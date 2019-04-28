package archive

import (
	"context"
	"go-common/app/service/main/archive/api"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchivepagePBKey(t *testing.T) {
	var (
		aid = int64(1)
	)
	convey.Convey("pagePBKey", t, func(ctx convey.C) {
		p1 := pagePBKey(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivevideoPBKey2(t *testing.T) {
	var (
		aid = int64(1)
		cid = int64(1)
	)
	convey.Convey("videoPBKey2", t, func(ctx convey.C) {
		p1 := videoPBKey2(aid, cid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveaddPageCache3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10097272)
		ps  = []*api.Page{}
	)
	convey.Convey("addPageCache3", t, func(ctx convey.C) {
		err := d.addPageCache3(c, aid, ps)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveaddVideoCache3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10097272)
		cid = int64(10097272)
		p   = &api.Page{}
	)
	convey.Convey("addVideoCache3", t, func(ctx convey.C) {
		err := d.addVideoCache3(c, aid, cid, p)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveDelVideoCache3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10097272)
		cid = int64(10097272)
	)
	convey.Convey("DelVideoCache3", t, func(ctx convey.C) {
		err := d.DelVideoCache3(c, aid, cid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchivepageCache3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10097272)
	)
	convey.Convey("pageCache3", t, func(ctx convey.C) {
		_, err := d.pageCache3(c, aid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchivepagesCache3(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{10097272}
	)
	convey.Convey("pagesCache3", t, func(ctx convey.C) {
		_, _, err := d.pagesCache3(c, aids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchivevideoCache3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10097272)
		cid = int64(10097272)
	)
	convey.Convey("videoCache3", t, func(ctx convey.C) {
		_, err := d.videoCache3(c, aid, cid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
