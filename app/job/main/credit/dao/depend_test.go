package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/credit/model"

	. "github.com/smartystreets/goconvey/convey"

	"gopkg.in/h2non/gock.v1"
)

func Test_SendPendant(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.SendPendant(context.TODO(), 27515305, []int64{83, 84}, 30)
		So(err, ShouldBeNil)
	})
}

func Test_SendMedal(t *testing.T) {
	Convey("should return err be nil", t, func() {
		defer gock.OffAll()
		httpMock("POST", d.sendMedalURL).Reply(200).JSON(`{"code":0}`)
		err := d.SendMedal(context.TODO(), 27515305, 4)
		So(err, ShouldBeNil)
	})
}

func Test_DelTag(t *testing.T) {
	Convey("should return err be nil", t, func() {
		defer gock.OffAll()
		httpMock("POST", d.delTagURL).Reply(200).JSON(`{"code":0}`)
		err := d.DelTag(context.TODO(), "1", "4")
		So(err, ShouldBeNil)
	})
}

func Test_ReportDM(t *testing.T) {
	Convey("should return err be nil", t, func() {
		defer gock.OffAll()
		httpMock("POST", d.delDMURL).Reply(200).JSON(`{"code":0}`)
		err := d.ReportDM(context.TODO(), "1", "4", 1)
		So(err, ShouldBeNil)
	})
}
func Test_BlockAccount(t *testing.T) {
	var bi = &model.BlockedInfo{}
	bi.UID = 27515305
	Convey("should return err be nil", t, func() {
		err := d.BlockAccount(context.TODO(), bi)
		So(err, ShouldBeNil)
	})
}

func TestDaoUnBlockAccount(t *testing.T) {
	var bi = &model.BlockedInfo{}
	bi.UID = 27515305
	Convey("should return err be nil", t, func() {
		err := d.UnBlockAccount(context.TODO(), bi)
		So(err, ShouldBeNil)
	})
}

func Test_SendMsg(t *testing.T) {
	Convey("should return err be nil", t, func() {
		defer gock.OffAll()
		httpMock("POST", d.sendMsgURL).Reply(200).JSON(`{"code":0}`)
		err := d.SendMsg(context.TODO(), 88889017, "风纪委员任期结束", "您的风纪委员资格已到期，任职期间的总结报告已生成。感谢您对社区工作的大力支持！#{点击查看}{"+"\"http://www.bilibili.com/judgement/\""+"}")
		So(err, ShouldBeNil)
	})
}

func TestDaoSms(t *testing.T) {
	Convey("should return err be nil", t, func() {
		defer gock.OffAll()
		httpMock("GET", "http://ops-mng.bilibili.co/api/sendsms").Reply(200).JSON(`{"code":0}`)
		err := d.Sms(context.TODO(), "13888888888", "1232131", "hahh1")
		So(err, ShouldBeNil)
	})
}

func Test_AddMoral(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.AddMoral(context.TODO(), 27515305, -10, 1, "credit", "reason", "风纪委处罚", "")
		So(err, ShouldBeNil)
	})
}

func Test_AddMoney(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.AddMoney(context.TODO(), 27956255, 10, "风纪委奖励")
		So(err, ShouldBeNil)
	})
}

func TestDaoCheckFilter(t *testing.T) {
	Convey("should return err be nil", t, func() {
		res, err := d.CheckFilter(context.TODO(), "credit", "lalsdal1", "127.0.0.1")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
func TestDaoUpAppealState(t *testing.T) {
	Convey("UpAppealState", t, func(ctx C) {
		var (
			c   = context.Background()
			bid = model.AppealBusinessID
			oid = int64(1)
			eid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx C) {
			defer gock.OffAll()
			httpMock("POST", d.upAppealStateURL).Reply(200).JSON(`{"code":0}`)
			err := d.UpAppealState(c, bid, oid, eid)
			ctx.Convey("Then err should be nil.", func(ctx C) {
				ctx.So(err, ShouldBeNil)
			})
		})
	})
}
