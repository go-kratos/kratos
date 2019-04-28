package service

import (
	"context"
	"reflect"
	"testing"

	"go-common/app/service/main/point/dao"
	"go-common/app/service/main/point/model"
	xsql "go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestServicePointInfo(t *testing.T) {
	convey.Convey("PointInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
			pm  = &model.PointInfo{
				Mid:          1,
				PointBalance: 2,
				Ver:          1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "PointInfoCache", func(_ *dao.Dao, _ context.Context, _ int64) (*model.PointInfo, error) {
				return nil, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "PointInfo", func(_ *dao.Dao, _ context.Context, _ int64) (*model.PointInfo, error) {
				return pm, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "SetPointInfoCache", func(_ *dao.Dao, _ context.Context, _ *model.PointInfo) error {
				return nil
			})
			pi, err := s.PointInfo(c, mid)
			ctx.Convey("Then err should be nil.pi should not be nil.", func(ctx convey.C) {
				t.Logf("pi:%+v", pi)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(pi, convey.ShouldNotBeNil)
			})
			ctx.Reset(func() {
				monkey.UnpatchAll()
			})
		})
	})
}

func TestServicePointHistory(t *testing.T) {
	convey.Convey("PointHistory", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(4780461)
			cursor = int(0)
			ps     = int(20)
			total  = int(2)
			pms    []*model.PointHistory
			pm     = &model.PointHistory{
				ID:           13,
				Mid:          4780461,
				Point:        60,
				ChangeType:   1,
				PointBalance: 418,
			}
		)
		pms = append(pms, pm)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "PointHistoryCount", func(_ *dao.Dao, _ context.Context, _ int64) (int, error) {
				return total, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "PointHistory", func(_ *dao.Dao, _ context.Context, _ int64, _, _ int) ([]*model.PointHistory, error) {
				return pms, nil
			})
			phs, total, ncursor, err := s.PointHistory(c, mid, cursor, ps)
			ctx.Convey("Then err should be nil.phs,total,ncursor should not be nil.", func(ctx convey.C) {
				t.Logf("phs:%+v", phs)
				t.Logf("total:%+v", total)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ncursor, convey.ShouldNotBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(phs, convey.ShouldNotBeNil)
			})
			ctx.Reset(func() {
				monkey.UnpatchAll()
			})
		})
	})
}

