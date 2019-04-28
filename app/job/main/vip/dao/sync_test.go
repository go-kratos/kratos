package dao

import (
	"context"
	"go-common/app/job/main/vip/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSyncAddUser(t *testing.T) {
	var r = &model.VipUserInfo{Mid: 7993}
	convey.Convey("clean data before", t, func() {
		d.db.Exec(context.Background(), "delete from vip_user_info where mid = ?", r.Mid)
	})
	convey.Convey("SyncAddUser", t, func() {
		err := d.SyncAddUser(context.Background(), r)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("SyncUpdateUser", t, func() {
		r.Status = 2
		eff, err := d.SyncUpdateUser(context.Background(), r, r.Ver)
		convey.So(err, convey.ShouldBeNil)
		convey.So(eff, convey.ShouldNotBeNil)
	})
	convey.Convey("clean data", t, func() {
		d.db.Exec(context.Background(), "delete from vip_user_info where mid = ?", r.Mid)
	})
}
