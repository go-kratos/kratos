package academy

import (
	"context"
	"go-common/app/interface/main/creative/model/academy"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAcademyOccupations(t *testing.T) {
	convey.Convey("Occupations", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Occupations(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAcademySkills(t *testing.T) {
	convey.Convey("Skills", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Skills(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAcademySkillArcs(t *testing.T) {
	convey.Convey("SkillArcs", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			pids   = []int64{}
			skids  = []int64{}
			sids   = []int64{}
			offset = int(0)
			limit  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.SkillArcs(c, pids, skids, sids, offset, limit)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAcademySkillArcCount(t *testing.T) {
	convey.Convey("SkillArcCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			pids  = []int64{}
			skids = []int64{}
			sids  = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.SkillArcCount(c, pids, skids, sids)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAcademyPlayAdd(t *testing.T) {
	convey.Convey("PlayAdd", t, func(ctx convey.C) {
		var (
			c = context.Background()
			p = &academy.Play{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.PlayAdd(c, p)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAcademyPlays(t *testing.T) {
	convey.Convey("Plays", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(0)
			offset = int(0)
			limit  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Plays(c, mid, offset, limit)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAcademyPlayCount(t *testing.T) {
	convey.Convey("PlayCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.PlayCount(c, mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
