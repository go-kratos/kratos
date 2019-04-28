package up

import (
	"context"
	"database/sql"

	"go-common/app/service/main/up/model"
)

const (
	// insert
	_inUPSwitchSQL = "INSERT INTO up_switch (mid, attribute) VALUES (?,?) ON DUPLICATE KEY UPDATE attribute=?"
	// select
	_getUPSwitchSQL = "SELECT id, mid, attribute FROM up_switch WHERE mid=?"
)

// SetSwitch add or update up switch.
func (d *Dao) SetSwitch(c context.Context, u *model.UpSwitch) (id int64, err error) {
	res, err := d.db.Exec(c, _inUPSwitchSQL, u.MID, u.Attribute, u.Attribute)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// RawUpSwitch get up switch attribute.
func (d *Dao) RawUpSwitch(c context.Context, mid int64) (u *model.UpSwitch, err error) {
	row := d.db.QueryRow(c, _getUPSwitchSQL, mid)
	u = &model.UpSwitch{}
	if err = row.Scan(&u.ID, &u.MID, &u.Attribute); err != nil {
		if err == sql.ErrNoRows {
			u = nil
			err = nil
			return
		}
	}
	return
}
