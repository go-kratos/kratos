package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTxUpdateUpSpyState(t *testing.T) {
	convey.Convey("TxUpdateUpSpyState", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			state = int(3)
			mid   = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			d.Exec(context.Background(), "INSERT INTO up_spy_statistics(mid, account_state) VALUES(1001, 4) ON DUPLICATE KEY UPDATE account_state = 4")
			rows, err := d.TxUpdateUpSpyState(tx, state, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpdateAvSpyState(t *testing.T) {
	convey.Convey("TxUpdateAvSpyState", t, func(ctx convey.C) {
		var (
			tx, _    = d.BeginTran(context.Background())
			state    = int(5)
			archives = []int64{1000}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			d.Exec(context.Background(), "INSERT INTO archive_spy_statistics(archive_id, deducted) VALUES(1000, 5)")
			rows, err := d.TxUpdateAvSpyState(tx, state, archives)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpSpyCount(t *testing.T) {
	convey.Convey("UpSpyCount", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.UpSpyCount(c)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpSpies(t *testing.T) {
	convey.Convey("UpSpies", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			query  = "WHERE"
			offset = int(0)
			limit  = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UpSpies(c, query, offset, limit)
			ctx.Convey("Then err should be nil.spies should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoArchiveSpyCount(t *testing.T) {
	convey.Convey("ArchiveSpyCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.ArchiveSpyCount(c, query)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoArchiveSpies(t *testing.T) {
	convey.Convey("ArchiveSpies", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			query  = ""
			offset = int(0)
			limit  = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			spies, err := d.ArchiveSpies(c, query, offset, limit)
			ctx.Convey("Then err should be nil.spies should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(spies, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCheatFansCount(t *testing.T) {
	convey.Convey("CheatFansCount", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.CheatFansCount(c)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCheatFans(t *testing.T) {
	convey.Convey("CheatFans", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			from  = int64(0)
			limit = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			fans, err := d.CheatFans(c, from, limit)
			ctx.Convey("Then err should be nil.fans should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(fans, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCheatUp(t *testing.T) {
	convey.Convey("DelCheatUp", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1000)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Exec(c, "INSERT INTO cheat_fans_info(mid, is_deleted) values(1000, 0) ON DUPLICATE KEY UPDATE is_deleted = 0")
			rows, err := d.DelCheatUp(c, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertCheatFansInfo(t *testing.T) {
	convey.Convey("InsertCheatFansInfo", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(123, 'tt', '2018-09-01', 100, 100, '2018-09-01')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertCheatFansInfo(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpRealFansCount(t *testing.T) {
	convey.Convey("GetUpRealFansCount", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			host = ""
			mid  = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GetUpRealFansCount(c, host, mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpCheatFansCount(t *testing.T) {
	convey.Convey("GetUpCheatFansCount", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			host = ""
			mid  = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GetUpCheatFansCount(c, host, mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
