package service

import (
	"context"
	xlog "log"
	"time"

	"go-common/app/job/main/vip/model"
	comol "go-common/app/service/main/coupon/model"
	"go-common/library/log"
)

func (s *Service) couponnotifyproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("couponnotifyproc  panic(%v)", x)
			go s.couponnotifyproc()
			log.Info("couponnotifyproc  recover")
		}
	}()
	for {
		f := <-s.notifycouponchan
		time.AfterFunc(2*time.Second, f)
	}
}

func (s *Service) couponnotify(f func()) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("couponnotifyproc panic(%v)", x)
		}
	}()
	select {
	case s.notifycouponchan <- f:
	default:
		xlog.Panic("s.couponnotifyproc chan full!")
	}
}

// CouponNotify coupon notify.
func (s *Service) CouponNotify(c context.Context, o *model.VipPayOrderNewMsg) (err error) {
	var (
		state      int8
		retrytimes = 3
	)
	if o == nil {
		return
	}
	if o.Status == model.SUCCESS {
		state = comol.AllowanceUseSuccess
	} else {
		state = comol.AllowanceUseFaild
	}
	for i := 0; i < retrytimes; i++ {
		if err = s.couponRPC.CouponNotify(c, &comol.ArgNotify{
			Mid:     o.Mid,
			OrderNo: o.OrderNo,
			State:   state,
		}); err != nil {
			log.Error("rpc.CouponNotify(%+v) error(%+v)", o, err)
			continue
		}
		break
	}
	return
}
