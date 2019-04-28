package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/coupon/model"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestAnalysisFile
func TestAnalysisFile(t *testing.T) {
	Convey("TestAnalysisFile ", t, func() {
		res, total, err := s.AnalysisFile(c, "/data/lv4.csv")
		t.Logf("res(%v) total(%d)", res, total)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestTokeni
func TestTokeni(t *testing.T) {
	Convey("TestTokeni ", t, func() {
		token := s.tokeni(100)
		t.Logf("token(%s)", token)
		So(token, ShouldNotBeBlank)
	})
}

// go test  -test.v -test.run TestOutFile
func TestOutFile(t *testing.T) {
	Convey("TestOutFile ", t, func() {
		err := s.OutFile(context.Background(), []byte("haha"), "/data/test.csv")
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestBatchSalary
func TestBatchSalary(t *testing.T) {
	Convey("TestbatchSalary ", t, func() {
		r, err := s.dao.BatchInfo(c, "allowance_lv41-4")
		So(err, ShouldBeNil)
		_, err = s.batchSalary(context.Background(), []int64{1, 2, 3}, "127.0.0.1", r)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestActivitySalaryCoupon
func TestActivitySalaryCoupon(t *testing.T) {
	Convey("TestActivitySalaryCoupon ", t, func() {
		err := s.ActivitySalaryCoupon(c, &model.ArgBatchSalaryCoupon{
			FileURL:     "/data/1.csv",
			Count:       1,
			BranchToken: "allowance_lv41-4",
			SliceSize:   1000,
		})
		So(err, ShouldBeNil)
	})
}
