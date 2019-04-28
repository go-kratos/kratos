package kfc

import (
	"context"

	kfcmdl "go-common/app/admin/main/activity/model/kfc"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// SearchList .
func (d *Dao) SearchList(c context.Context, code string, mid int64, pn, ps int) (list []*kfcmdl.BnjKfcCoupon, err error) {
	db := d.DB
	if code != "" {
		db = db.Where("coupon_code = ?", code)
	}
	if mid != 0 {
		db = db.Where("mid = ?", mid)
	}
	offset := (pn - 1) * ps
	if err = db.Offset(offset).Limit(ps).Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		err = errors.Wrap(err, "find error")
	}
	return
}
