package service

import (
	"context"
	"flag"
	"testing"
	"time"

	"go-common/app/job/main/coupon/conf"
	"go-common/app/job/main/coupon/model"
	"go-common/library/database/sql"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	c context.Context
)

func init() {
	var (
		err error
	)
	flag.Set("conf", "../cmd/coupon-job.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	c = context.Background()
	if s == nil {
		s = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

// go test  -test.v -test.run TestCheckCouponDeliver
func TestCheckCouponDeliver(t *testing.T) {
	Convey("TestCheckCouponDeliver ", t, func() {
		var (
			err error
		)
		arg := &model.NotifyParam{
			Mid:         1,
			CouponToken: "676289266420180402162120",
			NotifyURL:   "http://bangumi.bilibili.com/pay/inner/notify_ticket",
		}
		err = s.CheckCouponDeliver(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}

func TestCheckInUseCoupon(t *testing.T) {
	Convey("TestCheckInUseCoupon ", t, func() {
		s.CheckInUseCoupon()
	})
}

func TestNotifyproc(t *testing.T) {
	Convey("TestNotifyproc ", t, func() {
		var err error
		time.Sleep(time.Duration(s.c.Properties.NotifyTimeInterval))
		for i := 0; i < 10; i++ {
			arg := &model.NotifyParam{
				Mid:         1,
				CouponToken: "729792667120180402161647",
				NotifyURL:   "http://bangumi.bilibili.com/pay/inner/notify_ticket",
			}

			if err = s.CheckCouponDeliver(context.TODO(), arg); err != nil {
				arg.NotifyCount++
				s.notifyChan <- arg
			}
			So(err, ShouldBeNil)
		}
	})
}

func TestUpdateCoupon(t *testing.T) {
	Convey("TestUpdateCoupon ", t, func() {
		cp := &model.CouponInfo{
			CouponToken: "729792667120180402161647",
			Mid:         1,
			CouponType:  1,
			Ver:         4,
		}
		data := &model.CallBackRet{
			Ver: 3,
		}
		err := s.updateCouponState(c, cp, 2, data)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUpdateBalance

func TestUpdateBalance(t *testing.T) {
	Convey("TestUpdateBalance ", t, func() {
		var (
			tx      *sql.Tx
			mid     int64 = 1
			orderNo       = "9372774783174654609"
			ls      []*model.CouponBalanceChangeLog
			err     error
		)
		ls, err = s.dao.ConsumeCouponLog(c, mid, orderNo, model.Consume)
		So(err, ShouldBeNil)
		tx, err = s.dao.BeginTran(c)
		So(err, ShouldBeNil)
		err = s.UpdateBalance(c, tx, mid, model.Cartoon, ls, orderNo)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUpdateOrderState
func TestUpdateOrderState(t *testing.T) {
	Convey("TestUpdateOrderState ", t, func() {
		var (
			orderNo      = "6462644254161152528"
			faildOrderNo = "9176715513161453816"
			err          error
			o            *model.CouponOrder
		)
		data := &model.CallBackRet{
			Ver:    123456,
			IsPaid: 1,
		}
		o, err = s.dao.ByOrderNo(c, orderNo)
		So(err, ShouldBeNil)
		err = s.UpdateOrderState(c, o, model.PaySuccess, data)
		So(err, ShouldBeNil)
		o, err = s.dao.ByOrderNo(c, faildOrderNo)
		So(err, ShouldBeNil)
		err = s.UpdateOrderState(c, o, model.PayFaild, data)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestCouponCartoonDeliver
func TestCouponCartoonDeliver(t *testing.T) {
	Convey("TestCouponCartoonDeliver ", t, func() {
		var (
			err error
		)
		arg := &model.NotifyParam{
			CouponToken: "5586615697161708066",
			Mid:         1,
			Type:        2,
		}
		err = s.CouponCartoonDeliver(c, arg)
		So(err, ShouldBeNil)
	})
}
