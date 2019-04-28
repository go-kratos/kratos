package dao

import (
	"context"
	"database/sql"

	"go-common/library/ecode"
)

// Secret .
func (d *Dao) Secret(c context.Context, sappKey string) (res string, err error) {
	err = d.DB.Table("dm_apps").Where("appkey = ?", sappKey).Select("app_secret").Row().Scan(&res)
	if err == sql.ErrNoRows {
		err = ecode.NothingFound
	}
	return
}
