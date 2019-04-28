package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	pointmol "go-common/app/service/main/point/model"
	"go-common/app/service/main/vip/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

//PointRule point rule def.
func (s *Service) PointRule(c context.Context) (pr []*model.PointExchangePrice, err error) {
	pr = s.pointprices
	return
}

//BuyVipWithPoint buy vip with point .
func (s *Service) BuyVipWithPoint(c context.Context, mid int64, months int16) (err error) {
	var (
		point int32
		hv    *model.HandlerVip
	)
	if point, err = s.calcPoint(c, mid, months); err != nil {
		err = errors.WithStack(err)
		return
	}
	if hv, err = s.vipOpenByPoint(c, point, mid, months); err != nil {
		err = errors.WithStack(err)
		return
	}
	s.cache(func() {
		s.dao.DelVipInfoCache(context.TODO(), int64(hv.Mid))
	})
	s.asyncBcoin(func() {
		s.ProcesserHandler(context.TODO(), hv, "127.0.0.1")
	})
	return
}

func (s *Service) calcPoint(c context.Context, mid int64, months int16) (point int32, err error) {
	var (
		pe *model.PointExchangePrice
		ok bool
	)
	if months <= 0 {
		err = ecode.VipMonthsNotFoundErr
		return
	}
	if months > _ymonth && months%_ymonth != 0 {
		err = ecode.VipMonthsNotFoundErr
		return
	}
	if pe, ok = s.pointpricemap[months]; !ok {
		err = ecode.VipMonthsNotFoundErr
		return
	}
	point = pe.CurrentPoint
	return
}

func (s *Service) vipOpenByPoint(c context.Context, point int32, mid int64, month int16) (hv *model.HandlerVip, err error) {
	var (
		bo     = new(model.VipChangeBo)
		tx     *xsql.Tx
		status int8
	)
	if tx, err = s.dao.StartTx(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err == nil {
			if err = tx.Commit(); err != nil {
				tx.Rollback()
			}
		} else {
			tx.Rollback()
		}
	}()
	arg := &pointmol.ArgPointConsume{
		Mid:        int64(mid),
		RelationID: s.relationID(),
		ChangeType: model.ExchangeVip,
		Point:      int64(point),
	}
	if status, err = s.pointRPC.ConsumePoint(c, arg); err != nil {
		log.Error("%+v", errors.Wrapf(err, "s.pointRPC.ConsumePoint(+%v)", arg))
	}
	if status != model.PointConsumeSuc {
		err = ecode.VipPointExchangeErr
		return
	}
	hv = new(model.HandlerVip)
	bo.Mid = mid
	bo.ChangeType = model.PointChange
	bo.ChangeTime = xtime.Time(time.Now().Unix())
	bo.Days = int64(month) * model.VipDaysMonth
	bo.Months = month
	bo.Remark = _pointRemark
	bo.RelationID = arg.RelationID
	if hv, err = s.UpdateVipWithHistory(c, tx, bo); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// relationID get relation id
func (s *Service) relationID() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%010d", s.r.Int63n(9999999999)))
	b.WriteString(time.Now().Format("150405"))
	b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
	return b.String()
}
