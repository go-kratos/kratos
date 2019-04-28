package dao

import (
	"context"

	"go-common/app/service/main/vip/model"
	xsql "go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_platformAll  = "SELECT id,platform_name,platform,device,mobi_app,panel_type FROM vip_platform_config WHERE is_del=0 ORDER BY id"
	_platformByID = "SELECT id,platform_name,platform,device,mobi_app,panel_type FROM vip_platform_config WHERE id =?;"
)

// PlatformAll .
func (d *Dao) PlatformAll(c context.Context) (res []*model.ConfPlatform, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _platformAll); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	res = make([]*model.ConfPlatform, 0)
	for rows.Next() {
		r := new(model.ConfPlatform)
		if err = rows.Scan(&r.ID, &r.PlatformName, &r.Platform, &r.Device, &r.MobiApp, &r.PanelType); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//PlatformByID get info by open id.
func (d *Dao) PlatformByID(c context.Context, id int64) (r *model.ConfPlatform, err error) {
	r = new(model.ConfPlatform)
	if err = d.db.QueryRow(c, _platformByID, id).
		Scan(&r.ID, &r.PlatformName, &r.Platform, &r.Device, &r.MobiApp, &r.PanelType); err != nil {
		if err == xsql.ErrNoRows {
			r = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao platform by id(%d)", id)
	}
	return
}
