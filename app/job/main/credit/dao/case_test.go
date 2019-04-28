package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/credit/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AddBlockedCase(t *testing.T) {
	Convey("should return err be nil", t, func() {
		ca := &model.Case{}
		ca.Mid = 6660
		ca.BlockedDay = 7
		err := d.AddBlockedCase(context.TODO(), ca)
		So(err, ShouldBeNil)
	})
}

func Test_UpGrantCase(t *testing.T) {
	Convey("should return err be nil", t, func() {
		now := time.Now()
		err := d.UpGrantCase(context.TODO(), []int64{1, 2}, xtime.Time(now.Unix()), xtime.Time(now.Unix()))
		So(err, ShouldBeNil)
	})
}

func Test_Grantcase(t *testing.T) {
	Convey("should return err be nil", t, func() {
		sc, err := d.Grantcase(context.TODO(), 3)
		So(err, ShouldBeNil)
		So(sc, ShouldNotBeNil)
		So(sc, ShouldResemble, make(map[int64]*model.SimCase))
	})
}

func Test_CaseVote(t *testing.T) {
	Convey("should return err be nil", t, func() {
		cv, err := d.CaseVote(context.TODO(), 1004)
		So(err, ShouldBeNil)
		So(cv, ShouldNotBeNil)
	})
}

func Test_CaseRelationIDCount(t *testing.T) {
	Convey("should return err be nil", t, func() {
		count, err := d.CaseRelationIDCount(context.TODO(), 2, "2-8-113")
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func Test_CaseVotesCID(t *testing.T) {
	Convey("should return err be nil", t, func() {
		cv, err := d.CaseVotesCID(context.TODO(), 2)
		So(err, ShouldBeNil)
		So(cv, ShouldNotBeNil)
	})
}

func Test_CountCaseMID(t *testing.T) {
	Convey("should return err be nil", t, func() {
		count, err := d.CountCaseMID(context.TODO(), 1, 2)
		So(err, ShouldBeNil)
		So(count, ShouldNotBeNil)
	})
}
