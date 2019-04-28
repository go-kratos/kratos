package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/credit/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AddBlockInfo(t *testing.T) {
	var (
		r = model.BlockedInfo{}
	)
	r.OriginTitle = "go"
	r.OriginURL = "http:go"
	r.OriginType = 1
	r.OriginContent = "goc"
	r.OriginContentModify = "gocm"
	r.BlockedDays = 1
	r.BlockedForever = 1
	r.BlockedType = 1
	r.UID = 888890
	r.OperatorName = "lgs"
	r.PunishType = 3
	r.ReasonType = 3
	r.CaseID = 10
	Convey("should return err be nil", t, func() {
		_, err := d.AddBlockInfo(context.TODO(), &r, time.Now())
		So(err, ShouldBeNil)
	})
}

func Test_UpdateKPIPendentStatus(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.UpdateKPIPendentStatus(context.TODO(), 1)
		So(err, ShouldBeNil)
	})
}

func Test_UpdateKPIHandlerStatus(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.UpdateKPIHandlerStatus(context.TODO(), 1)
		So(err, ShouldBeNil)
	})
}

func Test_UpdateCase(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.UpdateCase(context.TODO(), model.CaseStatusDealed, model.JudgeTypeUndeal, 304)
		So(err, ShouldBeNil)
	})
}

func Test_InvalidJury(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.InvalidJury(context.TODO(), 1, 88889021)
		So(err, ShouldBeNil)
	})
}

func Test_UpdateVoteRight(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.UpdateVoteRight(context.TODO(), 88889021)
		So(err, ShouldBeNil)
	})
}

func Test_UpdateVoteTotal(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.UpdateVoteTotal(context.TODO(), 88889021)
		So(err, ShouldBeNil)
	})
}

func Test_UpdatePunishResult(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.UpdatePunishResult(context.TODO(), 1, 6)
		So(err, ShouldBeNil)
	})
}

func Test_BlockCount(t *testing.T) {
	Convey("should return err be nil and count>=0", t, func() {
		count, err := d.CountBlocked(context.TODO(), 88889021, time.Now())
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func TestDao_CaseByID(t *testing.T) {
	Convey("should return err be nil & res not be nil", t, func() {
		res, err := d.CaseByID(context.TODO(), 348)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
