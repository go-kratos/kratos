package service

import (
	"context"
	"flag"
	"testing"
	"time"

	"go-common/app/admin/main/coupon/conf"
	"go-common/app/admin/main/coupon/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	c = context.TODO()
	s *Service
)

func init() {
	var (
		err error
	)
	flag.Set("conf", "../cmd/coupon-admin.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	c = context.Background()
	if s == nil {
		s = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

// go test  -test.v -test.run TestAddBatchInfo
func TestAddBatchInfo(t *testing.T) {
	Convey("TestAddBatchInfo ", t, func() {
		var err error
		b := new(model.CouponBatchInfo)
		b.AppID = int64(1)
		b.Name = "name"
		b.MaxCount = int64(100)
		b.CurrentCount = int64(0)
		b.StartTime = time.Now().Unix()
		b.ExpireTime = time.Now().Unix() + int64(10000000)
		err = s.AddBatchInfo(c, b)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestBatchList
func TestBatchList(t *testing.T) {
	Convey("TestBatchList ", t, func() {
		var err error
		_, err = s.BatchList(c, &model.ArgBatchList{AppID: 1})
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestSalaryCoupon
func TestSalaryCoupon(t *testing.T) {
	Convey("TestSalaryCoupon ", t, func() {
		var (
			err   error
			mid   int64 = 1
			ct    int64 = 2
			token       = "test03"
			count       = 20
		)
		err = s.SalaryCoupon(c, mid, ct, count, token)
		So(err, ShouldBeNil)
	})
}

func TestService_CouponViewBatchAdd(t *testing.T) {
	Convey("view batch add ", t, func() {
		arg := new(model.ArgCouponViewBatch)
		arg.AppID = 1
		arg.Name = "观影券test"
		arg.StartTime = time.Now().Unix()
		arg.ExpireTime = time.Now().AddDate(0, 0, 12).Unix()
		arg.Operator = "admin"
		arg.MaxCount = -1
		arg.LimitCount = -1
		err := s.CouponViewBatchAdd(c, arg)
		So(err, ShouldBeNil)
	})
}

func TestService_CouponViewbatchSave(t *testing.T) {
	Convey("view batch save", t, func() {
		arg := new(model.ArgCouponViewBatch)
		arg.ID = 47
		arg.AppID = 1
		arg.Name = "观影券test"
		arg.StartTime = time.Now().Unix()
		arg.ExpireTime = time.Now().AddDate(0, 0, 12).Unix()
		arg.Operator = "admin"
		arg.MaxCount = 100
		arg.LimitCount = -1
		err := s.CouponViewbatchSave(c, arg)
		So(err, ShouldBeNil)
	})
}

func TestService_CouponViewBlock(t *testing.T) {
	Convey("coupon view block", t, func() {
		mid := int64(39)
		couponToken := "326209901520180419161313"
		err := s.CouponViewBlock(c, mid, couponToken)
		So(err, ShouldBeNil)
	})
}

func TestService_CouponViewUnblock(t *testing.T) {
	Convey("coupon view un block", t, func() {
		mid := int64(39)
		couponToken := "326209901520180419161313"
		err := s.CouponViewUnblock(c, mid, couponToken)
		So(err, ShouldBeNil)
	})
}

func TestService_CouponViewList(t *testing.T) {
	Convey("coupon view list", t, func() {
		arg := new(model.ArgSearchCouponView)
		arg.Mid = 39
		arg.AppID = 1
		res, count, err := s.CouponViewList(c, arg)
		t.Logf("res:%+v count:%+v", res, count)
		So(err, ShouldBeNil)
	})
}
