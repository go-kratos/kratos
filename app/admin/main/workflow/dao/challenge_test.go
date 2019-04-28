package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/workflow/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoChall(t *testing.T) {
	convey.Convey("Chall", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			cid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			chall, err := d.Chall(c, cid)
			ctx.Convey("Then err should be nil.chall should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(chall, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoChalls(t *testing.T) {
	convey.Convey("Challs", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			cids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			challs, err := d.Challs(c, cids)
			ctx.Convey("Then err should be nil.challs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(challs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoStateChalls(t *testing.T) {
	convey.Convey("StateChalls", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			cids  = []int64{1}
			state = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			challs, err := d.StateChalls(c, cids, state)
			ctx.Convey("Then err should be nil.challs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(challs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLastChallIDsByGids(t *testing.T) {
	convey.Convey("LastChallIDsByGids", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			gids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cids, err := d.LastChallIDsByGids(c, gids)
			ctx.Convey("Then err should be nil.cids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpChall(t *testing.T) {
	convey.Convey("TxUpChall", t, func(ctx convey.C) {
		var (
			tx    = d.ORM.Begin()
			chall = &model.Chall{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.TxUpChall(tx, chall)
			err1 := tx.Commit().Error
			defer func() {
				if err != nil {
					tx.Rollback()
				}
			}()
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(err1, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxBatchUpChallByIDs(t *testing.T) {
	convey.Convey("TxBatchUpChallByIDs", t, func(ctx convey.C) {
		var (
			tx    = d.ORM.Begin()
			cids  = []int64{1}
			state = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxBatchUpChallByIDs(tx, cids, state)
			err1 := tx.Commit().Error
			defer func() {
				if err != nil {
					tx.Rollback()
				}
			}()
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(err1, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAttPathsByCids(t *testing.T) {
	convey.Convey("AttPathsByCids", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			cids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			paths, err := d.AttPathsByCids(c, cids)
			ctx.Convey("Then err should be nil.paths should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(paths, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAttPathsByCid(t *testing.T) {
	convey.Convey("AttPathsByCid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			cid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			paths, err := d.AttPathsByCid(c, cid)
			ctx.Convey("Then err should be nil.paths should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(paths, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpChallBusState(t *testing.T) {
	convey.Convey("UpChallBusState", t, func(ctx convey.C) {
		var (
			c               = context.Background()
			cid             = int64(1)
			busState        = int8(1)
			assigneeAdminid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpChallBusState(c, cid, busState, assigneeAdminid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBatchUpChallBusState(t *testing.T) {
	convey.Convey("BatchUpChallBusState", t, func(ctx convey.C) {
		var (
			c               = context.Background()
			cids            = []int64{1}
			busState        = int8(1)
			assigneeAdminid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.BatchUpChallBusState(c, cids, busState, assigneeAdminid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxChallsByBusStates(t *testing.T) {
	convey.Convey("TxChallsByBusStates", t, func(ctx convey.C) {
		var (
			tx        = d.ORM.Begin()
			business  = int8(1)
			oid       = int64(1)
			busStates = []int8{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cids, err := d.TxChallsByBusStates(tx, business, oid, busStates)
			err1 := tx.Commit().Error
			defer func() {
				if err != nil {
					tx.Rollback()
				}
			}()
			ctx.Convey("Then err should be nil.cids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(err1, convey.ShouldBeNil)
				ctx.So(cids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpChallsBusStateByIDs(t *testing.T) {
	convey.Convey("TxUpChallsBusStateByIDs", t, func(ctx convey.C) {
		var (
			tx              = d.ORM.Begin()
			cids            = []int64{1}
			busState        = int8(1)
			assigneeAdminid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxUpChallsBusStateByIDs(tx, cids, busState, assigneeAdminid)
			err1 := tx.Commit().Error
			defer func() {
				if err != nil {
					tx.Rollback()
				}
			}()
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(err1, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxUpChallExtraV2(t *testing.T) {
	convey.Convey("TxUpChallExtraV2", t, func(ctx convey.C) {
		var (
			tx       = d.ORM.Begin()
			business = int8(1)
			oid      = int64(1)
			adminid  = int64(1)
			extra    map[string]interface{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.TxUpChallExtraV2(tx, business, oid, adminid, extra)
			err1 := tx.Commit().Error
			defer func() {
				if err != nil {
					tx.Rollback()
				}
			}()
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(err1, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpExtraV3(t *testing.T) {
	convey.Convey("UpExtraV3", t, func(ctx convey.C) {
		var (
			gids    = []int64{1}
			adminid = int64(1)
			extra   = "test extra"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpExtraV3(gids, adminid, extra)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxUpChallTag(t *testing.T) {
	convey.Convey("TxUpChallTag", t, func(ctx convey.C) {
		var (
			tx  = d.ORM.Begin()
			cid = int64(1)
			tid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxUpChallTag(tx, cid, tid)
			err1 := tx.Commit().Error
			defer func() {
				if err != nil {
					tx.Rollback()
				}
			}()
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(err1, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBatchUpChallByIDs(t *testing.T) {
	convey.Convey("BatchUpChallByIDs", t, func(ctx convey.C) {
		var (
			cids          = []int64{1}
			dispatchState = uint32(1)
			adminid       = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.BatchUpChallByIDs(cids, dispatchState, adminid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBatchResetAssigneeAdminID(t *testing.T) {
	convey.Convey("BatchResetAssigneeAdminID", t, func(ctx convey.C) {
		var (
			cids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.BatchResetAssigneeAdminID(cids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxUpChallAssignee(t *testing.T) {
	convey.Convey("TxUpChallAssignee", t, func(ctx convey.C) {
		var (
			tx   = d.ORM.Begin()
			cids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxUpChallAssignee(tx, cids)
			err1 := tx.Commit().Error
			defer func() {
				if err != nil {
					tx.Rollback()
				}
			}()
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(err1, convey.ShouldBeNil)
			})
		})
	})
}
