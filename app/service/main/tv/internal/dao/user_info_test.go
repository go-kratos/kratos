package dao

import (
	"context"
	"go-common/app/service/main/tv/internal/model"
	"math/rand"
	"testing"
	"time"

	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRawUserInfoByMid(t *testing.T) {
	convey.Convey("RawUserInfoByMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ui, err := d.RawUserInfoByMid(c, mid)
			ctx.Convey("Then err should be nil.ui should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ui, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxInsertUserInfo(t *testing.T) {
	convey.Convey("TxInsertUserInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			ui    = &model.UserInfo{
				Mid:         int64(rand.Int() / 10000000000),
				Ver:         1,
				VipType:     1,
				OverdueTime: xtime.Time(time.Now().Unix() + 60*60*24*31),
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.TxInsertUserInfo(c, tx, ui)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxUpdateUserInfo(t *testing.T) {
	convey.Convey("TxUpdateUserInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			ui    = &model.UserInfo{
				Mid:         27515308,
				Ver:         1,
				VipType:     1,
				OverdueTime: xtime.Time(time.Now().Unix() + 60*60*24*31),
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.TxUpdateUserInfo(c, tx, ui)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxUpdateUserStatus(t *testing.T) {
	convey.Convey("TxUpdateUserStatus", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			ui    = &model.UserInfo{
				Mid:    27515308,
				Status: 1,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.TxUpdateUserStatus(c, tx, ui)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxUpdateUserPayType(t *testing.T) {
	convey.Convey("TxUpdateUserPayType", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			ui    = &model.UserInfo{
				Mid:     27515308,
				PayType: 1,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.TxUpdateUserPayType(c, tx, ui)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}
