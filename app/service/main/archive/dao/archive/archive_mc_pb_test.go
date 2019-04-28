package archive

import (
	"context"
	"go-common/app/service/main/archive/api"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchivedescKey(t *testing.T) {
	var (
		aid = int64(1)
	)
	convey.Convey("descKey", t, func(ctx convey.C) {
		p1 := descKey(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivearcPBKey(t *testing.T) {
	var (
		aid = int64(1)
	)
	convey.Convey("arcPBKey", t, func(ctx convey.C) {
		p1 := arcPBKey(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivearchive3Cache(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
	)
	convey.Convey("archive3Cache", t, func(ctx convey.C) {
		d.archive3Cache(c, aid)
	})
}

func TestArchivearchive3Caches(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{1, 2}
	)
	convey.Convey("archive3Caches", t, func(ctx convey.C) {
		cached, err := d.archive3Caches(c, aids)
		ctx.Convey("Then err should be nil.cached should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(cached, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveaddArchive3Cache(t *testing.T) {
	var (
		c = context.TODO()
		a = &api.Arc{}
	)
	convey.Convey("addArchive3Cache", t, func(ctx convey.C) {
		err := d.addArchive3Cache(c, a)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchivedescCache(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
	)
	convey.Convey("descCache", t, func(ctx convey.C) {
		desc, err := d.descCache(c, aid)
		ctx.Convey("Then err should be nil.desc should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(desc, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveaddDescCache(t *testing.T) {
	var (
		c    = context.TODO()
		aid  = int64(1)
		desc = ""
	)
	convey.Convey("addDescCache", t, func(ctx convey.C) {
		err := d.addDescCache(c, aid, desc)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
