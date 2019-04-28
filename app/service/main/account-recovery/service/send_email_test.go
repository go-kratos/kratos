package service

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/account-recovery/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceSendMailM(t *testing.T) {
	var (
		c        = context.Background()
		mailType = int(1)
		params   = "1234"
		linkMail = "2459593393@qq.com"
	)
	convey.Convey("SendMailM", t, func(ctx convey.C) {
		err := s.SendMailM(c, mailType, linkMail, params)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceSendMailMany(t *testing.T) {
	var (
		c        = context.Background()
		mailType = int(0)
		endDate  = time.Now()
		timeS    = xtime.Time(endDate.Unix())
		batchRes = []*model.BatchAppeal{
			{
				Rid:      "1",
				Mid:      "1",
				LinkMail: "2459593393@qq.com",
				Ctime:    timeS,
			},
			{
				Rid:      "2",
				Mid:      "2",
				LinkMail: "1772968069@qq.com",
				Ctime:    timeS,
			},
		}
		userMap = make(map[string]*model.User)
		user1   = model.User{UserID: "mid=1", Pwd: "mid=1的pwd"}
		user2   = model.User{UserID: "mid=2", Pwd: "mid=2的pwd"}
	)
	userMap["1"] = &user1
	userMap["2"] = &user2
	convey.Convey("SendMailMany", t, func(ctx convey.C) {
		err := s.SendMailMany(c, mailType, batchRes, userMap)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceSendMailLog(t *testing.T) {
	var (
		c              = context.Background()
		mid      int64 = 1
		mailType       = int(1)
		params         = "1234"
		linkMail       = "2459593393@qq.com"
	)
	convey.Convey("SendMailLog", t, func(ctx convey.C) {
		err := s.SendMailLog(c, mid, mailType, linkMail, params)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
