package dao

import (
	"context"
	"fmt"
	"go-common/app/admin/main/vip/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaojoinPoolCondition(t *testing.T) {
	convey.Convey("joinPoolCondition", t, func(ctx convey.C) {
		var (
			sqlStr = ""
			q      = &model.ResoucePoolBo{}
			pn     = int(0)
			ps     = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.joinPoolCondition(sqlStr, q, pn, ps)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaojoinHistoryCondition(t *testing.T) {
	convey.Convey("joinHistoryCondition", t, func(ctx convey.C) {
		var (
			sql     = ""
			u       = &model.UserChangeHistoryReq{}
			iscount bool
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.joinHistoryCondition(sql, u, iscount)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSelPoolByName(t *testing.T) {
	convey.Convey("SelPoolByName", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			name = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.SelPoolByName(c, name)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelCountPool(t *testing.T) {
	convey.Convey("SelCountPool", t, func(ctx convey.C) {
		var (
			c = context.Background()
			r = &model.ResoucePoolBo{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.SelCountPool(c, r)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSelPool(t *testing.T) {
	convey.Convey("SelPool", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			r  = &model.ResoucePoolBo{}
			pn = int(0)
			ps = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.SelPool(c, r, pn, ps)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSelPoolRow(t *testing.T) {
	convey.Convey("SelPoolRow", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.SelPoolRow(c, id)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddPool(t *testing.T) {
	convey.Convey("AddPool", t, func(ctx convey.C) {
		var (
			c = context.Background()
			r = &model.ResoucePoolBo{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			a, err := d.AddPool(c, r)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdatePool(t *testing.T) {
	convey.Convey("UpdatePool", t, func(ctx convey.C) {
		var (
			c = context.Background()
			r = &model.ResoucePoolBo{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			a, err := d.UpdatePool(c, r)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSelBatchRow(t *testing.T) {
	convey.Convey("SelBatchRow", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.SelBatchRow(c, id)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelBatchRows(t *testing.T) {
	convey.Convey("SelBatchRows", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			poolID = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.SelBatchRows(c, poolID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddBatch(t *testing.T) {
	convey.Convey("AddBatch", t, func(ctx convey.C) {
		var (
			c = context.Background()
			r = &model.ResouceBatchBo{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			a, err := d.AddBatch(c, r)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateBatch(t *testing.T) {
	convey.Convey("UpdateBatch", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			r   = &model.VipResourceBatch{}
			ver = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			a, err := d.UpdateBatch(c, r, ver)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUseBatch(t *testing.T) {
	convey.Convey("UseBatch", t, func(ctx convey.C) {
		var (
			r   = &model.VipResourceBatch{}
			ver = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			tx, err := d.BeginTran(context.Background())
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tx, convey.ShouldNotBeNil)
			a, err := d.UseBatch(tx, r, ver)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSelBusiness(t *testing.T) {
	convey.Convey("SelBusiness", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.SelBusiness(c, id)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelBusinessByQuery(t *testing.T) {
	convey.Convey("SelBusinessByQuery", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.QueryBusinessInfo{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			r, err := d.SelBusinessByQuery(c, arg)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAllVersion(t *testing.T) {
	convey.Convey("AllVersion", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.AllVersion(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateVersion(t *testing.T) {
	convey.Convey("UpdateVersion", t, func(ctx convey.C) {
		var (
			c = context.Background()
			v = &model.VipAppVersion{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ret, err := d.UpdateVersion(c, v)
			ctx.Convey("Then err should be nil.ret should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ret, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBussinessList(t *testing.T) {
	convey.Convey("BussinessList", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			pn     = int(0)
			ps     = int(0)
			status = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.BussinessList(c, pn, ps, status)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBussinessCount(t *testing.T) {
	convey.Convey("BussinessCount", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			status = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.BussinessCount(c, status)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateBusiness(t *testing.T) {
	convey.Convey("UpdateBusiness", t, func(ctx convey.C) {
		var (
			c = context.Background()
			r = &model.VipBusinessInfo{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			a, err := d.UpdateBusiness(c, r)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddBusiness(t *testing.T) {
	convey.Convey("AddBusiness", t, func(ctx convey.C) {
		var (
			c = context.Background()
			r = &model.VipBusinessInfo{
				AppKey: fmt.Sprintf("a:%d", time.Now().Unix()),
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			a, err := d.AddBusiness(c, r)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoHistoryCount(t *testing.T) {
	convey.Convey("HistoryCount", t, func(ctx convey.C) {
		var (
			c = context.Background()
			u = &model.UserChangeHistoryReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.HistoryCount(c, u)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoHistoryList(t *testing.T) {
	convey.Convey("HistoryList", t, func(ctx convey.C) {
		var (
			c = context.Background()
			u = &model.UserChangeHistoryReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.HistoryList(c, u)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
