package service

import (
	"context"
	"time"

	"go-common/app/job/main/vip/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

//ScanSalaryLog scan salary log.
func (s *Service) ScanSalaryLog(c context.Context) (err error) {
	var (
		dv    = time.Now().Format("2006_01")
		olds  []*model.OldSalaryLog
		size  = 1000
		endID = 0
	)
	if endID, err = s.dao.SalaryLogMaxID(context.TODO(), dv); err != nil {
		err = errors.WithStack(err)
		return
	}
	page := endID / size
	if endID%size != 0 {
		page++
	}
	for i := 0; i < page; {
		startID := i * size
		eID := (i + 1) * size
		if olds, err = s.dao.SelOldSalaryList(context.TODO(), startID, eID, dv); err != nil {
			err = errors.WithStack(err)
			return
		}
		i++
		for _, v := range olds {
			l := &model.VideoCouponSalaryLog{
				Mid:         v.Mid,
				CouponCount: v.CouponCount,
				State:       v.State,
				Type:        v.Type,
				CouponType:  model.SalaryCouponType,
			}
			if err = s.dao.AddSalaryLog(context.TODO(), l, dv); err != nil {
				err = errors.WithStack(err)
				log.Error("+%v", err)
				continue
			}
		}
	}
	return
}
