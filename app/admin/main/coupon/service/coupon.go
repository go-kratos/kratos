package service

import (
	"bytes"
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"go-common/app/admin/main/coupon/model"
	col "go-common/app/service/main/coupon/model"
	coumol "go-common/app/service/main/coupon/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

// AddBatchInfo add batch info.
func (s *Service) AddBatchInfo(c context.Context, b *model.CouponBatchInfo) (err error) {
	if b.StartTime >= b.ExpireTime {
		err = ecode.CouPonBatchTimeErr
		return
	}
	b.BatchToken = s.token()
	b.Ctime = xtime.Time(time.Now().Unix())
	if _, err = s.dao.AddBatchInfo(c, b); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// BatchList batch list.
func (s *Service) BatchList(c context.Context, arg *model.ArgBatchList) (res []*model.CouponBatchResp, err error) {
	var bs []*model.CouponBatchInfo
	if bs, err = s.dao.BatchList(c, arg.AppID, arg.Type); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, v := range bs {
		r := new(model.CouponBatchResp)
		r.ID = v.ID
		r.AppID = v.AppID
		r.AppName = s.allAppInfo[v.AppID]
		r.Name = v.Name
		r.BatchToken = v.BatchToken
		r.MaxCount = v.MaxCount
		r.CurrentCount = v.CurrentCount
		r.StartTime = v.StartTime
		r.ExpireTime = v.ExpireTime
		r.ExpireDay = v.ExpireDay
		r.Operator = v.Operator
		r.LimitCount = v.LimitCount
		r.UseLimitExplain = model.NoLimitExplain
		r.Amount = v.Amount
		r.FullAmount = v.FullAmount
		if r.PlatfromLimit, err = xstr.SplitInts(v.PlatformLimit); err != nil {
			log.Error("xstr.SplitInts() err[%+v] ", v.PlatformLimit, err)
			err = nil
		}
		if r.PlatfromLimit == nil {
			r.PlatfromLimit = []int64{}
		}
		r.ProdLimExplainFmt(v.ProdLimMonth, v.ProdLimRenewal) //ProductLimitExplain
		r.ProdLimMonth = v.ProdLimMonth
		r.ProdLimRenewal = v.ProdLimRenewal
		r.State = batchState(v)
		res = append(res, r)
	}
	return
}

func batchState(v *model.CouponBatchInfo) (state int8) {
	state = v.State
	now := time.Now().Unix()
	if v.ExpireDay == -1 {
		if v.ExpireTime <= now {
			state = model.CodeBatchExpire
		}
	}
	return
}

// get coupon token
func (s *Service) token() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%07d", s.r.Int63n(9999999)))
	b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
	b.WriteString(time.Now().Format("20060102150405"))
	return b.String()
}

// AllAppInfo app app info.
func (s *Service) AllAppInfo(c context.Context) map[int64]string {
	return s.allAppInfo
}

// SalaryCoupon salary coupon.
func (s *Service) SalaryCoupon(c context.Context, mid int64, ct int64, count int, token string) (err error) {
	arg := new(coumol.ArgSalaryCoupon)
	arg.Count = count
	arg.CouponType = ct
	arg.Mid = mid
	arg.Origin = model.SystemAdminSalary
	arg.BatchToken = token
	if err = s.couRPC.SalaryCoupon(c, arg); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//RunSalaryCoupon run salary coupon.
func (s *Service) RunSalaryCoupon(c context.Context, mids []int64, token string, appID int64, couponType int64, origin int64, mt string) {
	go func() {
		var (
			err     error
			msgmids = []int64{}
		)
		defer func() {
			if x := recover(); x != nil {
				log.Error("RunSalaryCoupon.GoRun arg[%s] panic[%+v]", token, x)
				log.Error("%s", debug.Stack())
			}
		}()
		for _, v := range mids {
			for i := 0; i < _maxretry; i++ {
				if err = s.couRPC.SalaryCoupon(context.Background(), &col.ArgSalaryCoupon{
					Mid:        v,
					CouponType: couponType,
					Origin:     origin,
					Count:      1,
					BatchToken: token,
					AppID:      appID,
				}); err != nil {
					time.Sleep(200 * time.Millisecond)
					continue
				}
				break
			}
			if err != nil {
				log.Error("RunSalaryCoupon faild arg[%s,%d] err[%+v] ", token, v, err)
				continue
			}
			time.Sleep(10 * time.Millisecond)
			log.Info("RunSalaryCoupon suc arg[%s,%d]", token, v)
			msgmids = append(msgmids, v)
		}
		if len(mt) > 0 && s.c.Prop.SalaryNormalMsgOpen && len(msgmids) > 0 {
			if cerr := s.msgchan.Do(c, func(c context.Context) {
				s.sendMsg(msgmids, metadata.String(c, metadata.RemoteIP), 1, true, mt)
			}); cerr != nil {
				log.Error("s.sendMsg err(%+v)", cerr)
			}
		}
	}()
}
