package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/answer/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestHit(t *testing.T) {
	convey.Convey("hit", t, func() {
		p1 := hit(0)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestDaoPendantHistory(t *testing.T) {
	convey.Convey("PendantHistory", t, func() {
		hid, status, err := d.PendantHistory(context.Background(), 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(status, convey.ShouldNotBeNil)
		convey.So(hid, convey.ShouldNotBeNil)
	})
}

func TestDaoHistory(t *testing.T) {
	convey.Convey("History", t, func() {
		res, err := d.History(context.Background(), 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestDaoHistoryByHid(t *testing.T) {
	convey.Convey("HistoryByHid", t, func() {
		res, err := d.HistoryByHid(context.Background(), 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestDaoOldHistory(t *testing.T) {
	convey.Convey("OldHistory", t, func() {
		res, err := d.OldHistory(context.Background(), 0, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestDaoSharingIndexByHid(t *testing.T) {
	convey.Convey("SharingIndexByHid", t, func() {
		res, err := d.SharingIndexByHid(context.Background(), 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestDaoAddPendantHistory(t *testing.T) {
	convey.Convey("AddPendantHistory", t, func() {
		affected, err := d.AddPendantHistory(context.Background(), time.Now().Unix(), time.Now().Unix())
		convey.So(err, convey.ShouldBeNil)
		convey.So(affected, convey.ShouldNotBeNil)
	})
}

func TestDaoHistorySuit(t *testing.T) {
	var (
		err      error
		affected int64
		hid      = time.Now().Unix()
	)
	ans := &model.AnswerHistory{
		Hid: hid,
	}
	convey.Convey("AddHistory", t, func() {
		affected, _, err = d.AddHistory(context.Background(), 0, ans)
		convey.So(err, convey.ShouldBeNil)
		convey.So(affected, convey.ShouldNotBeNil)
	})
	convey.Convey("HistoryByHid", t, func() {
		ans, err = d.HistoryByHid(context.Background(), hid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(affected, convey.ShouldNotBeNil)
	})
	convey.Convey("SetHistory", t, func() {
		affected, err = d.SetHistory(context.Background(), 0, ans)
		convey.So(err, convey.ShouldBeNil)
		convey.So(affected, convey.ShouldNotBeNil)
	})
}

func TestDaoUpdateLevel(t *testing.T) {
	convey.Convey("UpdateLevel", t, func() {
		ret, err := d.UpdateLevel(context.Background(), 0, 0, 0, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ret, convey.ShouldNotBeNil)
	})
}

func TestDaoUpdateCaptcha(t *testing.T) {
	convey.Convey("UpdateCaptcha", t, func() {
		ret, err := d.UpdateCaptcha(context.Background(), 0, 0, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ret, convey.ShouldNotBeNil)
	})
}

func TestDaoUpdateStepTwoTime(t *testing.T) {
	convey.Convey("UpdateStepTwoTime", t, func() {
		ret, err := d.UpdateStepTwoTime(context.Background(), 0, 0, time.Now())
		convey.So(err, convey.ShouldBeNil)
		convey.So(ret, convey.ShouldNotBeNil)
	})
}

func TestDaoUpdateExtraStartTime(t *testing.T) {
	convey.Convey("UpdateExtraStartTime", t, func() {
		ret, err := d.UpdateExtraStartTime(context.Background(), 0, 0, time.Now())
		convey.So(err, convey.ShouldBeNil)
		convey.So(ret, convey.ShouldNotBeNil)
	})
}

func TestDaoUpdateExtraRet(t *testing.T) {
	convey.Convey("UpdateExtraRet", t, func() {
		ret, err := d.UpdateExtraRet(context.Background(), 0, 0, 0, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ret, convey.ShouldNotBeNil)
	})
}

func TestDaoUpPendantHistory(t *testing.T) {
	convey.Convey("UpPendantHistory", t, func() {
		ret, err := d.UpPendantHistory(context.Background(), 0, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ret, convey.ShouldNotBeNil)
	})
}

func TestDaoQidsExtraByState(t *testing.T) {
	convey.Convey("QidsExtraByState", t, func() {
		_, err := d.QidsExtraByState(context.Background(), 1)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoExtraByIds(t *testing.T) {
	convey.Convey("ExtraByIds", t, func() {
		res, err := d.ExtraByIds(context.Background(), []int64{1, 2, 3})
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestDaoTypes(t *testing.T) {
	convey.Convey("Types", t, func() {
		_, err := d.Types(context.Background())
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoByIds(t *testing.T) {
	convey.Convey("ByIds", t, func() {
		_, err := d.ByIds(context.Background(), []int64{1, 2, 3})
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoQidsByState(t *testing.T) {
	convey.Convey("QidsByState", t, func() {
		_, err := d.QidsByState(context.Background(), 1)
		convey.So(err, convey.ShouldBeNil)
	})
}
