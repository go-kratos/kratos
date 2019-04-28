package service

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/account-recovery/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceQueryAccount(t *testing.T) {
	var (
		c   = context.Background()
		req = &model.QueryInfoReq{
			QType:  "uid",
			QValue: "1",
			CToken: "f0bd2dd65c444697b7b1330bb757bf3d",
			Code:   "gp74",
		}
	)
	convey.Convey("QueryAccount", t, func(ctx convey.C) {
		res, err := s.QueryAccount(c, req)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestServiceCommitInfo(t *testing.T) {
	var (
		c     = context.Background()
		uinfo = &model.UserInfoReq{
			LoginAddrs:   "地址1,地址2",
			RegTime:      xtime.Time(time.Now().Unix()),
			RegType:      1,
			RegAddr:      "注册地址",
			Unames:       "昵称1",
			Pwds:         "密码1",
			Phones:       "手机1",
			Emails:       "邮件1",
			SafeQuestion: 0,
			SafeAnswer:   "",
			CardID:       "ISN-367890",
			CardType:     1,
			Captcha:      "6789",
			LinkMail:     "2459593393@qq.com",
			Mid:          8,
		}
	)
	convey.Convey("CommitInfo", t, func(ctx convey.C) {
		err := s.CommitInfo(c, uinfo)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceQueryCon(t *testing.T) {
	var (
		defint int64
		c      = context.Background()
		aq     = &model.QueryRecoveryInfoReq{
			RID:        1,
			UID:        0,
			Status:     &defint,
			Game:       &defint,
			Size:       2,
			StartTime:  1533206284,
			EndTime:    1539206284,
			IsAdvanced: false,
			Page:       1,
		}
	)
	convey.Convey("QueryCon", t, func(ctx convey.C) {
		res, err := s.QueryCon(c, aq)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestServiceJudge(t *testing.T) {
	var (
		c   = context.Background()
		req = &model.JudgeReq{
			Rid:      1,
			Status:   1,
			Operator: "hyy",
			OptTime:  xtime.Time(time.Now().Unix()),
			Remark:   "申诉通过",
		}
	)
	convey.Convey("Judge", t, func(ctx convey.C) {
		err := s.Judge(c, req)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceBatchJudge(t *testing.T) {
	var (
		c       = context.Background()
		ridsAry = []int64{1, 2}
		req     = &model.BatchJudgeReq{
			Rids:     "1,2",
			Status:   1,
			Operator: "hyy",
			OptTime:  xtime.Time(time.Now().Unix()),
			Remark:   "批量申诉驳回",
			RidsAry:  ridsAry,
		}
	)
	convey.Convey("BatchJudge", t, func(ctx convey.C) {
		err := s.BatchJudge(c, req)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceGetCaptchaMail(t *testing.T) {
	var (
		c   = context.Background()
		req = &model.CaptchaMailReq{Mid: 1, LinkMail: "2459593393@qq.com"}
	)
	convey.Convey("GetCaptchaMail", t, func(ctx convey.C) {
		state, err := s.GetCaptchaMail(c, req)
		ctx.Convey("Then err should be nil.state should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(state, convey.ShouldNotBeNil)
		})
	})
}

func TestServiceSendMail(t *testing.T) {
	var (
		c   = context.Background()
		req = &model.SendMailReq{RID: 1, Status: 1}
	)
	convey.Convey("SendMail", t, func(ctx convey.C) {
		err := s.SendMail(c, req)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServicebatchJudge(t *testing.T) {
	var (
		c        = context.Background()
		status   = int64(1)
		rids     = []int64{1, 2}
		operator = "abcd"
		optTime  xtime.Time
		remark   = "通过"
	)
	convey.Convey("batchJudge", t, func(ctx convey.C) {
		err := s.batchJudge(c, status, rids, operator, optTime, remark)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServicedeal(t *testing.T) {
	var (
		c   = context.Background()
		rid = int64(1)
	)
	convey.Convey("deal", t, func(ctx convey.C) {
		err := s.deal(c, rid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceagree(t *testing.T) {
	var (
		c   = context.Background()
		rid = int64(1)
	)
	convey.Convey("agree", t, func(ctx convey.C) {
		err := s.agree(c, rid, "账号找回服务")
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServicereject(t *testing.T) {
	var (
		c   = context.Background()
		rid = int64(1)
	)
	convey.Convey("reject", t, func(ctx convey.C) {
		err := s.reject(c, rid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceagreeMore(t *testing.T) {
	var (
		c       = context.Background()
		ridsStr = "1,2"
	)
	convey.Convey("agreeMore", t, func(ctx convey.C) {
		err := s.agreeMore(c, ridsStr, "账号找回服务")
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServicerejectMore(t *testing.T) {
	var (
		c       = context.Background()
		ridsStr = "1,2"
	)
	convey.Convey("rejectMore", t, func(ctx convey.C) {
		err := s.rejectMore(c, ridsStr)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceWebToken(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("WebToken", t, func(ctx convey.C) {
		token, err := s.WebToken(c)
		ctx.Convey("Then err should be nil.token should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(token, convey.ShouldNotBeNil)
		})
	})
}

func TestServicehideUID(t *testing.T) {
	var (
		mid = "1"
	)
	convey.Convey("hideUID", t, func(ctx convey.C) {
		uid := hideUID(mid)
		ctx.Convey("Then uid should not be nil.", func(ctx convey.C) {
			ctx.So(uid, convey.ShouldNotBeNil)
		})
	})
}

func TestServiceGameList(t *testing.T) {
	var (
		c    = context.Background()
		mids = "1,2"
	)
	convey.Convey("GameList", t, func(ctx convey.C) {
		res, err := s.GameList(c, mids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestServiceQueryRecoveryAddit(t *testing.T) {
	convey.Convey("QueryRecoveryAddit", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			bizData, err := s.QueryRecoveryAddit(c, rids)
			ctx.Convey("Then err should be nil.bizData should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bizData, convey.ShouldNotBeNil)
			})
		})
	})
}
