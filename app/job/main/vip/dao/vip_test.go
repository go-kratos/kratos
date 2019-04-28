package dao

import (
	"context"
	"go-common/app/job/main/vip/model"
	xtime "go-common/library/time"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddChangeHistoryBatch(t *testing.T) {
	convey.Convey("AddChangeHistoryBatch", t, func(ctx convey.C) {
		var (
			res = []*model.VipChangeHistory{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddChangeHistoryBatch(res)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelChangeHistory(t *testing.T) {
	convey.Convey("SelChangeHistory", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			startID = int64(0)
			endID   = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelChangeHistory(c, startID, endID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelOldChangeHistory(t *testing.T) {
	convey.Convey("SelOldChangeHistory", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			startID = int64(0)
			endID   = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelOldChangeHistory(c, startID, endID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelOldChangeHistoryMaxID(t *testing.T) {
	convey.Convey("SelOldChangeHistoryMaxID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelOldChangeHistoryMaxID(c)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelChangeHistoryMaxID(t *testing.T) {
	convey.Convey("SelChangeHistoryMaxID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelChangeHistoryMaxID(c)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoVipStatus(t *testing.T) {
	convey.Convey("VipStatus", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.VipStatus(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelMaxID(t *testing.T) {
	convey.Convey("SelMaxID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelMaxID(c)
			ctx.Convey("Then err should be nil.mID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateVipUserInfo(t *testing.T) {
	convey.Convey("UpdateVipUserInfo", t, func(ctx convey.C) {
		var (
			r = &model.VipUserInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.StartTx(context.Background())
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tx, convey.ShouldNotBeNil)
			_, err = d.UpdateVipUserInfo(tx, r)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateVipUser(t *testing.T) {
	convey.Convey("UpdateVipUser", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(0)
			status  = int8(0)
			vipType = int8(0)
			payType = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UpdateVipUser(c, mid, status, vipType, payType)
			ctx.Convey("Then err should be nil.eff should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDupUserDiscountHistory(t *testing.T) {
	convey.Convey("DupUserDiscountHistory", t, func(ctx convey.C) {
		var (
			r = &model.VipUserDiscountHistory{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.StartTx(context.Background())
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tx, convey.ShouldNotBeNil)
			_, err = d.DupUserDiscountHistory(tx, r)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddUserInfo(t *testing.T) {
	convey.Convey("AddUserInfo", t, func(ctx convey.C) {
		var (
			r = &model.VipUserInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.StartTx(context.Background())
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tx, convey.ShouldNotBeNil)
			_, err = d.AddUserInfo(tx, r)
			ctx.Convey("Then err should be nil.eff should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelVipUserInfo(t *testing.T) {
	convey.Convey("SelVipUserInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelVipUserInfo(c, mid)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateUserInfo(t *testing.T) {
	convey.Convey("UpdateUserInfo", t, func(ctx convey.C) {
		var (
			r = &model.VipUserInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.StartTx(context.Background())
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tx, convey.ShouldNotBeNil)
			_, err = d.UpdateUserInfo(tx, r)
			ctx.Convey("Then err should be nil.eff should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelUserDiscountMaxID(t *testing.T) {
	convey.Convey("SelUserDiscountMaxID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelUserDiscountMaxID(c)
			ctx.Convey("Then err should be nil.maxID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelUserDiscountHistorys(t *testing.T) {
	convey.Convey("SelUserDiscountHistorys", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			sID        = int(0)
			eID        = int(0)
			discountID = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelUserDiscountHistorys(c, sID, eID, discountID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelOldUserInfoMaxID(t *testing.T) {
	convey.Convey("SelOldUserInfoMaxID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelOldUserInfoMaxID(c)
			ctx.Convey("Then err should be nil.maxID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelUserInfoMaxID(t *testing.T) {
	convey.Convey("SelUserInfoMaxID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelUserInfoMaxID(c)
			ctx.Convey("Then err should be nil.maxID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelUserInfos(t *testing.T) {
	convey.Convey("SelUserInfos", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sID = int(0)
			eID = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelUserInfos(c, sID, eID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelVipList(t *testing.T) {
	convey.Convey("SelVipList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int(0)
			endID = int(0)
			ot    = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelVipList(c, id, endID, ot)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelVipUsers(t *testing.T) {
	convey.Convey("SelVipUsers", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int(0)
			endID = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelVipUsers(c, id, endID, xtime.Time(time.Now().Unix()), xtime.Time(time.Now().Unix()))
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelVipUserInfos(t *testing.T) {
	convey.Convey("SelVipUserInfos", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			id     = int(0)
			endID  = int(0)
			status = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelVipUserInfos(c, id, endID, xtime.Time(time.Now().Unix()), xtime.Time(time.Now().Unix()), status)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoOldVipInfo(t *testing.T) {
	convey.Convey("OldVipInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.OldVipInfo(c, mid)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelEffectiveScopeVipList(t *testing.T) {
	convey.Convey("SelEffectiveScopeVipList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int(0)
			endID = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelEffectiveScopeVipList(c, id, endID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelOldUserInfoMaps(t *testing.T) {
	convey.Convey("SelOldUserInfoMaps", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sID = int(0)
			eID = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelOldUserInfoMaps(c, sID, eID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func Test_SelChangeHistoryMaps(t *testing.T) {
	convey.Convey("SelChangeHistoryMaps", t, func(ctx convey.C) {
		_, err := d.SelChangeHistoryMaps(context.Background(), []int64{7593623})
		ctx.So(err, convey.ShouldBeNil)
	})
}

func Test_SelOldChangeHistoryMaps(t *testing.T) {
	convey.Convey("SelOldChangeHistoryMaps", t, func(ctx convey.C) {
		_, err := d.SelOldChangeHistoryMaps(context.Background(), []int64{7593623})
		ctx.So(err, convey.ShouldBeNil)
	})
}

func Test_SelVipByIds(t *testing.T) {
	convey.Convey("SelVipByIds", t, func(ctx convey.C) {
		_, err := d.SelVipByIds(context.Background(), []int64{7593623})
		ctx.So(err, convey.ShouldBeNil)
	})
}
