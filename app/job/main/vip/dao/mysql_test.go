package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SelVipList(t *testing.T) {
	var (
		res []*model.VipUserInfo
		err error
		id  = 0
		mID = 10000
		ot  = time.Now().Format("2006-01-02")
	)
	Convey("should return true where err != nil and res not empty", t, func() {
		res, err = d.SelVipList(context.TODO(), id, mID, ot)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func Test_UpdateVipUserInfo(t *testing.T) {
	vuf := new(model.VipUserInfo)
	vuf.ID = 1
	vuf.Type = 1
	vuf.Status = 0
	Convey("should return true where err == nil", t, func() {
		tx, err := d.StartTx(context.TODO())
		So(err, ShouldBeNil)
		_, err = d.UpdateVipUserInfo(tx, vuf)
		So(err, ShouldBeNil)
	})
}

func Test_SelAppInfo(t *testing.T) {
	Convey("should return true where err == nil and res not empty", t, func() {
		res, err := d.SelAppInfo(context.TODO())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})

}

func Test_SelMaxID(t *testing.T) {
	Convey("should return true where err == nil ", t, func() {
		r, err := d.SelMaxID(context.TODO())
		So(r, ShouldBeGreaterThan, 100)
		So(err, ShouldBeNil)
	})
}

func TestDao_SelOrderMaxID(t *testing.T) {
	Convey("sel order max id", t, func() {
		_, err := d.SelOrderMaxID(context.TODO())
		So(err, ShouldBeNil)
	})
}
func TestDao_SelOldBcoinMaxID(t *testing.T) {
	Convey("sel old bcoin Max id", t, func() {
		_, err := d.SelOldBcoinMaxID(context.TODO())
		So(err, ShouldBeNil)
	})
}

func TestDao_SelBcoinMaxID(t *testing.T) {
	Convey("sel bcoin maxID", t, func() {
		_, err := d.SelBcoinMaxID(context.TODO())
		So(err, ShouldBeNil)
	})
}
func TestDao_SelChangeHistoryMaxID(t *testing.T) {
	Convey("sel change hsitory max id", t, func() {
		_, err := d.SelChangeHistoryMaxID(context.TODO())
		So(err, ShouldBeNil)
	})
}

func TestDao_SelOldOrderMaxID(t *testing.T) {
	Convey("sel old Order maxID", t, func() {
		_, err := d.SelOldOrderMaxID(context.TODO())
		So(err, ShouldBeNil)
	})
}

func TestDao_VipStatus(t *testing.T) {
	Convey("vip status", t, func() {
		_, err := d.VipStatus(context.TODO(), 1)
		So(err, ShouldBeNil)
	})
}

func TestDao_SelOldBcoinSalaryDataMaps(t *testing.T) {
	Convey("vip SelOldBcoinSalaryDataMaps", t, func() {
		_, err := d.SelOldBcoinSalaryDataMaps(context.TODO(), []int64{7593623})
		So(err, ShouldBeNil)
	})
}

func TestDao_SelBcoinSalary(t *testing.T) {
	arg := &model.QueryBcoinSalary{
		StartID: 1,
		EndID:   1,
		Status:  1,
	}
	Convey("vip SelBcoinSalary", t, func() {
		_, err := d.SelBcoinSalary(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}
func TestDao_SelOldBcoinSalary(t *testing.T) {
	arg := &model.VipBcoinSalaryMsg{
		Mid:           7593623,
		Status:        1,
		GiveNowStatus: 1,
		Payday:        "10",
		Amount:        1,
		Memo:          "memo",
	}
	Convey("vip AddBcoinSalary", t, func() {
		err := d.AddBcoinSalary(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
	Convey("vip UpdateBcoinSalary", t, func() {
		err := d.UpdateBcoinSalary(context.TODO(), "10", 7593623, 2)
		So(err, ShouldBeNil)
	})
	Convey("del bcoin salary", t, func() {
		err := d.DelBcoinSalary(context.TODO(), "10", 7593623)
		So(err, ShouldBeNil)
	})
	Convey("vip SelOldBcoinSalary", t, func() {
		_, err := d.SelOldBcoinSalary(context.TODO(), 1, 1)
		So(err, ShouldBeNil)
	})
}

func TestDao_SelBcoinSalaryData(t *testing.T) {
	Convey("SelBcoinSalaryData", t, func() {
		_, err := d.SelBcoinSalaryData(context.TODO(), 1, 1)
		So(err, ShouldBeNil)
	})
}

func TestDao_SelBcoinSalaryDataMaps(t *testing.T) {
	Convey("SelBcoinSalaryDataMaps", t, func() {
		_, err := d.SelBcoinSalaryDataMaps(context.TODO(), []int64{1})
		So(err, ShouldBeNil)
	})
}

func TestDao_SelEffectiveVipList(t *testing.T) {
	Convey("SelEffectiveVipList", t, func() {
		_, err := d.SelEffectiveVipList(context.TODO(), 1, 1)
		So(err, ShouldBeNil)
	})
}
