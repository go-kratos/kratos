package dao

import (
	"testing"

	"go-common/app/admin/ep/saga/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestQueryUserByUserName(t *testing.T) {
	convey.Convey("Test QueryUserByUserName", t, func(ctx convey.C) {

		contactInfo, err := d.QueryUserByUserName("zhanglin")
		ctx.Convey("The err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(contactInfo.NickName, convey.ShouldEqual, "木木哥")
			ctx.So(contactInfo.UserID, convey.ShouldEqual, "003207")
		})
	})
}

func TestQueryUserByID(t *testing.T) {
	convey.Convey("Test QueryUserByID", t, func(ctx convey.C) {

		contactInfo, err := d.QueryUserByID("003207")
		ctx.Convey("The err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(contactInfo.UserName, convey.ShouldEqual, "zhanglin")
			ctx.So(contactInfo.NickName, convey.ShouldEqual, "木木哥")
		})
	})
}

func TestUserIds(t *testing.T) {
	convey.Convey("Test UserIds", t, func(ctx convey.C) {

		var (
			userNames = []string{"zhanglin", "wuwei"}
		)

		userIds, err := d.UserIds(userNames)
		ctx.Convey("The err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(userIds, convey.ShouldEqual, "003207|001396")
		})
	})
}

func TestWechatContact(t *testing.T) {
	convey.Convey("test wechat contact", t, func(ctx convey.C) {
		var (
			ID = "111"
		)

		contact := &model.ContactInfo{
			ID:          ID,
			UserName:    "lisi",
			UserID:      "222",
			NickName:    "sansan",
			VisibleSaga: true,
		}

		contactUpdate := &model.ContactInfo{
			ID:          ID,
			UserName:    "zhan",
			UserID:      "333",
			NickName:    "sansan",
			VisibleSaga: true,
		}

		ctx.Convey("set contact", func(ctx convey.C) {
			err := d.CreateContact(contact)
			ctx.Convey("set err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("update contact", func(ctx convey.C) {
			err := d.UptContact(contactUpdate)
			ctx.Convey("get err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("delete contact", func(ctx convey.C) {
			err := d.DelContact(contactUpdate)
			ctx.Convey("get err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestWechatCreateLog(t *testing.T) {
	convey.Convey("test wechat create log", t, func(ctx convey.C) {

		info := &model.WechatCreateLog{
			Name:   "lisi",
			Owner:  "zhanglin",
			ChatID: "333",
			Cuser:  "333",
		}

		info1 := &model.WechatCreateLog{
			Name: "lisi",
		}

		ctx.Convey("set wechat log", func(ctx convey.C) {
			err := d.AddWechatCreateLog(info)
			ctx.Convey("set err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("query wechat log", func(ctx convey.C) {
			infos, total, err := d.QueryWechatCreateLog(false, nil, info1)
			ctx.Convey("query err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldBeGreaterThanOrEqualTo, 1)
				ctx.So(infos[0].Name, convey.ShouldEqual, "lisi")
			})
		})
	})
}
