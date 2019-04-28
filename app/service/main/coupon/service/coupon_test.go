package service

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"go-common/app/service/main/coupon/conf"
	"go-common/app/service/main/coupon/model"
	"go-common/library/database/sql"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	_mid int64 = 1
	c          = context.TODO()
	s    *Service
)

func init() {
	var (
		err error
	)
	flag.Set("conf", "../cmd/coupon-service.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	c = context.Background()
	if s == nil {
		s = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

// go test  -test.v -test.run TestUserCoupon
func TestUserCoupon(t *testing.T) {
	Convey("TestUserCoupon ", t, func() {
		var (
			startTime  = time.Now().Unix()
			expireTime = time.Now().Unix() + int64(1000000)
			err        error
		)
		err = s.AddCoupon(c, _mid, startTime, expireTime, int64(1), int64(1))
		So(err, ShouldBeNil)
		err = s.dao.DelCouponsCache(c, _mid, int8(1))
		So(err, ShouldBeNil)
		cs, err := s.UserCoupon(c, _mid, int8(1))
		t.Logf("cs(%v)", len(cs))
		So(err, ShouldBeNil)
		ret, couponToken, err := s.UseCoupon(c, _mid, int64(1), "test use", fmt.Sprintf("%d", time.Now().Unix()), int8(1), int64(1))
		t.Logf("use couponToken(%s)(%v)", couponToken, err)
		So(ret == model.UseSuccess, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestCouponPage
func TestCouponPage(t *testing.T) {
	Convey("TestCouponPage ", t, func() {
		var (
			err error
		)
		count, cs, err := s.CouponPage(c, _mid, model.NotUsed, 1, 10)
		t.Logf("count(%d)", count)
		for _, v := range cs {
			t.Logf("v(%v)", v)
		}
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestAddCartoonCoupon
func TestAddCartoonCoupon(t *testing.T) {
	Convey("TestAddCartoonCoupon ", t, func() {
		var (
			err        error
			bi         *model.CouponBatchInfo
			batchToken       = "test1"
			mid        int64 = 1
			ct         int64 = 2 // cartoon
			origin     int64 = 1
			count            = 6
		)
		bi, err = s.dao.BatchInfo(c, batchToken)
		So(err, ShouldBeNil)
		So(bi, ShouldNotBeNil)
		err = s.AddCartoonCoupon(c, bi, mid, ct, origin, count)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestSalaryCoupon
func TestSalaryCoupon(t *testing.T) {
	Convey("TestSalaryCoupon ", t, func() {
		var (
			err        error
			batchToken       = "test02"
			mid        int64 = 1
			ct         int64 = 2 // cartoon
			origin     int64 = 1
			count            = 6
		)
		err = s.SalaryCoupon(c, &model.ArgSalaryCoupon{
			Mid:        mid,
			CouponType: ct,
			Origin:     origin,
			Count:      count,
			BatchToken: batchToken,
			AppID:      0,
		})
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestCarToonCouponCount
func TestCarToonCouponCount(t *testing.T) {
	Convey("TestCarToonCouponCount ", t, func() {
		var (
			err       error
			mid       int64 = 1
			ct        int8  = 2
			count     int
			_novalmid int64 = 9999999
		)
		count, err = s.CarToonCouponCount(c, mid, ct)
		t.Logf("count(%d)", count)
		So(err, ShouldBeNil)
		count, err = s.CarToonCouponCount(c, _novalmid, ct)
		t.Logf("count(%d)", count)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUpdateBalance
func TestUpdateBalance(t *testing.T) {
	Convey("TestUpdateBalance ", t, func() {
		var (
			err     error
			mid     int64 = 1
			ct      int8  = 2
			count   int64 = 2
			tx      *sql.Tx
			cs      []*model.CouponBalanceInfo
			orderNo = "test001"
		)
		tx, err = s.dao.BeginTran(c)
		So(err, ShouldBeNil)
		cs, err = s.dao.CouponBlances(c, mid, ct, time.Now().Unix())
		So(err, ShouldBeNil)
		t.Logf("cs(%d)", len(cs))
		err = s.UpdateBalance(c, tx, mid, count, cs, orderNo, ct)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestConsumeCoupon
func TestConsumeCoupon(t *testing.T) {
	Convey("TestConsumeCoupon ", t, func() {
		var (
			err     error
			mid     int64 = 1
			ct      int8  = 2
			count   int64 = 2
			cs      []*model.CouponBalanceInfo
			orderNo       = "test002"
			remake        = "我是传奇2"
			tips          = "共2话"
			ver     int64 = 1
			token   string
		)
		cs, err = s.dao.CouponBlances(c, mid, ct, time.Now().Unix())
		So(err, ShouldBeNil)
		t.Logf("cs(%d)", len(cs))
		token, err = s.ConsumeCoupon(c, mid, ct, cs, count, orderNo, remake, tips, ver)
		So(err, ShouldBeNil)
		t.Logf("token(%s)", token)
	})
}

// go test  -test.v -test.run TestCartoonUse
func TestCartoonUse(t *testing.T) {
	Convey("TestCartoonUse ", t, func() {
		var (
			err     error
			mid     int64 = 1
			ct      int8  = 2
			count   int64 = 3
			orderNo       = "test0059"
			remake        = "我是传奇5"
			tips          = "共2话"
			ver     int64 = 1
			token   string
			ret     int8
		)
		ret, token, err = s.CartoonUse(c, mid, orderNo, ct, ver, remake, tips, count)
		t.Logf("token(%s)", token)
		t.Logf("ret(%d)", ret)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestCouponCartoonPage
func TestCouponCartoonPage(t *testing.T) {
	Convey("TestCouponCartoonPage ", t, func() {
		var (
			err  error
			data *model.CouponCartoonPageResp
		)
		data, err = s.CouponCartoonPage(c, _mid, model.NotUsed, 1, 10)
		if data != nil {
			for _, v := range data.List {
				t.Logf("v(%v)", v)
			}
		}
		t.Logf("data(%v)", data)
		So(err, ShouldBeNil)
		data, err = s.CouponCartoonPage(c, _mid, model.Used, 1, 10)
		if data != nil {
			for _, v := range data.List {
				t.Logf("v(%v)", v)
			}
		}
		t.Logf("data(%v)", data)
		So(err, ShouldBeNil)
		data, err = s.CouponCartoonPage(c, _mid, model.Expire, 1, 10)
		if data != nil {
			for _, v := range data.List {
				t.Logf("v(%v)", v)
			}
		}
		t.Logf("data(%v)", data)
		So(err, ShouldBeNil)
	})
}
