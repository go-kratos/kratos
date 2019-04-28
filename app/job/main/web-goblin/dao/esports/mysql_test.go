package esports

import (
	"context"
	"testing"

	mdlesp "go-common/app/job/main/web-goblin/model/esports"
	arcmdl "go-common/app/service/main/archive/api"

	"github.com/smartystreets/goconvey/convey"
)

func TestEsportsContests(t *testing.T) {
	var (
		c     = context.Background()
		stime = int64(1539590040)
		etime = int64(1539590040)
	)
	convey.Convey("Contests", t, func(ctx convey.C) {
		res, err := d.Contests(c, stime, etime)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(len(res), convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestEsportsTeams(t *testing.T) {
	var (
		c      = context.Background()
		homeID = int64(1)
		awayID = int64(2)
	)
	convey.Convey("Teams", t, func(ctx convey.C) {
		res, err := d.Teams(c, homeID, awayID)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(len(res), convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestEsportsArcs(t *testing.T) {
	var (
		c     = context.Background()
		id    = int64(1)
		limit = int(50)
	)
	convey.Convey("Arcs", t, func(ctx convey.C) {
		res, err := d.Arcs(c, id, limit)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestEsportsUpArcScore(t *testing.T) {
	var (
		c        = context.Background()
		partArcs = []*mdlesp.Arc{}
		arcs     map[int64]*arcmdl.Arc
	)
	convey.Convey("UpArcScore", t, func(ctx convey.C) {
		err := d.UpArcScore(c, partArcs, arcs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestEsportsscore(t *testing.T) {
	var (
		arc = &arcmdl.Arc{}
	)
	convey.Convey("score", t, func(ctx convey.C) {
		res := d.score(arc)
		ctx.Convey("Then res should not be nil.", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
