package dao

import (
	"context"
	"go-common/app/interface/main/tv/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_modPage  = "SELECT id,page_id,type,title,icon,source,capacity,flexible, more, `order`,moretype,morepage,src_type FROM tv_modules WHERE page_id = ? AND deleted = 0 AND valid = 1 ORDER BY `order`"
	_passedSn = "SELECT id FROM tv_ep_season WHERE is_deleted = 0 AND `check` = 1 AND valid = 1"
)

// ModPage gets all the modules in one page
func (d *Dao) ModPage(ctx context.Context, pid int) (modules []*model.Module, err error) {
	var rows *sql.Rows
	rows, err = d.db.Query(ctx, _modPage, pid)
	if err != nil {
		log.Error("ModPage.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.Module{}
		if err = rows.Scan(&a.ID, &a.PageID, &a.Type, &a.Title, &a.Icon, &a.Source, &a.Capacity, &a.Flexible, &a.More, &a.Order, &a.MoreType, &a.MorePage, &a.SrcType); err != nil {
			log.Error("ModPage row.Scan error(%v)", err)
			return
		}
		modules = append(modules, a)
	}
	return
}

// PassedSns gets all passed seasons, to prepare their index_show data
func (d *Dao) PassedSns(ctx context.Context) (ids []int64, err error) {
	rows, err := d.db.Query(ctx, _passedSn)
	if err != nil {
		log.Error("PassedSns.Query error(%v)", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sid int64
		if err = rows.Scan(&sid); err != nil {
			log.Error("PassedSns row.Scan error(%v)", err)
			return nil, err
		}
		ids = append(ids, sid)
	}
	if err = rows.Err(); err != nil {
		log.Error("PassedSns Rows Err %v,")
	}
	return
}
