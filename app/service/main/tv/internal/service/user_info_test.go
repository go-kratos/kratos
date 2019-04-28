package service

import (
	"context"
	"go-common/app/service/main/tv/internal/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceUserInfo(t *testing.T) {
	convey.Convey("UserInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ui, err := s.UserInfo(c, mid)
			ctx.Convey("Then err should be nil.ui should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ui, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceYstUserInfo(t *testing.T) {
	convey.Convey("YstUserInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &model.YstUserInfoReq{
				Mid: 27515308,
			}
		)
		req.Sign, _ = s.dao.Signer().Sign(req)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ui, err := s.YstUserInfo(c, req)
			ctx.Convey("Then err should be nil.ui should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ui, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceExpireUserAsync(t *testing.T) {
	convey.Convey("ExpireUserAsync", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			s.ExpireUserAsync(c, mid)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestServiceExpireUser(t *testing.T) {
	convey.Convey("ExpireUser", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := s.ExpireUser(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceincrVipDuration(t *testing.T) {
	convey.Convey("incrVipDuration", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = s.dao.BeginTran(c)
			ui    = &model.UserInfo{
				Mid: 27515308,
			}
			d = time.Hour * 24 * 31
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := s.incrVipDuration(c, tx, ui, d)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

//
//func TestServiceRenewPanel(t *testing.T) {
//	convey.Convey("RenewPanel", t, func(ctx convey.C) {
//		var (
//			c   = context.Background()
//			mid = int64(27515308)
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			p, err := s.RenewPanel(c, mid)
//			ctx.Convey("Then err should be nil.p should not be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//				ctx.So(p, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}

//func TestServiceRenewVip(t *testing.T) {
//	convey.Convey("RenewVip", t, func(ctx convey.C) {
//		var (
//			c   = context.Background()
//			mid = int64(0)
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			err := s.RenewVip(c, mid)
//			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//			})
//		})
//	})
//}
