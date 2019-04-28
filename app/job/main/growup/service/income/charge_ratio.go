package income

import (
	"context"

	incomeD "go-common/app/job/main/growup/dao/income"
	model "go-common/app/job/main/growup/model/income"
)

// ChargeRatioSvr charge ratio service
type ChargeRatioSvr struct {
	dao *incomeD.Dao
}

// NewChargeRatioSvr new charge ratio service
func NewChargeRatioSvr(dao *incomeD.Dao) *ChargeRatioSvr {
	return &ChargeRatioSvr{dao: dao}
}

// ArchiveChargeRatio get av charge ratio
func (p *ChargeRatioSvr) ArchiveChargeRatio(c context.Context, limit int64) (rs map[int]map[int64]*model.ArchiveChargeRatio, err error) {
	rs = make(map[int]map[int64]*model.ArchiveChargeRatio)
	var id int64
	for {
		var ros map[int]map[int64]*model.ArchiveChargeRatio
		ros, id, err = p.dao.ArchiveChargeRatio(c, id, limit)
		if err != nil {
			return
		}
		if len(ros) == 0 {
			break
		}
		for ctype, m := range ros {
			if _, ok := rs[ctype]; ok {
				for aid, ratio := range m {
					rs[ctype][aid] = ratio
				}
			} else {
				rs[ctype] = m
			}
		}
	}
	return
}

// UpChargeRatio get up charge ratio
func (p *ChargeRatioSvr) UpChargeRatio(c context.Context, limit int64) (rs map[int]map[int64]*model.UpChargeRatio, err error) {
	rs = make(map[int]map[int64]*model.UpChargeRatio)
	var id int64
	for {
		var ros map[int]map[int64]*model.UpChargeRatio
		ros, id, err = p.dao.UpChargeRatio(c, id, limit)
		if err != nil {
			return
		}
		if len(ros) == 0 {
			break
		}
		for ctype, m := range ros {
			if _, ok := rs[ctype]; ok {
				for mid, ratio := range m {
					rs[ctype][mid] = ratio
				}
			} else {
				rs[ctype] = m
			}
		}
	}
	return
}
