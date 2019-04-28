package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/account-recovery/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetStatusByRid(t *testing.T) {
	var (
		c   = context.Background()
		rid = int64(1)
	)
	convey.Convey("GetStatusByRid", t, func(ctx convey.C) {
		status, err := d.GetStatusByRid(c, rid)
		ctx.Convey("Then err should be nil.status should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(status, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetSuccessCount(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1234)
	)
	convey.Convey("GetSuccessCount", t, func(ctx convey.C) {
		count, err := d.GetSuccessCount(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBatchGetRecoverySuccess(t *testing.T) {
	var (
		c    = context.Background()
		mids = []int64{1234}
	)
	convey.Convey("BatchGetRecoverySuccess", t, func(ctx convey.C) {
		countMap, err := d.BatchGetRecoverySuccess(c, mids)
		ctx.Convey("Then err should be nil.countMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(countMap, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateSuccessCount(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1234)
	)
	convey.Convey("UpdateSuccessCount", t, func(ctx convey.C) {
		err := d.UpdateSuccessCount(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoBatchUpdateSuccessCount(t *testing.T) {
	var (
		c    = context.Background()
		mids = "1234"
	)
	convey.Convey("BatchUpdateSuccessCount", t, func(ctx convey.C) {
		err := d.BatchUpdateSuccessCount(c, mids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGetNoDeal(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1234)
	)
	convey.Convey("GetNoDeal", t, func(ctx convey.C) {
		count, err := d.GetNoDeal(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateStatus(t *testing.T) {
	var (
		c        = context.Background()
		status   = int64(1)
		rid      = int64(1)
		operator = "abcd"
		optTime  xtime.Time
		remark   = ""
	)
	convey.Convey("UpdateStatus", t, func(ctx convey.C) {
		err := d.UpdateStatus(c, status, rid, operator, optTime, remark)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdateUserType(t *testing.T) {
	var (
		c      = context.Background()
		status = int64(1)
		rid    = int64(1)
	)
	convey.Convey("UpdateUserType", t, func(ctx convey.C) {
		err := d.UpdateUserType(c, status, rid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoInsertRecoveryInfo(t *testing.T) {
	var (
		c     = context.Background()
		uinfo = &model.UserInfoReq{
			LoginAddrs: "中国-福州,中国-上海,澳大利亚",
			//RegTime:timeS,  //变成2018
			RegTime:      1533206284, //变成2018  //数据库设置为int(11),so数据库必须设置为tiimestamp
			RegType:      int8(1),
			RegAddr:      "中国上海",
			Unames:       "昵称AA,昵称BB,昵称CC",
			Pwds:         "密码1,密码2",
			Phones:       "12345678901,54321678923",
			Emails:       "2456@sina.com,789@qq.com",
			SafeQuestion: int8(1),
			SafeAnswer:   "心态呀",
			CardID:       "ISN-1234567890-0987",
			CardType:     int8(1),
			LinkMail:     "345678@qq.com",
			Mid:          1234,
		}
	)
	convey.Convey("InsertRecoveryInfo", t, func(ctx convey.C) {
		lastID, err := d.InsertRecoveryInfo(c, uinfo)
		ctx.Convey("Then err should be nil.lastID should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(lastID, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateSysInfo(t *testing.T) {
	var (
		c   = context.Background()
		sys = &model.SysInfo{
			SysLoginAddrs: "中国-福州,中国-上海,澳大利亚",
			SysReg:        "对",
			SysUNames:     "对,错,错",
			SysPwds:       "对,错",
			SysPhones:     "对,错",
			SysEmails:     "对,错",
			SysSafe:       "对",
			SysCard:       "对",
		}
		userType = int64(1)
		rid      = int64(1)
	)
	convey.Convey("UpdateSysInfo", t, func(ctx convey.C) {
		err := d.UpdateSysInfo(c, sys, userType, rid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoQueryByID(t *testing.T) {
	var (
		c                   = context.Background()
		rid                 = int64(1)
		fromTime xtime.Time = 1533120949
		endTime  xtime.Time = 1535636392
	)
	convey.Convey("QueryByID", t, func(ctx convey.C) {
		res, err := d.QueryByID(c, rid, fromTime, endTime)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoQueryInfoByLimit(t *testing.T) {
	var (
		c   = context.Background()
		req = &model.DBRecoveryInfoParams{
			ExistGame:   false,
			ExistStatus: false,
			ExistMid:    false,
			Mid:         0,
			Game:        0,
			Status:      1,
			FirstRid:    20,
			LastRid:     0,
			Size:        2,
			StartTime:   1533120949,
			EndTime:     1535636392,
			SubNum:      1,
			CurrPage:    1,
		}
	)
	convey.Convey("QueryInfoByLimit", t, func(ctx convey.C) {
		res, total, err := d.QueryInfoByLimit(c, req)
		ctx.Convey("Then err should be nil.res,total should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(total, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetUinfoByRid(t *testing.T) {
	var (
		c   = context.Background()
		rid = int64(240)
	)
	convey.Convey("GetUinfoByRid", t, func(ctx convey.C) {
		mid, linkMail, ctime, err := d.GetUinfoByRid(c, rid)
		ctx.Convey("Then err should be nil.mid,linkMail,ctime should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ctime, convey.ShouldNotBeNil)
			ctx.So(linkMail, convey.ShouldNotBeNil)
			ctx.So(mid, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetUinfoByRidMore(t *testing.T) {
	var (
		c       = context.Background()
		ridsStr = "1,2"
	)
	convey.Convey("GetUinfoByRidMore", t, func(ctx convey.C) {
		bathRes, err := d.GetUinfoByRidMore(c, ridsStr)
		ctx.Convey("Then err should be nil.bathRes should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(bathRes, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetUnCheckInfo(t *testing.T) {
	var (
		c   = context.Background()
		rid = int64(1)
	)
	convey.Convey("GetUnCheckInfo", t, func(ctx convey.C) {
		r, err := d.GetUnCheckInfo(c, rid)
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(r, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetMailStatus(t *testing.T) {
	var (
		c   = context.Background()
		rid = int64(1)
	)
	convey.Convey("GetMailStatus", t, func(ctx convey.C) {
		mailStatus, err := d.GetMailStatus(c, rid)
		ctx.Convey("Then err should be nil.mailStatus should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(mailStatus, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateMailStatus(t *testing.T) {
	var (
		c   = context.Background()
		rid = int64(1)
	)
	convey.Convey("UpdateMailStatus", t, func(ctx convey.C) {
		err := d.UpdateMailStatus(c, rid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGetAllByCon(t *testing.T) {
	var (
		defint int64
		c      = context.Background()
		aq     = &model.QueryRecoveryInfoReq{
			//RID: 1,
			//UID:2,
			Status:    &defint,
			Game:      &defint,
			Size:      10,
			Page:      1,
			StartTime: 1533052800,
			EndTime:   1536924163,

			//StartTime  time.Time `json:"start_time" form:"start_time"`
			//EndTime    time.Time `json:"end_time" form:"end_time"`
			//IsAdvanced bool      `json:"-"`
			//Page       int64     `form:"page"`
		}
	)
	convey.Convey("UpdateMailStatus", t, func(ctx convey.C) {
		resultData, total, err := d.GetAllByCon(c, aq)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(resultData, convey.ShouldNotBeNil)
			ctx.So(total, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldBeNil)
			ctx.Println(total, "  len=", len(resultData))
		})
	})
}

func TestDaoInsertRecoveryAddit(t *testing.T) {
	var (
		c           = context.Background()
		rid   int64 = 1
		files       = "http://uat-i0.hdslb.com/bfs/account/recovery/bca2.zip,http://uat-i0.hdslb.com/bfs/account/recovery/abcd.zip"
		extra       = `{"GameArea":"ios-A服","GameNames":"崩坏3","GamePlay":"1"}`
	)
	convey.Convey("InsertRecoveryAddit", t, func(ctx convey.C) {
		err := d.InsertRecoveryAddit(c, rid, files, extra)
		err = d.InsertRecoveryAddit(c, 2, files, extra)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdateRecoveryAddit(t *testing.T) {
	var (
		c           = context.Background()
		rid   int64 = 1
		files       = []string{"http://uat-i0.hdslb.com/bfs/aaaa.zip", "http://uat-i0.hdslb.com/bfs/dddd.zip"}
		extra       = `{"GameArea":"ios-A服","GameNames":"崩坏3","GamePlay":"1"}`
	)
	convey.Convey("UpdateMailStatus", t, func(ctx convey.C) {
		err := d.UpdateRecoveryAddit(c, rid, files, extra)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGetRecoveryAddit(t *testing.T) {
	var (
		c         = context.Background()
		rid int64 = 1
	)
	convey.Convey("UpdateMailStatus", t, func(ctx convey.C) {
		addit, err := d.GetRecoveryAddit(c, rid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.Println(addit, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBatchGetRecoveryAddit(t *testing.T) {
	var (
		c    = context.Background()
		rids = []int64{1, 2}
	)
	convey.Convey("BatchGetRecoveryAddit", t, func(ctx convey.C) {
		addits, err := d.BatchGetRecoveryAddit(c, rids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(addits, convey.ShouldNotBeNil)
		})
	})
}

func TestBatchGetLastSuccess(t *testing.T) {
	var (
		c    = context.Background()
		mids = []int64{1234}
	)
	convey.Convey("BatchGetLastSuccess", t, func(ctx convey.C) {
		lastSuccessMap, err := d.BatchGetLastSuccess(c, mids)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(lastSuccessMap, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetLastSuccess(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1234)
	)
	convey.Convey("GetLastSuccess", t, func(ctx convey.C) {
		res, err := d.GetLastSuccess(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
