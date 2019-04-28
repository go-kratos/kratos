package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/growup/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInsertWhitelist(t *testing.T) {
	convey.Convey("InsertWhitelist", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
			typ = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "DELETE FROM up_white_list where mid = 1001")
			rows, err := d.InsertWhitelist(c, mid, typ)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPendings(t *testing.T) {
	convey.Convey("Pendings", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mids  = []int64{1001}
			table = "up_info_video"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state) VALUES(1001, 2) ON DUPLICATE KEY UPDATE account_state = 2")
			ms, err := d.Pendings(c, mids, table)
			ctx.Convey("Then err should be nil.ms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ms, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUnusualUps(t *testing.T) {
	convey.Convey("UnusualUps", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mids  = []int64{1002}
			table = "up_info_video"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state) VALUES(1002, 5) ON DUPLICATE KEY UPDATE account_state = 5")
			ms, err := d.UnusualUps(c, mids, table)
			ctx.Convey("Then err should be nil.ms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ms, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertUpVideo(t *testing.T) {
	convey.Convey("InsertUpVideo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			now = xtime.Time(time.Now().Unix())
			v   = &model.UpInfo{
				MID:                  1003,
				Nickname:             "aa",
				AccountType:          1,
				AccountState:         1,
				OriginalArchiveCount: 1,
				MainCategory:         1,
				Fans:                 1,
				SignType:             1,
				Reason:               "",
				ApplyAt:              now,
				SignedAt:             now,
				RejectAt:             now,
				ForbidAt:             now,
				QuitAt:               now,
				DismissAt:            now,
				ExpiredIn:            now,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "DELETE FROM up_info_video WHERE mid = 1003")
			rows, err := d.InsertUpVideo(c, v)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertUpColumn(t *testing.T) {
	convey.Convey("InsertUpColumn", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			up = &model.UpInfo{
				MID:      1004,
				Nickname: "aa",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "DELETE FROM up_info_column WHERE mid = 1004")
			rows, err := d.InsertUpColumn(c, up)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertBgmUpInfo(t *testing.T) {
	convey.Convey("InsertBgmUpInfo", t, func(ctx convey.C) {
		var (
			c = context.Background()
			m = &model.UpInfo{
				MID:          1005,
				AccountState: 3,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "DELETE FROM up_info_bgm WHERE mid = 1005")
			rows, err := d.InsertBgmUpInfo(c, m)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCategoryInfo(t *testing.T) {
	convey.Convey("CategoryInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			nickname, categoryID, err := d.CategoryInfo(c, mid)
			Exec(c, "INSERT INTO up_category_info(mid,nick_name) VALUES(1001, 'tt')")
			ctx.Convey("Then err should be nil.nickname,categoryID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(categoryID, convey.ShouldNotBeNil)
				ctx.So(nickname, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoStat(t *testing.T) {
	convey.Convey("Stat", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_base_statistics(mid, fans, avs) VALUES(1001, 10, 10)")
			fans, avs, err := d.Stat(c, mid)
			ctx.Convey("Then err should be nil.fans,avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
				ctx.So(fans, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpsCount(t *testing.T) {
	convey.Convey("UpsCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "up_info_video"
			query = "id > 0"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state) VALUES(1006, 3)")
			count, err := d.UpsCount(c, table, query)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpsVideoInfo(t *testing.T) {
	convey.Convey("UpsVideoInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = "id > 0"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.UpsVideoInfo(c, query)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpsColumnInfo(t *testing.T) {
	convey.Convey("UpsColumnInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = "id > 0"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.UpsColumnInfo(c, query)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpsBgmInfo(t *testing.T) {
	convey.Convey("UpsBgmInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = "id > 0"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.UpsBgmInfo(c, query)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoReject(t *testing.T) {
	convey.Convey("Reject", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			table     = "up_info_video"
			state     = int(6)
			reason    = "test"
			rejectAt  = xtime.Time(time.Now().Unix())
			expiredIn = xtime.Time(time.Now().Unix())
			mids      = []int64{1007}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state) VALUES(1007, 3) ON DUPLICATE KEY UPDATE account_state = 3")
			rows, err := d.Reject(c, table, state, reason, rejectAt, expiredIn, mids)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPass(t *testing.T) {
	convey.Convey("Pass", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			table    = "up_info_video"
			state    = int(3)
			signedAt = xtime.Time(time.Now().Unix())
			mids     = []int64{1008}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state) VALUES(1008, 2) ON DUPLICATE KEY UPDATE account_state = 2")
			rows, err := d.Pass(c, table, state, signedAt, mids)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDismiss(t *testing.T) {
	convey.Convey("Dismiss", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			table     = "up_info_video"
			newState  = int(6)
			oldState  = int(3)
			reason    = "test"
			dismissAt = xtime.Time(time.Now().Unix())
			quitAt    = xtime.Time(time.Now().Unix())
			mid       = int64(1008)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state) VALUES(1008, 3) ON DUPLICATE KEY UPDATE account_state = 3")
			rows, err := d.Dismiss(c, table, newState, oldState, reason, dismissAt, quitAt, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxDismiss(t *testing.T) {
	convey.Convey("TxDismiss", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			tx, _     = d.BeginTran(c)
			table     = "up_info_video"
			newState  = int(6)
			oldState  = int(3)
			reason    = ""
			dismissAt = xtime.Time(time.Now().Unix())
			quitAt    = xtime.Time(time.Now().Unix())
			mid       = int64(1008)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			Exec(c, "INSERT INTO up_info_video(mid, account_state) VALUES(1008, 3) ON DUPLICATE KEY UPDATE account_state = 3")
			rows, err := d.TxDismiss(tx, table, newState, oldState, reason, dismissAt, quitAt, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoForbid(t *testing.T) {
	convey.Convey("Forbid", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			table     = "up_info_video"
			newState  = int(7)
			oldState  = int(3)
			reason    = ""
			forbidAt  = xtime.Time(time.Now().Unix())
			expiredIn = xtime.Time(time.Now().Unix())
			mid       = int64(1008)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state) VALUES(1008, 3) ON DUPLICATE KEY UPDATE account_state = 3")
			rows, err := d.Forbid(c, table, newState, oldState, reason, forbidAt, expiredIn, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxForbid(t *testing.T) {
	convey.Convey("TxForbid", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			tx, _     = d.BeginTran(c)
			table     = "up_info_video"
			newState  = int(6)
			oldState  = int(3)
			reason    = "test"
			forbidAt  = xtime.Time(time.Now().Unix())
			expiredIn = xtime.Time(time.Now().Unix())
			mid       = int64(1008)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			Exec(c, "INSERT INTO up_info_video(mid, account_state) VALUES(1008, 3) ON DUPLICATE KEY UPDATE account_state = 3")
			rows, err := d.TxForbid(tx, table, newState, oldState, reason, forbidAt, expiredIn, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateAccountState(t *testing.T) {
	convey.Convey("UpdateAccountState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "up_info_video"
			state = int(4)
			mid   = int64(1008)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state) VALUES(1008, 3) ON DUPLICATE KEY UPDATE account_state = 3")
			rows, err := d.UpdateAccountState(c, table, state, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelUpInfo(t *testing.T) {
	convey.Convey("DelUpInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "up_info_video"
			mid   = int64(1008)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, is_deleted) VALUES(1008, 0) ON DUPLICATE KEY UPDATE is_deleted = 0")
			rows, err := d.DelUpInfo(c, table, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRecUpInfo(t *testing.T) {
	convey.Convey("RecUpInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "up_info_video"
			mid   = int64(1008)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, is_deleted) VALUES(1008, 1) ON DUPLICATE KEY UPDATE is_deleted = 1")
			rows, err := d.RecUpInfo(c, table, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoupdateUpInfoDel(t *testing.T) {
	convey.Convey("updateUpInfoDel", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			table     = "up_info_video"
			mid       = int64(1008)
			isDeleted = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, is_deleted) VALUES(1008, 1) ON DUPLICATE KEY UPDATE is_deleted = 1")
			rows, err := d.updateUpInfoDel(c, table, mid, isDeleted)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelUpAccount(t *testing.T) {
	convey.Convey("DelUpAccount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1008)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_account(mid, is_deleted) VALUES(1008, 0) ON DUPLICATE KEY UPDATE is_deleted = 0")
			rows, err := d.DelUpAccount(c, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCreditRecord(t *testing.T) {
	convey.Convey("DelCreditRecord", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(1000)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO credit_score_record(id, is_deleted) VALEUS(1000, 0) ON DUPLICATE KEY UPDATE is_deleted = 0")
			rows, err := d.DelCreditRecord(c, id)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxDelCreditRecord(t *testing.T) {
	convey.Convey("TxDelCreditRecord", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			id    = int64(1000)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			Exec(c, "INSERT INTO credit_score_record(id, is_deleted) VALEUS(1000, 0) ON DUPLICATE KEY UPDATE is_deleted = 0")
			rows, err := d.TxDelCreditRecord(tx, id)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpInfo(t *testing.T) {
	convey.Convey("UpInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(1008)
			state = int64(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state, is_deleted) VALUES(1008, 3, 0) ON DUPLICATE KEY UPDATE account_state = 3, is_deleted = 0")
			info, err := d.UpInfo(c, mid, state)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpInfoByState(t *testing.T) {
	convey.Convey("GetUpInfoByState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "up_info_video"
			mids  = []int64{1008}
			state = int(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			info, err := d.GetUpInfoByState(c, table, mids, state)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpState(t *testing.T) {
	convey.Convey("GetUpState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "up_info_video"
			mid   = int64(1008)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			state, err := d.GetUpState(c, table, mid)
			ctx.Convey("Then err should be nil.state should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(state, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBGMCount(t *testing.T) {
	convey.Convey("BGMCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1008)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO background_music(sid, mid) VALUES(100, 1008)")
			count, err := d.BGMCount(c, mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
