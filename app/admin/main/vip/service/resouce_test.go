package service

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/vip/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GrandResouce(t *testing.T) {
	Convey("should return true err == nil and failMids is empty", t, func() {
		var (
			mids     []int
			remark   = "系统备注"
			batchID  = 5
			username = "system"
		)
		mids = append(mids, 110)
		failMids, err := s.GrandResouce(context.TODO(), remark, int64(batchID), mids, username)
		So(err, ShouldBeNil)
		So(failMids, ShouldBeEmpty)
	})
}

func TestService_BatchInfo(t *testing.T) {
	Convey("batch info", t, func() {
		id := 12
		_, err := s.BatchInfo(context.TODO(), id)
		So(err, ShouldBeNil)
	})
}
func TestService_PoolInfo(t *testing.T) {
	Convey("pool info", t, func() {
		id := 12
		_, err := s.PoolInfo(context.TODO(), id)
		So(err, ShouldBeNil)
	})
}

func Test_UpdateResouce(t *testing.T) {
	Convey("should ", t, func() {
		pojo := new(model.ResoucePoolBo)
		pojo.ID = 25
		pojo.PoolName = "test123123112311"
		pojo.BusinessID = 3
		pojo.Reason = "zhaozhihao"
		pojo.CodeExpireTime = xtime.Time(time.Now().Unix())
		pojo.StartTime = xtime.Time(time.Now().AddDate(0, 0, -1).Unix())
		pojo.EndTime = xtime.Time(time.Now().AddDate(0, 0, 10).Unix())
		pojo.Contacts = "阿斯顿发"
		pojo.ContactsNumber = "123124123"
		err := s.UpdatePool(context.TODO(), pojo)
		So(err, ShouldBeNil)
	})
}

func Test_SaveBatchCode(t *testing.T) {
	Convey("testing ", t, func() {
		arg := new(model.BatchCode)
		arg.ID = 26
		arg.PoolID = 25
		arg.Type = 1
		arg.BusinessID = 3
		arg.BatchName = "测试123"
		arg.SurplusCount = 100000
		arg.Count = 20000
		arg.Unit = 366
		arg.LimitDay = 9
		arg.MaxCount = 5
		arg.StartTime = xtime.Time(time.Now().Unix())
		arg.EndTime = xtime.Time(time.Now().AddDate(0, 0, 1).Unix())
		arg.Price = 10
		arg.Reason = "zhaozhihao"
		err := s.SaveBatchCode(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}

func Test_vipInfo(t *testing.T) {
	Convey("testing", t, func() {
		mid := 123
		res, err := s.VipInfo(context.TODO(), int64(mid))
		t.Logf("%+v", res)
		So(err, ShouldBeNil)
	})
}

func Test_drawback(t *testing.T) {
	Convey("drawback", t, func() {
		mid := 2089809
		username := "zhaozhihao"
		remark := "zhaozhihao"
		day := 10
		err := s.Drawback(context.TODO(), day, int64(mid), username, remark)
		So(err, ShouldBeNil)
	})
}

func TestService_ExportCode(t *testing.T) {
	Convey(" export code", t, func() {
		codes, err := s.ExportCode(context.TODO(), 13)
		So(codes, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
