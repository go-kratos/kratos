package service

import (
	"context"
	"go-common/app/job/main/vip/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//  go test  -test.v -test.run TestSyncUserInfoByMid
func TestSyncUserInfoByMid(t *testing.T) {
	Convey("SyncUserInfoByMid err == nil", t, func() {
		err := s.SyncUserInfoByMid(context.TODO(), 1002)
		So(err, ShouldBeNil)
	})
}

//  go test  -test.v -test.run TestUpdateDatabusUserInfo

func TestUpdateDatabusUserInfo(t *testing.T) {
	Convey("TestUpdateDatabusUserInfo err == nil", t, func() {
		var (
			mid int64 = 2089809
			old *model.VipUserInfoOld
			msg = new(model.VipUserInfoMsg)
			err error
		)
		old, err = s.dao.OldVipInfo(context.TODO(), mid)
		So(err, ShouldBeNil)
		msg.Mid = old.Mid
		msg.Type = old.Type
		msg.Status = old.Status
		msg.StartTime = old.StartTime.Time().Format("2006-01-02 15:04:05")
		msg.OverdueTime = old.OverdueTime.Time().Format("2006-01-02 15:04:05")
		msg.AnnualVipOverdueTime = old.AnnualVipOverdueTime.Time().Format("2006-01-02 15:04:05")
		msg.RecentTime = old.RecentTime.Time().Format("2006-01-02 15:04:05")
		msg.Wander = old.Wander
		msg.AutoRenewed = old.AutoRenewed
		msg.IsAutoRenew = old.IsAutoRenew
		msg.IosOverdueTime = old.IosOverdueTime.Time().Format("2006-01-02 15:04:05")
		userInfo := convertMsgToUserInfo(msg)
		err = s.addUserInfo(context.TODO(), userInfo)
		So(err, ShouldBeNil)
	})
}
