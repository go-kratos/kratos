package dao

import (
	"context"
	"go-common/app/admin/main/growup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetActivityByName(t *testing.T) {
	convey.Convey("GetActivityByName", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			name = "test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO creative_activity(name) VALUES('test')")
			_, err := d.GetActivityByName(c, name)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoActivityCount(t *testing.T) {
	convey.Convey("ActivityCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.ActivityCount(c, query)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetActivities(t *testing.T) {
	convey.Convey("GetActivities", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = ""
			from  = int(0)
			limit = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			acs, err := d.GetActivities(c, query, from, limit)
			ctx.Convey("Then err should be nil.acs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(acs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxGetActivityByName(t *testing.T) {
	convey.Convey("TxGetActivityByName", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			name  = "test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			id, err := d.TxGetActivityByName(tx, name)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxInsertActivity(t *testing.T) {
	convey.Convey("TxInsertActivity", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			ac    = &model.CActivity{
				Name: "test1",
			}
			update = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			Exec(context.Background(), "DELETE FROM creative_activity WHERE name = 'test1'")
			rows, err := d.TxInsertActivity(tx, ac, update)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetActivityBonus(t *testing.T) {
	convey.Convey("GetActivityBonus", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{111}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO activity_bonus(activity_id,ranking,bonus_money) VALUES(111, 1, 1)")
			_, err := d.GetActivityBonus(c, ids)
			ctx.Convey("Then err should be nil.brs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxInsertActivityBonusBatch(t *testing.T) {
	convey.Convey("TxInsertActivityBonusBatch", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			vals  = "(1112, 10, 10)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			rows, err := d.TxInsertActivityBonusBatch(tx, vals)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpActivityStateCount(t *testing.T) {
	convey.Convey("UpActivityStateCount", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			id     = int64(111)
			states = []int64{1, 2, 3}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_activity(activity_id, mid, state) VALUES(111, 1001, 2) ON DUPLICATE KEY UPDATE state = 2")
			count, err := d.UpActivityStateCount(c, id, states)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListUpActivity(t *testing.T) {
	convey.Convey("ListUpActivity", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(111)
			from  = int(0)
			limit = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.ListUpActivity(c, id, from, limit)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListUpActivitySuccess(t *testing.T) {
	convey.Convey("ListUpActivitySuccess", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(111)
			mid   = int64(1001)
			from  = int(0)
			limit = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.ListUpActivitySuccess(c, id, mid, from, limit)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateUpActivityState(t *testing.T) {
	convey.Convey("UpdateUpActivityState", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			activityID = int64(111)
			mids       = []int64{1001}
			oldState   = int(2)
			newState   = int(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_activity(activity_id, mid, state) VALUES(111, 1001, 2) ON DUPLICATE KEY UPDATE state = 2")
			rows, err := d.UpdateUpActivityState(c, activityID, mids, oldState, newState)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpdateUpActivityState(t *testing.T) {
	convey.Convey("TxUpdateUpActivityState", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			tx, _      = d.BeginTran(c)
			activityID = int64(111)
			mids       = []int64{1001}
			oldState   = int(2)
			newState   = int(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			Exec(c, "INSERT INTO up_activity(activity_id, mid, state) VALUES(111, 1001, 2) ON DUPLICATE KEY UPDATE state = 2")
			rows, err := d.TxUpdateUpActivityState(tx, activityID, mids, oldState, newState)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
