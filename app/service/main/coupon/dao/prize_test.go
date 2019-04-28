package dao

import (
	"context"
	"math/rand"
	"testing"

	"go-common/app/service/main/coupon/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUserCard(t *testing.T) {
	convey.Convey("UserCard", t, func(convCtx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(88895150)
			actID    = int64(1)
			cardType = int8(1)
			noMID    = int64(-1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.UserCard(c, mid, actID, cardType)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldNotBeNil)
			})
			r, err = d.UserCard(c, noMID, actID, cardType)
			convCtx.Convey("Then err should be nil.r should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUserCards(t *testing.T) {
	convey.Convey("UserCards", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(88895150)
			actID = int64(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.UserCards(c, mid, actID)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddUserCard(t *testing.T) {
	convey.Convey("AddUserCard", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(context.Background())
			uc    = &model.CouponUserCard{MID: rand.Int63n(9999999)}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.AddUserCard(c, tx, uc)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateUserCard(t *testing.T) {
	convey.Convey("UpdateUserCard", t, func(convCtx convey.C) {
		var (
			c           = context.Background()
			mid         = int64(0)
			state       = int8(0)
			couponToken = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.UpdateUserCard(c, mid, state, couponToken)
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoClose(t *testing.T) {
	convey.Convey("TestDaoClose", t, func(convCtx convey.C) {
		d.Close()
	})
}
