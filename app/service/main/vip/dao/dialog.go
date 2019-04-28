package dao

import (
	"context"
	"time"

	"go-common/app/service/main/vip/model"
	xsql "go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_dialogAll = "SELECT id,app_id,platform,start_time,end_time,title,content,follow,left_button,left_link,right_button,right_link FROM vip_conf_dialog WHERE stage=true AND start_time<=? AND (end_time='1970-01-01 08:00:00' OR end_time >?)"
)

// DialogAll .
func (d *Dao) DialogAll(c context.Context) (res []*model.ConfDialog, err error) {
	var (
		rows *xsql.Rows
		curr = time.Now().Format("2006-01-02 15:04:05")
	)
	if rows, err = d.db.Query(c, _dialogAll, curr, curr); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.ConfDialog)
		if err = rows.Scan(&r.ID, &r.AppID, &r.Platform, &r.StartTime, &r.EndTime, &r.Title, &r.Content, &r.Follow, &r.LeftButton, &r.LeftLink, &r.RightButton, &r.RightLink); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}
