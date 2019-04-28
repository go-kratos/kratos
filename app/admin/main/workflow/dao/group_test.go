package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/workflow/model/param"

	"github.com/jinzhu/gorm"
	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGroupByID(t *testing.T) {
	convey.Convey("GroupByID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			gid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			g, err := d.GroupByID(c, gid)
			ctx.Convey("Then err should be nil.g should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(g, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGroupByOid(t *testing.T) {
	convey.Convey("GroupByOid", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			oid      = int64(1)
			business = int8(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			g, err := d.GroupByOid(c, oid, business)
			ctx.Convey("Then err should be nil.g should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(g, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxGroupsByOidsStates(t *testing.T) {
	convey.Convey("TxGroupsByOidsStates", t, func(ctx convey.C) {
		var (
			tx       = &gorm.DB{}
			oids     = []int64{}
			business = int8(0)
			state    = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			groups, err := d.TxGroupsByOidsStates(tx, oids, business, state)
			ctx.Convey("Then err should be nil.groups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(groups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGroups(t *testing.T) {
	convey.Convey("Groups", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			gids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			groups, err := d.Groups(c, gids)
			ctx.Convey("Then err should be nil.groups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(groups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxGroups(t *testing.T) {
	convey.Convey("TxGroups", t, func(ctx convey.C) {
		var (
			tx   = d.ORM.Begin()
			gids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			groups, err := d.TxGroups(tx, gids)
			err1 := tx.Commit().Error
			defer func() {
				if err != nil {
					tx.Rollback()
				}
			}()
			ctx.Convey("Then err should be nil.groups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(err1, convey.ShouldBeNil)
				ctx.So(groups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpGroup(t *testing.T) {
	convey.Convey("TxUpGroup", t, func(ctx convey.C) {
		var (
			tx       = d.ORM.Begin()
			oid      = int64(1)
			business = int8(1)
			tid      = int64(0)
			note     = "test note"
			rid      = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxUpGroup(tx, oid, business, tid, note, rid)
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

func TestDaoUpGroupRole(t *testing.T) {
	convey.Convey("UpGroupRole", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			grsp = &param.GroupRoleSetParam{
				GID:       []int64{1},
				AdminID:   1,
				AdminName: "anonymous",
				BID:       1,
				RID:       1,
				TID:       1,
				Note:      "test note",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpGroupRole(c, grsp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxUpGroupState(t *testing.T) {
	convey.Convey("TxUpGroupState", t, func(ctx convey.C) {
		var (
			tx    = d.ORM.Begin()
			gid   = int64(1)
			state = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxUpGroupState(tx, gid, state)
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

func TestDaoTxUpGroupHandling(t *testing.T) {
	convey.Convey("TxUpGroupHandling", t, func(ctx convey.C) {
		var (
			tx       = d.ORM.Begin()
			gid      = int64(1)
			handling = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxUpGroupHandling(tx, gid, handling)
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

func TestDaoTxBatchUpGroupHandling(t *testing.T) {
	convey.Convey("TxBatchUpGroupHandling", t, func(ctx convey.C) {
		var (
			tx       = d.ORM.Begin()
			gids     = []int64{1, 2}
			handling = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxBatchUpGroupHandling(tx, gids, handling)
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

func TestDaoTxBatchUpGroupState(t *testing.T) {
	convey.Convey("TxBatchUpGroupState", t, func(ctx convey.C) {
		var (
			tx    = d.ORM.Begin()
			gids  = []int64{1, 2}
			state = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxBatchUpGroupState(tx, gids, state)
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

func TestDaoTxSetGroupStateTid(t *testing.T) {
	convey.Convey("TxSetGroupStateTid", t, func(ctx convey.C) {
		var (
			tx    = d.ORM.Begin()
			gids  = []int64{1, 2}
			state = int8(0)
			tid   = int64(0)
			rid   = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxSetGroupStateTid(tx, gids, state, rid, tid)
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

func TestDaoTxSimpleSetGroupState(t *testing.T) {
	convey.Convey("TxSimpleSetGroupState", t, func(ctx convey.C) {
		var (
			tx    = d.ORM.Begin()
			gids  = []int64{1, 2}
			state = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxSimpleSetGroupState(tx, gids, state)
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
