package service

import (
	"context"
	"fmt"

	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
)

// ProductLimit product limit.
func (s *Service) ProductLimit(c context.Context, a *model.ArgProductLimit) (err error) {
	var count, limitCount int64
	if a.PanelType != model.PanelTypeEle {
		return
	}
	if limitCount = s.c.AssociateConf.BilibiliBuyDurationMap[fmt.Sprintf("%d", a.Months)]; limitCount == 0 {
		return ecode.VipAssociateGrantDurationErr
	}
	if count, err = s.dao.CountProductBuy(c, a.Mid, a.Months, a.PanelType); err != nil {
		return
	}
	if count >= limitCount {
		err = ecode.VipActivityProductLimit
		return
	}
	return
}
