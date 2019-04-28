package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/playlist/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyPlArc(t *testing.T) {
	var (
		pid = int64(1)
	)
	convey.Convey("keyPlArc", t, func(ctx convey.C) {
		p1 := keyPlArc(pid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaokeyPlArcDesc(t *testing.T) {
	var (
		pid = int64(1)
	)
	convey.Convey("keyPlArcDesc", t, func(ctx convey.C) {
		p1 := keyPlArcDesc(pid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoArcsCache(t *testing.T) {
	var (
		c     = context.Background()
		pid   = int64(1)
		start = int(1)
		end   = int(20)
	)
	convey.Convey("ArcsCache", t, func(ctx convey.C) {
		_, err := d.ArcsCache(c, pid, start, end)
		ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddArcCache(t *testing.T) {
	var (
		c   = context.Background()
		pid = int64(1)
		arc = &model.ArcSort{Aid: 13825646, Sort: 100, Desc: "abc"}
	)
	convey.Convey("AddArcCache", t, func(ctx convey.C) {
		err := d.AddArcCache(c, pid, arc)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetArcsCache(t *testing.T) {
	var (
		c    = context.Background()
		pid  = int64(1)
		arcs = []*model.ArcSort{}
	)
	convey.Convey("SetArcsCache", t, func(ctx convey.C) {
		arcs = append(arcs, &model.ArcSort{Aid: 13825646, Sort: 100, Desc: "abc"})
		err := d.SetArcsCache(c, pid, arcs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetArcDescCache(t *testing.T) {
	var (
		c    = context.Background()
		pid  = int64(1)
		aid  = int64(13825646)
		desc = "abc"
	)
	convey.Convey("SetArcDescCache", t, func(ctx convey.C) {
		err := d.SetArcDescCache(c, pid, aid, desc)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelArcsCache(t *testing.T) {
	var (
		c    = context.Background()
		pid  = int64(1)
		aids = []int64{13825646, 11, 100}
	)
	convey.Convey("DelArcsCache", t, func(ctx convey.C) {
		err := d.DelArcsCache(c, pid, aids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelCache(t *testing.T) {
	var (
		c   = context.Background()
		pid = int64(1)
	)
	convey.Convey("DelCache", t, func(ctx convey.C) {
		err := d.DelCache(c, pid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
