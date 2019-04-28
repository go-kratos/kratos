package anchorReward

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAnchorRewardpingMC(t *testing.T) {
	convey.Convey("pingMC", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.pingMC(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestAnchorRewardnewRewardKey(t *testing.T) {
	convey.Convey("newRewardKey", t, func(ctx convey.C) {
		var (
			uid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := newRewardKey(uid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAnchorRewardSetNewReward(t *testing.T) {
	convey.Convey("SetNewReward", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			uid    = int64(0)
			hasNew = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetNewReward(c, uid, hasNew)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestAnchorRewardGetNewReward(t *testing.T) {
	convey.Convey("GetNewReward", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			uid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			r, err := d.GetNewReward(c, uid)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAnchorRewardGetHasReward(t *testing.T) {
	convey.Convey("GetHasReward", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			uid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			r, err := d.GetHasReward(c, uid)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAnchorRewardhasRewardKey(t *testing.T) {
	convey.Convey("hasRewardKey", t, func(ctx convey.C) {
		var (
			uid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := hasRewardKey(uid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAnchorRewardSetHasReward(t *testing.T) {
	convey.Convey("SetHasReward", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			uid    = int64(0)
			isHave = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetHasReward(c, uid, isHave)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestAnchorRewardGetOrder(t *testing.T) {
	convey.Convey("GetOrder", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			exists, err := d.GetOrder(c, id)
			ctx.Convey("Then err should be nil.exists should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(exists, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAnchorRewardSaveOrder(t *testing.T) {
	convey.Convey("SaveOrder", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SaveOrder(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
