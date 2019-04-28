package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/growup/model"
	xtime "go-common/library/time"
)

func TestDaoGetAccountState(t *testing.T) {
	convey.Convey("GetAccountState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "up_info_video"
			mid   = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state, is_deleted) VALUES(1001, 3, 0) ON DUPLICATE KEY UPDATE account_state = 3, is_deleted = 0")
			state, err := d.GetAccountState(c, table, mid)
			ctx.Convey("Then err should be nil.state should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(state, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpSignedAt(t *testing.T) {
	convey.Convey("GetUpSignedAt", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state, is_deleted) VALUES(1001, 3, 0) ON DUPLICATE KEY UPDATE account_state = 3, is_deleted = 0")
			signedAt, err := d.GetUpSignedAt(c, "up_info_video", mid)
			ctx.Convey("Then err should be nil.signedAt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(signedAt, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertUpInfo(t *testing.T) {
	convey.Convey("InsertUpInfo", t, func(ctx convey.C) {
		var (
			c               = context.Background()
			table           = "up_info_video"
			totalCountField = "total_play_count"
			v               = &model.UpInfo{
				MID:          1002,
				AccountState: 3,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "DELETE FROM up_info_video WHERE mid = 1002")
			rows, err := d.InsertUpInfo(c, table, totalCountField, v)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxInsertBgmUpInfo(t *testing.T) {
	convey.Convey("TxInsertBgmUpInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			v     = &model.UpInfo{
				MID:          1002,
				AccountState: 3,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			Exec(c, "DELETE FROM up_info_bgm WHERE mid = 1002")
			rows, err := d.TxInsertBgmUpInfo(tx, v)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxInsertCreditScore(t *testing.T) {
	convey.Convey("TxInsertCreditScore", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			mid   = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			rows, err := d.TxInsertCreditScore(tx, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBlocked(t *testing.T) {
	convey.Convey("Blocked", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_blocked(mid) VALUES(1001)")
			id, err := d.Blocked(c, mid)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoWhite(t *testing.T) {
	convey.Convey("White", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_white_list(mid) VALUES(1001)")
			m, err := d.White(c, mid)
			ctx.Convey("Then err should be nil.m should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(m, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAvUpStatus(t *testing.T) {
	convey.Convey("AvUpStatus", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state, is_deleted) VALUES(1001, 3, 0) ON DUPLICATE KEY UPDATE account_state = 3, is_deleted = 0")
			status, err := d.AvUpStatus(c, mid)
			ctx.Convey("Then err should be nil.status should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(status, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBgmUpStatus(t *testing.T) {
	convey.Convey("BgmUpStatus", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_bgm(mid, account_state, is_deleted) VALUES(1001, 3, 0) ON DUPLICATE KEY UPDATE account_state = 3, is_deleted = 0")
			status, err := d.BgmUpStatus(c, mid)
			ctx.Convey("Then err should be nil.status should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(status, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoColumnUpStatus(t *testing.T) {
	convey.Convey("ColumnUpStatus", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_column(mid, account_state, is_deleted) VALUES(1001, 3, 0) ON DUPLICATE KEY UPDATE account_state = 3, is_deleted = 0")
			status, err := d.ColumnUpStatus(c, mid)
			ctx.Convey("Then err should be nil.status should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(status, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCategoryInfo(t *testing.T) {
	convey.Convey("CategoryInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			nickname, categoryID, err := d.CategoryInfo(c, mid)
			ctx.Convey("Then err should be nil.nickname,categoryID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(categoryID, convey.ShouldNotBeNil)
				ctx.So(nickname, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoFans(t *testing.T) {
	convey.Convey("Fans", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_base_statistics(mid, fans) VALUES(1001, 100)")
			fans, err := d.Fans(c, mid)
			ctx.Convey("Then err should be nil.fans should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(fans, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxQuit(t *testing.T) {
	convey.Convey("TxQuit", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			tx, _     = d.BeginTran(c)
			table     = "up_info_video"
			mid       = int64(1003)
			quitAt    = xtime.Time(time.Now().Unix())
			expiredIn = xtime.Time(time.Now().Unix())
			reason    = "test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			Exec(c, "INSERT INTO up_info_video(mid, account_state, is_deleted) VALUES(1003, 3, 0) ON DUPLICATE KEY UPDATE account_state = 3, is_deleted = 0")
			rows, err := d.TxQuit(tx, table, mid, quitAt, expiredIn, reason)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxInsertCreditRecord(t *testing.T) {
	convey.Convey("TxInsertCreditRecord", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			cr    = &model.CreditRecord{
				MID:    1001,
				Reason: 1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			Exec(c, "DELETE FROM credit_score_record WHERE mid = 1001")
			rows, err := d.TxInsertCreditRecord(tx, cr)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoNickname(t *testing.T) {
	convey.Convey("Nickname", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1004)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state, is_deleted, nickname) VALUES(1004, 3, 0, 'test') ON DUPLICATE KEY UPDATE account_state = 3, is_deleted = 0")
			nickname, err := d.Nickname(c, mid)
			ctx.Convey("Then err should be nil.nickname should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(nickname, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCreditScore(t *testing.T) {
	convey.Convey("CreditScore", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO credit_score(mid, score) VALUES(1001, 100)")
			score, err := d.CreditScore(c, mid)
			ctx.Convey("Then err should be nil.score should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(score, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBgmUpCount(t *testing.T) {
	convey.Convey("BgmUpCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO background_music(mid) VALUES(1001)")
			count, err := d.BgmUpCount(c, mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxDeductCreditScore(t *testing.T) {
	convey.Convey("TxDeductCreditScore", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			score = int(10)
			mid   = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			Exec(c, "INSERT INTO credit_score(mid, score) VALUES(1001, 100)")
			rows, err := d.TxDeductCreditScore(tx, score, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