func TestServiceOldPointHistory(t *testing.T) {
	convey.Convey("OldPointHistory", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(4780461)
			pn    = int(1)
			ps    = int(20)
			total = int(2)
			pms   []*model.OldPointHistory
			pm    = &model.OldPointHistory{
				ID:           13,
				Mid:          4780461,
				Point:        60,
				ChangeType:   1,
				PointBalance: 418,
			}
		)
		pms = append(pms, pm)

		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "PointHistoryCount", func(_ *dao.Dao, _ context.Context, _ int64) (int, error) {
				return total, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "OldPointHistory", func(_ *dao.Dao, _ context.Context, _ int64, _, _ int) ([]*model.OldPointHistory, error) {
				return pms, nil
			})
			phs, total, err := s.OldPointHistory(c, mid, pn, ps)
			ctx.Convey("Then err should be nil.phs,total should not be nil.", func(ctx convey.C) {
				t.Logf("phs:%+v", phs)
				t.Logf("total:%+v", total)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(phs, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}

func TestServicePointAddByBp(t *testing.T) {
	convey.Convey("PointAddByBp", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			pa = &model.ArgPointAdd{
				Mid:        1,
				ChangeType: 3,
				RelationID: "121",
				Bcoin:      10,
				Remark:     "test",
				OrderID:    "31",
			}
			hid = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "ExistPointOrder", func(_ *dao.Dao, _ context.Context, _ string) (int, error) {
				return hid, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "DelPointInfoCache", func(_ *dao.Dao, _ context.Context, _ int64) error {
				return nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "DelPointInfoCache", func(_ *dao.Dao, _ context.Context, _ int64) error {
				return nil
			})
			p, err := s.PointAddByBp(c, pa)
			ctx.Convey("Then err should be nil.p should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceupdatePointWithHistory(t *testing.T) {
	convey.Convey("updatePointWithHistory", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ph = &model.PointHistory{
				ID:           13,
				Mid:          4780461,
				Point:        60,
				ChangeType:   1,
				PointBalance: 418,
			}
			tx = &xsql.Tx{}
			pm = &model.PointInfo{
				Mid:          1,
				PointBalance: 2,
				Ver:          1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "BeginTran", func(_ *dao.Dao, _ context.Context) (*xsql.Tx, error) {
				return nil, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Commit", func(_ *xsql.Tx) error {
				return nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "InsertPointHistory", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ *model.PointHistory) (int64, error) {
				return 0, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "TxPointInfo", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ int64) (*model.PointInfo, error) {
				return pm, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "InsertPoint", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ *model.PointInfo) (int64, error) {
				return 1, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "UpdatePointInfo", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ *model.PointInfo, _ int64) (int64, error) {
				return 1, nil
			})

			pointBalance, activePoint, err := s.updatePointWithHistory(c, ph)
			ctx.Convey("Then err should be nil.pointBalance,activePoint should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(activePoint, convey.ShouldNotBeNil)
				ctx.So(pointBalance, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}

func TestServiceupdatePoint(t *testing.T) {
	convey.Convey("updatePoint", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx    = &xsql.Tx{}
			mid   = int64(1)
			point = int64(10)
			pm    = &model.PointInfo{
				Mid:          1,
				PointBalance: 2,
				Ver:          1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "TxPointInfo", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ int64) (*model.PointInfo, error) {
				return pm, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "InsertPoint", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ *model.PointInfo) (int64, error) {
				return 1, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "UpdatePointInfo", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ *model.PointInfo, _ int64) (int64, error) {
				return 1, nil
			})

			pb, err := s.updatePoint(c, tx, mid, point)
			ctx.Convey("Then err should be nil.pb should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(pb, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}

func TestServiceConsumePoint(t *testing.T) {
	convey.Convey("ConsumePoint", t, func(ctx convey.C) {
		var (
			tx = &xsql.Tx{}
			c  = context.Background()
			pc = &model.ArgPointConsume{
				Mid:        2,
				ChangeType: 4,
				RelationID: "1",
				Point:      1,
				Remark:     "测试消费",
			}
			pm = &model.PointInfo{
				Mid:          2,
				PointBalance: 2000,
				Ver:          1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "BeginTran", func(_ *dao.Dao, _ context.Context) (*xsql.Tx, error) {
				return nil, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Commit", func(_ *xsql.Tx) error {
				return nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Rollback", func(_ *xsql.Tx) error {
				return nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "InsertPointHistory", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ *model.PointHistory) (int64, error) {
				return 0, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "TxPointInfo", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ int64) (*model.PointInfo, error) {
				return pm, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "InsertPoint", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ *model.PointInfo) (int64, error) {
				return 1, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "UpdatePointInfo", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ *model.PointInfo, _ int64) (int64, error) {
				return 1, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "DelPointInfoCache", func(_ *dao.Dao, _ context.Context, _ int64) error {
				return nil
			})
			status, err := s.ConsumePoint(c, pc)
			ctx.Convey("Then err should be nil.status should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(status, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}

func TestServiceAddPoint(t *testing.T) {
	convey.Convey("AddPoint", t, func(ctx convey.C) {
		var (
			tx = &xsql.Tx{}
			c  = context.Background()
			pc = &model.ArgPoint{
				Mid:        52,
				ChangeType: 3,
				Remark:     "测试增加",
				Point:      300,
				Operator:   "yubaihai",
			}
			pm = &model.PointInfo{
				Mid:          52,
				PointBalance: 2,
				Ver:          1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "BeginTran", func(_ *dao.Dao, _ context.Context) (*xsql.Tx, error) {
				return nil, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Commit", func(_ *xsql.Tx) error {
				return nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "InsertPointHistory", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ *model.PointHistory) (int64, error) {
				return 0, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "TxPointInfo", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ int64) (*model.PointInfo, error) {
				return pm, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "InsertPoint", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ *model.PointInfo) (int64, error) {
				return 1, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "UpdatePointInfo", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, _ *model.PointInfo, _ int64) (int64, error) {
				return 1, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "DelPointInfoCache", func(_ *dao.Dao, _ context.Context, _ int64) error {
				return nil
			})
			status, err := s.AddPoint(c, pc)
			ctx.Convey("Then err should be nil.status should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(status, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}
