package service

import (
	"testing"
	"time"

	"go-common/app/job/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestServiceSalaryCoupon
func TestServiceSalaryCoupon(t *testing.T) {
	Convey("TestServiceSalaryCoupon", t, func() {
		var (
			err    error
			mid    int64 = 999
			st     int8  = model.TimingSalaryType
			vt     int8  = model.AnnualVip
			dv           = time.Now().Format("2006_01")
			atonce       = model.CouponSalaryTiming
		)
		err = s.salaryCoupon(c, mid, st, vt, dv, atonce)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceSalaryInsertAct
func TestServiceSalaryInsertAct(t *testing.T) {
	Convey("TestServiceSalaryCoupon", t, func() {
		var (
			err  error
			nvip = &model.VipUserInfoMsg{
				Mid:                  9995,
				Status:               1,
				OverdueTime:          "2018-06-11 18:27:12",
				AnnualVipOverdueTime: "2018-06-09 18:27:12",
			}
		)
		err = s.salaryInsertAct(c, nvip)
		So(err, ShouldBeNil)
		nvip.Mid = 88881
		nvip.OverdueTime = "2018-07-31 18:27:12"
		nvip.AnnualVipOverdueTime = "2018-07-31 18:27:12"
		err = s.salaryInsertAct(c, nvip)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceSalaryUpdateAct
func TestServiceSalaryUpdateAct(t *testing.T) {
	Convey("TestServiceSalaryUpdateAct", t, func() {
		var (
			err  error
			nvip = &model.VipUserInfoMsg{
				Mid:                  65,
				Status:               2,
				OverdueTime:          "2019-06-11 18:27:12",
				AnnualVipOverdueTime: "2019-06-11 18:27:12",
			}
			ovip = &model.VipUserInfoMsg{
				Mid:                  65,
				Status:               2,
				OverdueTime:          "2018-06-16 18:27:12",
				AnnualVipOverdueTime: "2018-06-09 18:27:12",
				Type:                 1,
			}
		)
		// vip -> a vip
		err = s.salaryUpdateAct(c, nvip, ovip)
		So(err, ShouldBeNil)
		// not vip -> vip
		ovip.OverdueTime = "2017-06-11 18:27:12"
		nvip.OverdueTime = "2018-07-31 18:27:12"
		nvip.AnnualVipOverdueTime = "2018-07-31 18:27:12"
		ovip.Mid = 66
		nvip.Mid = 66
		err = s.salaryUpdateAct(c, nvip, ovip)
		So(err, ShouldBeNil)
		// vip  - > a vip
		ovip.OverdueTime = "2018-08-19 18:27:12"
		nvip.AnnualVipOverdueTime = "2019-06-11 18:27:12"
		ovip.Mid = 66
		nvip.Mid = 66
		err = s.salaryUpdateAct(c, nvip, ovip)
		So(err, ShouldBeNil)
		nvip.Mid = 67
		err = s.salaryInsertAct(c, nvip)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceScanSalaryLog
func TestServiceScanSalaryLog(t *testing.T) {
	Convey("TestServiceScanSalaryLog", t, func() {
		err := s.ScanSalaryLog(c)
		So(err, ShouldBeNil)
	})
}
