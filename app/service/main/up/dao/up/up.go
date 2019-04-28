package up

import (
	"context"

	upgrpc "go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/dao/global"
	"go-common/app/service/main/up/model"
	"go-common/library/database/sql"
)

const (
	// insert
	_inUpInfoSQL = "INSERT INTO up (mid,attribute) VALUES (?,?) ON DUPLICATE KEY UPDATE attribute=?"
	// select
	_upInfoSQL          = "SELECT id, mid, attribute FROM up WHERE mid = ?"
	_upInfoActivitysSQL = "SELECT id, mid, activity FROM up_base_info WHERE id > ? AND business_type = 1  LIMIT ?"
)

// AddUp add up.
func (d *Dao) AddUp(c context.Context, u *model.Up) (id int64, err error) {
	res, err := d.db.Exec(c, _inUpInfoSQL, u.MID, u.Attribute, u.Attribute)
	if err != nil {
		return
	}
	id, err = res.RowsAffected()
	return
}

// RawUp get attribute.
func (d *Dao) RawUp(c context.Context, mid int64) (u *model.Up, err error) {
	row := d.db.QueryRow(c, _upInfoSQL, mid)
	u = &model.Up{}
	if err = row.Scan(&u.ID, &u.MID, &u.Attribute); err != nil {
		if err == sql.ErrNoRows {
			u = nil
			err = nil
			return
		}
	}
	return
}

// UpInfoActivitys list <id, UpActivity> k-v pairs
func (d *Dao) UpInfoActivitys(c context.Context, lastID int64, ps int) (mup map[int64]*upgrpc.UpActivity, err error) {
	rows, err := global.GetUpCrmDB().Query(c, _upInfoActivitysSQL, lastID, ps)
	if err != nil {
		return
	}
	defer rows.Close()
	mup = make(map[int64]*upgrpc.UpActivity, ps)
	for rows.Next() {
		var (
			id int64
			up = new(upgrpc.UpActivity)
		)
		if err = rows.Scan(&id, &up.Mid, &up.Activity); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
			return
		}
		mup[id] = up
	}
	return
}
