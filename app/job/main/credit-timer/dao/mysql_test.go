package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/credit-timer/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpdateKPI(t *testing.T) {
	r := &model.Kpi{}
	r.Day = time.Now()
	r.Mid = 111
	r.Rate = 1
	r.Rank = 10
	r.RankPer = 10
	r.RankTotal = 100
	Convey("should return err be nil", t, func() {
		err := d.UpdateKPI(context.TODO(), r)
		So(err, ShouldBeNil)
	})
}

func Test_UpdateKPIData(t *testing.T) {
	r := &model.KpiData{}
	r.Day = time.Now()
	r.Mid = 111
	Convey("should return err be nil", t, func() {
		err := d.UpdateKPIData(context.TODO(), r)
		So(err, ShouldBeNil)
	})
}

func Test_UpdateKPIPoint(t *testing.T) {
	r := &model.KpiPoint{}
	r.Day = time.Now()
	r.Mid = 1
	r.Point = 100
	r.ActiveDays = 10
	r.BlockedTotal = 11
	r.VoteRadio = 60
	r.VoteTotal = 1000
	Convey("should return err be nil", t, func() {
		err := d.UpdateKPIPoint(context.TODO(), r)
		So(err, ShouldBeNil)
	})
}

func Test_UpdateCaseEndTime(t *testing.T) {
	Convey("should return err be nil", t, func() {
		num, err := d.UpdateCaseEndTime(context.TODO(), time.Now())
		So(err, ShouldBeNil)
		So(num, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func Test_UpdateCaseEndVote(t *testing.T) {
	Convey("should return err be nil", t, func() {
		num, err := d.UpdateCaseEndVote(context.TODO(), 600, time.Now().Add(time.Minute*10))
		So(err, ShouldBeNil)
		So(num, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func Test_UpdateJury(t *testing.T) {
	Convey("should return err be nil", t, func() {
		num, err := d.UpdateJury(context.TODO(), time.Now())
		So(err, ShouldBeNil)
		So(num, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func Test_UpdateJuryExpired(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.UpdateJuryExpired(context.TODO(), 88889017, time.Now())
		So(err, ShouldBeNil)
	})
}

func Test_UpdateVote(t *testing.T) {
	Convey("should return err be nil", t, func() {
		num, err := d.UpdateVote(context.TODO(), time.Now())
		So(err, ShouldBeNil)
		So(num, ShouldBeGreaterThanOrEqualTo, 0)
	})

}

func Test_LoadConf(t *testing.T) {
	Convey("should return err be nil", t, func() {
		vTotal, err := d.LoadConf(context.TODO())
		So(err, ShouldBeNil)
		So(vTotal, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func Test_JuryList(t *testing.T) {
	Convey("should return err be nil", t, func() {
		mids, err := d.JuryList(context.TODO())
		So(err, ShouldBeNil)
		So(mids, ShouldNotResemble, []int64{})
	})
}

func Test_JuryKPI(t *testing.T) {
	begin := time.Now().Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	Convey("should return err be nil", t, func() {
		res, err := d.JuryKPI(context.TODO(), begin, end)
		So(err, ShouldBeNil)
		So(res, ShouldNotResemble, []int64{})
	})
}

func Test_CountVoteTotal(t *testing.T) {
	begin := time.Now().Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	Convey("should return err be nil", t, func() {
		count, err := d.CountVoteTotal(context.TODO(), 88889017, begin, end)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func Test_CountVoteRightViolate(t *testing.T) {
	begin := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	Convey("should return err be nil", t, func() {
		count, err := d.CountVoteRightViolate(context.TODO(), 88889017, begin, end)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func Test_CountVoteRightLegal(t *testing.T) {
	begin := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	Convey("should return err be nil", t, func() {
		count, err := d.CountVoteRightLegal(context.TODO(), 88889017, begin, end)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func Test_CountBlocked(t *testing.T) {
	begin := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	Convey("should return err be nil", t, func() {
		count, err := d.CountBlocked(context.TODO(), 88889017, begin, end)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func Test_KpiPointDay(t *testing.T) {
	day := time.Now().Format("2006-01-02")
	Convey("should return err be nil", t, func() {
		_, err := d.KPIPointDay(context.TODO(), day)
		So(err, ShouldBeNil)
		// So(kp, ShouldNotBeNil)
		// So(kp, ShouldResemble,[]model.KpiPoint{})
	})
}

func TestDao_KPIPoint(t *testing.T) {
	day := time.Now().Format("2006-01-02")
	Convey("should return err be nil", t, func() {
		kp, _ := d.KPIPoint(context.TODO(), 88889017, day)
		// So(err, ShouldBeNil)
		So(kp, ShouldNotBeNil)
		So(kp, ShouldResemble, model.KpiPoint{})
	})
}

func Test_KPIList(t *testing.T) {
	Convey("should return err be nil", t, func() {
		kpis, err := d.KPIList(context.TODO(), 88889017)
		So(err, ShouldBeNil)
		So(kpis, ShouldNotBeNil)
		// So(kpis, ShouldResemble,[]model.Kpi{})
	})
}

func Test_CountVoteActive(t *testing.T) {
	begin := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	Convey("should return err be nil", t, func() {
		count, err := d.CountVoteActive(context.TODO(), 88889017, begin, end)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func TestDao_CountOpinion(t *testing.T) {
	begin := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	Convey("should return err be nil", t, func() {
		count, err := d.CountOpinion(context.TODO(), 88889017, begin, end)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func Test_OpinionQuality(t *testing.T) {
	begin := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	Convey("should return err be nil", t, func() {
		likes, hates, err := d.OpinionQuality(context.TODO(), 88889017, begin, end)
		So(err, ShouldBeNil)
		So(likes, ShouldBeGreaterThanOrEqualTo, 0)
		So(hates, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func Test_CountVoteByTime(t *testing.T) {
	begin := time.Now().AddDate(0, 0, -30)
	end := time.Now().AddDate(0, 0, 1)
	Convey("should return err be nil", t, func() {
		count, err := d.CountVoteByTime(context.TODO(), 88889017, begin, end)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThanOrEqualTo, 0)
	})
}
