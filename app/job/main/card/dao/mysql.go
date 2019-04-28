package dao

import (
	"context"

	"github.com/pkg/errors"
)

const (
	_updateExpireTime = "UPDATE `user_card_equip` SET `expire_time` = ? WHERE `mid` = ?;"
)

// UpdateExpireTime update expire time.
func (d *Dao) UpdateExpireTime(c context.Context, t int64, mid int64) (err error) {
	if _, err = d.db.Exec(c, _updateExpireTime, t, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
