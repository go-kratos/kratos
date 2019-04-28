package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/admin/main/coupon/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCountCode(t *testing.T) {
	convey.Convey("CountCode", t, func(convCtx convey.C) {
		var (
			c = context.Background()
			a = &model.ArgCouponCode{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.CountCode(c, a)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCodeList(t *testing.T) {
	convey.Convey("CodeList", t, func(convCtx convey.C) {
		var (
			c = context.Background()
			a = &model.ArgCouponCode{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.CodeList(c, a)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateCodeBlock(t *testing.T) {
	convey.Convey("UpdateCodeBlock", t, func(convCtx convey.C) {
		var (
			c = context.Background()
			a = &model.CouponCode{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.UpdateCodeBlock(c, a)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCodeByID(t *testing.T) {
	convey.Convey("CodeByID", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.CodeByID(c, id)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBatchAddCode(t *testing.T) {
	convey.Convey("BatchAddCode", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			cs = []*model.CouponCode{
				{
					BatchToken: fmt.Sprintf("%d", time.Now().Unix()),
					Code:       fmt.Sprintf("%d", time.Now().Unix()),
					CouponType: 3,
					State:      model.CodeStateNotUse,
				},
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.BatchAddCode(c, cs)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaowhereSQL(t *testing.T) {
	convey.Convey("whereSQL", t, func(convCtx convey.C) {
		var (
			a = &model.ArgCouponCode{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			sql := whereSQL(a)
			convCtx.Convey("Then sql should not be nil.", func(convCtx convey.C) {
				convCtx.So(sql, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaopageSQL(t *testing.T) {
	convey.Convey("pageSQL", t, func(convCtx convey.C) {
		var (
			pn = int(0)
			ps = int(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			sql := pageSQL(pn, ps)
			convCtx.Convey("Then sql should not be nil.", func(convCtx convey.C) {
				convCtx.So(sql, convey.ShouldNotBeNil)
			})
		})
	})
}
