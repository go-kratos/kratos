package goblin

import (
	"context"

	"go-common/app/interface/main/tv/model"
)

const (
	_getChl = "SELECT id, title, `desc`, splash FROM tv_channel WHERE deleted = 0"
)

// ChlInfo .
func (d *Dao) ChlInfo(c context.Context) (chls []*model.Channel, err error) {
	rows, err := d.db.Query(c, _getChl)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		li := &model.Channel{}
		if err = rows.Scan(&li.ID, &li.Title, &li.Desc, &li.Splash); err != nil {
			return
		}
		chls = append(chls, li)
	}
	return
}
