package dao

import (
	"context"
	"testing"
	"time"

	gock "gopkg.in/h2non/gock.v1"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoLoginout(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("TestDaoLoginout", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.loginOutURL).Reply(200).JSON(`{"code":0}`)
		err := d.Loginout(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSendCleanCache(t *testing.T) {
	var (
		c      = context.TODO()
		mid    = int64(0)
		months = int16(0)
		days   = int64(0)
		no     = int(0)
		ip     = ""
	)
	convey.Convey("SendCleanCache", t, func(ctx convey.C) {
		err := d.SendCleanCache(c, mid, months, days, no, ip)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSendMultipMsg(t *testing.T) {
	var (
		c        = context.TODO()
		mids     = ""
		content  = ""
		title    = ""
		mc       = ""
		ip       = ""
		dataType = int(0)
	)
	convey.Convey("SendMultipMsg", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.msgURI).Reply(200).JSON(`{"code":0}`)
		err := d.SendMultipMsg(c, mids, content, title, mc, ip, dataType)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSendBcoinCoupon(t *testing.T) {
	var (
		c          = context.TODO()
		mids       = "20606508"
		activityID = "1"
		money      = int64(0)
		dueTime    = time.Now()
	)
	convey.Convey("SendBcoinCoupon", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.payURI).Reply(200).JSON(`{"code":0}`)
		err := d.SendBcoinCoupon(c, mids, activityID, money, dueTime)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
