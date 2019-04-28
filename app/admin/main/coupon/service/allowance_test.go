package service

import (
	"go-common/app/admin/main/coupon/model"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestAddAllowanceBatchInfo
func TestAddAllowanceBatchInfo(t *testing.T) {
	Convey("TestAddAllowanceBatchInfo ", t, func() {
		var err error
		b := &model.CouponBatchInfo{
			AppID:         1,
			Name:          "test1",
			MaxCount:      1000,
			CurrentCount:  1000,
			StartTime:     1532057501,
			ExpireTime:    1542057501,
			Ver:           1,
			Operator:      "yubaihai",
			LimitCount:    20,
			FullAmount:    100,
			Amount:        20,
			State:         0,
			CouponType:    3,
			ExpireDay:     7,
			PlatformLimit: "3,4",
		}
		_, err = s.AddAllowanceBatchInfo(c, b)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUpdateAllowanceBatchInfo
func TestUpdateAllowanceBatchInfo(t *testing.T) {
	Convey("TestUpdateAllowanceBatchInfo ", t, func() {
		var err error
		b := &model.CouponBatchInfo{
			ID:            2,
			AppID:         1,
			Name:          "test2",
			MaxCount:      10000,
			CurrentCount:  1000,
			StartTime:     1532057501,
			ExpireTime:    1542057501,
			Ver:           1,
			Operator:      "yubaihai",
			LimitCount:    200,
			FullAmount:    100,
			Amount:        20,
			State:         0,
			CouponType:    3,
			PlatformLimit: "3",
		}
		err = s.UpdateAllowanceBatchInfo(c, b)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUpdateBatchStatus
func TestUpdateBatchStatus(t *testing.T) {
	Convey("TestUpdateBatchStatus ", t, func() {
		So(s.UpdateBatchStatus(c, model.BatchStateNormal, "yubaihai", 150), ShouldBeNil)
	})
}

// go test  -test.v -test.run TestBatchInfo
func TestBatchInfo(t *testing.T) {
	Convey("TestBatchInfo ", t, func() {
		res, err := s.BatchInfo(c, "test2")
		t.Logf("res(%v)", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestAllowanceSalary
func TestAllowanceSalary(t *testing.T) {
	Convey("TestAllowanceSalary ", t, func() {
		count, err := s.AllowanceSalary(c, nil, nil, []int64{332}, "allowance_test1", "vip")
		time.Sleep(time.Second * 1)
		t.Logf("count(%v)", count)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUpdateAllowanceState
func TestUpdateAllowanceState(t *testing.T) {
	Convey("TestUpdateAllowanceState ", t, func() {
		err := s.UpdateAllowanceState(c, 1, model.NotUsed, "097060140820180713120943")
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestAllowanceList
func TestAllowanceList(t *testing.T) {
	Convey("TestAllowanceList", t, func() {
		res, err := s.AllowanceList(c, &model.ArgAllowanceSearch{
			Mid:   1,
			AppID: 0,
		})
		t.Logf("count(%v)", len(res))
		So(err, ShouldBeNil)
	})
}
