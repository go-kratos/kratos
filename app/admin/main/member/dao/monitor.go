package dao

import (
	"context"

	"go-common/app/admin/main/member/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_monitorName = "user_monitor"
)

// Monitors is.
func (d *Dao) Monitors(ctx context.Context, mid int64, includeDeleted bool, pn, ps int) (mns []*model.Monitor, total int, err error) {
	query := d.member.Table(_monitorName).Order("id DESC")
	if !includeDeleted {
		query = query.Where("is_deleted=?", false)
	}
	if mid > 0 {
		query = query.Where("mid=?", mid)
	}
	query = query.Count(&total)
	query = query.Offset((pn - 1) * ps).Limit(ps)
	if err = query.Find(&mns).Error; err != nil {
		if err == sql.ErrNoRows {
			return []*model.Monitor{}, 0, nil
		}
		err = errors.Wrap(err, "monitors")
		return
	}
	return
}

// AddMonitor is.
func (d *Dao) AddMonitor(ctx context.Context, mid int64, operator, remark string) error {
	mn := &model.Monitor{
		Mid: mid,
	}
	ups := map[string]interface{}{
		"is_deleted": false,
		"operator":   operator,
		"remark":     remark,
	}
	if err := d.member.Table(_monitorName).Where("mid=?", mid).Assign(ups).FirstOrCreate(mn).Error; err != nil {
		return errors.Wrap(err, "add monitor")
	}
	return nil
}

// DelMonitor is.
func (d *Dao) DelMonitor(ctx context.Context, mid int64, operator, remark string) error {
	ups := map[string]interface{}{
		"is_deleted": true,
		"operator":   operator,
		"remark":     remark,
	}
	if err := d.member.Table(_monitorName).Where("mid=?", mid).UpdateColumns(ups).Error; err != nil {
		return errors.Wrap(err, "del monitor")
	}
	return nil
}
