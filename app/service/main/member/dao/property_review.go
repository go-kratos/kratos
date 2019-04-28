package dao

import (
	"context"

	"go-common/app/service/main/member/model"

	"github.com/pkg/errors"
)

const (
	_addUserMonitor        = "INSERT INTO user_monitor(mid, operator, remark) VALUES(?,?,?) ON DUPLICATE KEY UPDATE operator=VALUES(operator), remark=VALUES(remark), is_deleted=0"
	_isInUserMonitor       = "SELECT count(1) FROM user_monitor WHERE mid=? and is_deleted=0"
	_addPropertyReview     = "INSERT INTO user_property_review(mid, old, new, state, property, is_monitor, extra) VALUES (?,?,?,?,?,?,?)"
	_archivePropertyReview = "UPDATE user_property_review SET state=?, operator=?, remark=? WHERE mid=? AND property=? AND state=0"
)

// AddUserMonitor is.
func (d *Dao) AddUserMonitor(ctx context.Context, mid int64, operator, remark string) error {
	if _, err := d.db.Exec(ctx, _addUserMonitor, mid, operator, remark); err != nil {
		return errors.Wrapf(err, "dao add user monitor")
	}
	return nil
}

// IsInUserMonitor is.
func (d *Dao) IsInUserMonitor(ctx context.Context, mid int64) (bool, error) {
	row := d.db.QueryRow(ctx, _isInUserMonitor, mid)
	inMonitor := false
	if err := row.Scan(&inMonitor); err != nil {
		return false, errors.Wrapf(err, "dao is in user monitor")
	}
	return inMonitor, nil
}

// AddPropertyReview is.
func (d *Dao) AddPropertyReview(ctx context.Context, r *model.UserPropertyReview) error {
	if _, err := d.db.Exec(ctx, _addPropertyReview, r.Mid, r.Old, r.New, r.State, r.Property, r.IsMonitor, r.Extra); err != nil {
		return errors.Wrapf(err, "dao add user property review")
	}
	return nil
}

// ArchivePropertyReview is.
func (d *Dao) ArchivePropertyReview(ctx context.Context, mid int64, property int8) error {
	if _, err := d.db.Exec(ctx, _archivePropertyReview, model.ReviewStateArchived, "system", "已存在待审核单，本条归档处理", mid, property); err != nil {
		return errors.Wrapf(err, "dao archive property review")
	}
	return nil
}
