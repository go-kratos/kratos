package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDelAvRatio(t *testing.T) {
	convey.Convey("DelAvRatio", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			limit = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelAvRatio(c, limit)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelIncome(t *testing.T) {
	convey.Convey("DelIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			limit = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelIncome(c, limit)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelActivity(t *testing.T) {
	convey.Convey("DelActivity", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			limit = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelActivity(c, limit)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCommonTagInfo(t *testing.T) {
	convey.Convey("CommonTagInfo", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			startAt = time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO tag_info(category_id,ratio,start_at,end_at,is_common,activity_id) VALUES(1,2,'2018-06-24','2018-06-24',1,0) ON DUPLICATE KEY UPDATE category_id=VALUES(category_id),ratio=VALUES(ratio)")
			tagInfos, err := d.CommonTagInfo(c, startAt)
			ctx.Convey("Then err should be nil.tagInfos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tagInfos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoNoCommonTagInfo(t *testing.T) {
	convey.Convey("NoCommonTagInfo", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			startAt = time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO tag_info(id,category_id,ratio,start_at,end_at,is_common,activity_id) VALUES(3,3,3,'2018-06-24','2018-06-24',0,0) ON DUPLICATE KEY UPDATE start_at=VALUES(start_at),end_at=VALUES(end_at),is_common=VALUES(is_common)")
			Exec(c, "INSERT INTO tag_up_info(tag_id,mid) VALUES(3,3) ON DUPLICATE KEY UPDATE tag_id=VALUES(tag_id)")
			tagInfos, err := d.NoCommonTagInfo(c, startAt)
			ctx.Convey("Then err should be nil.tagInfos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tagInfos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAIDsByMID(t *testing.T) {
	convey.Convey("AIDsByMID", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(12)
			tagID  = int(2)
			offset = int64(0)
			limit  = int64(100)
			month  = "06"
			date   = time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_daily_charge_06(av_id,inc_charge,tag_id,mid,date) VALUES(1,200,2,12,'2018-06-24') ON DUPLICATE KEY UPDATE av_id=VALUES(av_id), date=VALUES(date)")
			aids, err := d.AIDsByMID(c, mid, tagID, offset, limit, month, date)
			ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(aids, convey.ShouldBeNil)
			})
		})
	})
}

//	_commonAvSQL         = "SELECT id,av_id,inc_charge,is_deleted FROM av_daily_charge_%s WHERE tag_id = ? AND date = ? LIMIT ?,?"
func TestDaoAIDs(t *testing.T) {
	convey.Convey("AIDs", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			tagID  = int(10)
			offset = int64(0)
			limit  = int64(10)
			month  = "10"
			date   = time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_daily_charge_10(av_id,inc_charge,tag_id,mid,date) VALUES(1,200,10,1,'2018-06-24') ON DUPLICATE KEY UPDATE av_id=VALUES(av_id), date=VALUES(date),tag_id=VALUES(tag_id)")
			aids, err := d.AIDs(c, tagID, offset, limit, month, date)
			ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(aids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertRatio(t *testing.T) {
	convey.Convey("InsertRatio", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1,2,3)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertRatio(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpdateTagInfo(t *testing.T) {
	convey.Convey("TxUpdateTagInfo", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			tagID = int64(1)
			query = "total_income=1"
		)
		defer tx.Commit()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.TxUpdateTagInfo(tx, tagID, query)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpdateTagUpInfo(t *testing.T) {
	convey.Convey("TxUpdateTagUpInfo", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			query = "(1,2,3)"
		)
		defer tx.Commit()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.TxUpdateTagUpInfo(tx, query)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoActivityTagInfo(t *testing.T) {
	convey.Convey("ActivityTagInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			infos, err := d.ActivityTagInfo(c, date)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetActivityAVInfo(t *testing.T) {
	convey.Convey("GetActivityAVInfo", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			pn         = int(1)
			ps         = int(2)
			activities = []int64{3}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			info, total, err := d.GetActivityAVInfo(c, pn, ps, activities)
			ctx.Convey("Then err should not be nil.info,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoActivityMIDExist(t *testing.T) {
	convey.Convey("ActivityMIDExist", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tagID = int64(2)
			mid   = int64(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			info, err := d.ActivityMIDExist(c, tagID, mid)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAllTypes(t *testing.T) {
	convey.Convey("GetAllTypes", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rmap, err := d.GetAllTypes(c)
			ctx.Convey("Then err should not be nil.rmap should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rmap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetTagTotalIncome(t *testing.T) {
	convey.Convey("GetTagTotalIncome", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			tagIDs = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			infos, err := d.GetTagTotalIncome(c, tagIDs)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldBeNil)
			})
		})
	})
}
