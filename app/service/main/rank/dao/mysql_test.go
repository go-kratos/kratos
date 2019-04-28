package dao

import (
	"context"
	"testing"

	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoMaxOid(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("MaxOid", t, func(ctx convey.C) {
		oid, err := d.MaxOid(c)
		ctx.Convey("Then err should be nil.oid should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(oid, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoArchiveMetas(t *testing.T) {
	var (
		c     = context.Background()
		id    = int64(1)
		limit = int(1)
	)
	convey.Convey("ArchiveMetas", t, func(ctx convey.C) {
		p1, err := d.ArchiveMetas(c, id, limit)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoArchiveMetasIncrs(t *testing.T) {
	var (
		c     = context.Background()
		id    = int64(1)
		begin xtime.Time
		end   xtime.Time
		limit = int(1)
	)
	convey.Convey("ArchiveMetasIncrs", t, func(ctx convey.C) {
		p1, err := d.ArchiveMetasIncrs(c, id, begin, end, limit)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoArchiveTypes(t *testing.T) {
	var (
		c   = context.Background()
		ids = []int64{1}
	)
	convey.Convey("ArchiveTypes", t, func(ctx convey.C) {
		p1, err := d.ArchiveTypes(c, ids)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoArchiveStats(t *testing.T) {
	var (
		c    = context.Background()
		aids = []int64{1}
	)
	convey.Convey("ArchiveStats", t, func(ctx convey.C) {
		p1, err := d.ArchiveStats(c, aids)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoArchiveStatsIncrs(t *testing.T) {
	var (
		c     = context.Background()
		tbl   = int(1)
		id    = int64(1)
		begin xtime.Time
		end   xtime.Time
		limit = int(1)
	)
	convey.Convey("ArchiveStatsIncrs", t, func(ctx convey.C) {
		p1, err := d.ArchiveStatsIncrs(c, tbl, id, begin, end, limit)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoArchiveTVs(t *testing.T) {
	var (
		c    = context.Background()
		aids = []int64{1}
	)
	convey.Convey("ArchiveTVs", t, func(ctx convey.C) {
		p1, err := d.ArchiveTVs(c, aids)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoArchiveTVsIncrs(t *testing.T) {
	var (
		c     = context.Background()
		id    = int64(1)
		begin xtime.Time
		end   xtime.Time
		limit = int(1)
	)
	convey.Convey("ArchiveTVsIncrs", t, func(ctx convey.C) {
		p1, err := d.ArchiveTVsIncrs(c, id, begin, end, limit)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
