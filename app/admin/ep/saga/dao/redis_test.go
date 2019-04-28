package dao

import (
	"context"
	"testing"

	"go-common/app/admin/ep/saga/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoWeixinTokenKeyRedis(t *testing.T) {
	convey.Convey("weixinTokenKeyRedis", t, func(ctx convey.C) {
		var (
			key = "111"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := weixinTokenKeyRedis(key)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, "saga_weixin_token__111")
			})
		})
	})
}

func TestDaoAccessTokenRedis(t *testing.T) {
	convey.Convey("AccessTokenRedis", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			key   = "111"
			token = "sdfgsdgfdg"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {

			err := d.SetAccessTokenRedis(c, key, token, -1)
			ctx.Convey("Set Access Token. Then err should be nil. ", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})

			token, err := d.AccessTokenRedis(c, key)
			ctx.Convey("get access token. Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(token, convey.ShouldEqual, token)
			})
		})
	})
}

func TestDaoRequireVisibleUsersRedis(t *testing.T) {
	convey.Convey("RequireVisibleUsersRedis", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			userID      = "222"
			contactInfo = &model.ContactInfo{
				ID:          "111",
				UserName:    "zhanglin",
				UserID:      userID,
				NickName:    "mumuge",
				VisibleSaga: true,
			}
			userMap = make(map[string]model.RequireVisibleUser)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {

			err := d.SetRequireVisibleUsersRedis(c, contactInfo)
			ctx.Convey("Set Visible Users. Then err should be nil. ", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})

			err = d.RequireVisibleUsersRedis(c, &userMap)
			ctx.Convey("get Visible Users. Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(userMap[userID].UserName, convey.ShouldEqual, "zhanglin")
				ctx.So(userMap[userID].NickName, convey.ShouldEqual, "mumuge")
			})

			err = d.DeleteRequireVisibleUsersRedis(c)
			ctx.Convey("delete Visible Users. Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetItemRedis(t *testing.T) {
	convey.Convey("GetItemRedis", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			key         = "333"
			contactInfo = &model.ContactInfo{
				ID:          "111",
				UserName:    "zhanglin",
				UserID:      "333",
				NickName:    "mumuge",
				VisibleSaga: true,
			}
			getContactInfo *model.ContactInfo
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {

			err := d.SetItemRedis(c, key, contactInfo, 0)
			ctx.Convey("Set Item. Then err should be nil. ", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})

			err = d.ItemRedis(c, key, &getContactInfo)
			ctx.Convey("get Item. Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(getContactInfo.UserName, convey.ShouldEqual, "zhanglin")
				ctx.So(getContactInfo.UserID, convey.ShouldEqual, "333")
				ctx.So(getContactInfo.NickName, convey.ShouldEqual, "mumuge")
				ctx.So(getContactInfo.VisibleSaga, convey.ShouldEqual, true)
			})
		})
	})
}
