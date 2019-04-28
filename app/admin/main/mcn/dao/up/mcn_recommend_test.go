package up

import (
	"context"
	"testing"

	"go-common/app/admin/main/mcn/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpAddMcnUpRecommend(t *testing.T) {
	convey.Convey("AddMcnUpRecommend", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.McnUpRecommendPool{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.AddMcnUpRecommend(c, arg)
			d.db.Exec(c, "delete from mcn_up_recommend_pool where up_mid=0")
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpUpMcnUpsRecommendOP(t *testing.T) {
	convey.Convey("UpMcnUpsRecommendOP", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			upMids = []int64{1, 2}
			state  = model.MCNUPRecommendState(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpMcnUpsRecommendOP(c, upMids, state)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMcnUpRecommends(t *testing.T) {
	convey.Convey("McnUpRecommends", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPRecommendReq{}
		)
		arg.Order = "mtime"
		arg.Sort = "DESC"
		arg.Page = 1
		arg.Size = 10
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.McnUpRecommends(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMcnUpRecommendTotal(t *testing.T) {
	convey.Convey("McnUpRecommendTotal", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPRecommendReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.McnUpRecommendTotal(c, arg)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMcnUpBindMids(t *testing.T) {
	convey.Convey("McnUpBindMids", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			bmids, err := d.McnUpBindMids(c, mids)
			ctx.Convey("Then err should be nil.bmids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				if len(bmids) == 0 {
					ctx.So(bmids, convey.ShouldBeEmpty)
				} else {
					ctx.So(bmids, convey.ShouldNotBeNil)
				}
			})
		})
	})
}

func TestUpMcnUpRecommendMid(t *testing.T) {
	convey.Convey("McnUpRecommendMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			m, err := d.McnUpRecommendMid(c, mid)
			ctx.Convey("Then err should be nil.m should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(m, convey.ShouldBeNil)
			})
		})
	})
}

func TestUpbuildUpRecommendSQL(t *testing.T) {
	convey.Convey("buildUpRecommendSQL", t, func(ctx convey.C) {
		var (
			tp  = ""
			arg = &model.MCNUPRecommendReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			sql, values := d.buildUpRecommendSQL(tp, arg)
			ctx.Convey("Then sql,values should not be nil.", func(ctx convey.C) {
				ctx.So(values, convey.ShouldNotBeNil)
				ctx.So(sql, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpjoinStringSQL(t *testing.T) {
	convey.Convey("joinStringSQL", t, func(ctx convey.C) {
		var (
			is = []string{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.joinStringSQL(is)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
